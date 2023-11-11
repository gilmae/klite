package store

import (
	"encoding/binary"
	"fmt"
	"math"
	"unsafe"

	"github.com/gilmae/klite/data"
	"github.com/google/go-cmp/cmp"
)

const (
	IndexRootPageOffset = 0
	IndexRootPageSize   = uint16(unsafe.Sizeof(uint32(0)))
	StoreHeadPageOffset = IndexRootPageOffset + IndexRootPageSize
	StoreHeadPageSize   = uint16(unsafe.Sizeof(uint32(0)))
	StoreTailPageOffset = StoreHeadPageOffset + StoreHeadPageSize
	StoreTailPageSize   = uint16(unsafe.Sizeof(uint32(0)))
	NextKeyOffset       = StoreTailPageOffset + StoreTailPageSize
	NextKeySize         = uint16(unsafe.Sizeof(uint32(0)))
	StreamHeader        = NextKeyOffset + NextKeySize
)

type Stream struct {
	pager data.Pager
	page  *data.Page
	index data.Tree
}

func NewStream(p data.Pager, rootPageNum uint32) *Stream {
	stream := &Stream{pager: p}
	stream.page, _ = stream.pager.Page(rootPageNum)

	stream.index = *data.NewTree(p, stream.IndexPage())
	return stream
}

func InitialiseStream(p data.Pager) *Stream {
	stream := &Stream{pager: p}
	streamRootPage := stream.pager.GetNextUnusedPageNum()
	stream.page, _ = stream.pager.Page(streamRootPage)

	indexRootPageNum := stream.pager.GetNextUnusedPageNum()
	indexRootPage, _ := stream.pager.Page(indexRootPageNum)
	data.NewLeaf(indexRootPage)
	stream.SetIndexPage(indexRootPageNum)
	stream.index = *data.NewTree(p, indexRootPageNum)

	storeHeadPageNum := stream.pager.GetNextUnusedPageNum()
	storeHeadPage, _ := stream.pager.Page(storeHeadPageNum)
	InititaliseNode(storeHeadPage)
	stream.SetStoreHeadPage(storeHeadPageNum)
	stream.SetStoreTailPage(storeHeadPageNum)

	stream.setNextKey(0)

	return stream
}

func (s *Stream) IndexPage() uint32 {
	return binary.LittleEndian.Uint32((*s.page)[IndexRootPageOffset : IndexRootPageOffset+IndexRootPageSize])
}

func (s *Stream) SetIndexPage(pageNum uint32) {
	binary.LittleEndian.PutUint32((*s.page)[IndexRootPageOffset:IndexRootPageOffset+IndexRootPageSize], pageNum)
}

func (s *Stream) StoreHeadPage() uint32 {
	return binary.LittleEndian.Uint32((*s.page)[StoreHeadPageOffset : StoreHeadPageOffset+StoreHeadPageSize])
}

func (s *Stream) SetStoreHeadPage(pageNum uint32) {
	binary.LittleEndian.PutUint32((*s.page)[StoreHeadPageOffset:StoreHeadPageOffset+StoreHeadPageSize], pageNum)
}

func (s *Stream) StoreTailPage() uint32 {
	return binary.LittleEndian.Uint32((*s.page)[StoreTailPageOffset : StoreTailPageOffset+StoreHeadPageSize])
}

func (s *Stream) SetStoreTailPage(pageNum uint32) {
	binary.LittleEndian.PutUint32((*s.page)[StoreTailPageOffset:StoreTailPageOffset+StoreTailPageSize], pageNum)
}

func (s *Stream) NextKey() uint32 {
	return binary.LittleEndian.Uint32((*s.page)[NextKeyOffset : NextKeyOffset+NextKeySize])
}

func (s *Stream) setNextKey(key uint32) {
	binary.LittleEndian.PutUint32((*s.page)[NextKeyOffset:NextKeyOffset+NextKeySize], key)
}

func (s *Stream) Add(payload []byte) (uint32, error) {
	key := s.NextKey()
	dataWritten := 0

	curPageNum := s.StoreTailPage()
	curPage, err := s.pager.Page(curPageNum)
	if err != nil {
		return 0, err
	}

	startPageNum := curPageNum
	curNode := NewNode(curPage)
	startingOffset := curNode.NextFreePosition()
	for dataWritten < len(payload) {
		bytesAvailable := curNode.SpaceRemaining()
		if bytesAvailable <= 0 {
			nextNodePageNum := s.pager.GetNextUnusedPageNum()
			nextNodePage, err := s.pager.Page(nextNodePageNum)
			if err != nil {
				return 0, err
			}
			nextNode := InititaliseNode(nextNodePage)

			curNode.SetNext(nextNodePageNum)
			nextNode.SetPrevious(curPageNum)
			s.SetStoreTailPage(nextNodePageNum)

			curPageNum = nextNodePageNum
			curNode = nextNode
			curPage = nextNodePage
		} else {
			bytesToWrite := int(math.Min(float64(bytesAvailable), float64(len(payload))))
			if bytesToWrite+dataWritten > len(payload) {
				bytesToWrite = len(payload) - dataWritten
			}

			bytesWritten, err := curNode.Write(payload[dataWritten : dataWritten+bytesToWrite])
			if err == nil {
				dataWritten += int(bytesWritten)
			}
		}
	}

	s.index.Insert(key, data.NewIndexItem(startPageNum, startingOffset, uint32(len(payload))))
	s.setNextKey(key + 1)
	return key, nil
}

func (s *Stream) Get(key uint32) ([]byte, error) {
	indexItem := s.index.Get(key)

	if cmp.Equal(indexItem, data.IndexItem{}) {
		return nil, fmt.Errorf("key not found")
	}

	curPageNum := indexItem.PageNum
	curPage, err := s.pager.Page(curPageNum)

	curOffset := indexItem.Offset

	if err != nil {
		return nil, err
	}

	buffer := make([]byte, indexItem.Length)
	curNode := NewNode(curPage)

	totalNumBytesRead := uint32(0)

	for totalNumBytesRead < indexItem.Length {
		numBytesRead, _ := curNode.Read(curOffset, indexItem.Length-totalNumBytesRead, buffer[totalNumBytesRead:])
		totalNumBytesRead += numBytesRead

		if totalNumBytesRead < indexItem.Length {
			nextPageNum := curNode.Next()
			nextPage, _ := s.pager.Page(nextPageNum)
			curNode = NewNode(nextPage)
			curOffset = HeaderSize
		}
	}
	return buffer, nil
}
