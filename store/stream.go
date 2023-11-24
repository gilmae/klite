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
	IndexRootPageOffset        = 0
	IndexRootPageSize          = uint16(unsafe.Sizeof(uint32(0)))
	StoreHeadPageOffset        = IndexRootPageOffset + IndexRootPageSize
	StoreHeadPageSize          = uint16(unsafe.Sizeof(uint32(0)))
	StoreTailPageOffset        = StoreHeadPageOffset + StoreHeadPageSize
	StoreTailPageSize          = uint16(unsafe.Sizeof(uint32(0)))
	NextKeyOffset              = StoreTailPageOffset + StoreTailPageSize
	NextKeySize                = uint16(unsafe.Sizeof(uint32(0)))
	LastValueWrittenPageOffset = NextKeyOffset + NextKeySize
	LastValueWrittenPageSize   = uint16(unsafe.Sizeof(uint32(0)))
	LastValueWrittenPosOffset  = LastValueWrittenPageOffset + LastValueWrittenPageSize
	LastValueWrittenPosSize    = uint16(unsafe.Sizeof(uint16(0)))
	StreamHeader               = LastValueWrittenPosOffset + LastValueWrittenPosSize
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

func InitialiseStream(p data.Pager) (*Stream, uint32) {
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

	return stream, streamRootPage
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

func (s *Stream) LastValueWrittenPage() uint32 {
	return binary.LittleEndian.Uint32((*s.page)[LastValueWrittenPageOffset : LastValueWrittenPageOffset+LastValueWrittenPageSize])
}

func (s *Stream) setLastValueWrittenPage(key uint32) {
	binary.LittleEndian.PutUint32((*s.page)[LastValueWrittenPageOffset:LastValueWrittenPageOffset+LastValueWrittenPageSize], key)
}

func (s *Stream) LastValueWrittenPos() uint16 {
	return binary.LittleEndian.Uint16((*s.page)[LastValueWrittenPosOffset : LastValueWrittenPosOffset+LastValueWrittenPosSize])
}

func (s *Stream) setLastValueWrittenPos(key uint16) {
	binary.LittleEndian.PutUint16((*s.page)[LastValueWrittenPosOffset:LastValueWrittenPosOffset+LastValueWrittenPosSize], key)
}

func (s *Stream) Add(payload []byte) (uint32, error) {
	/*
		1. Get next write position
		2. If no room for header, close tail page and create new one
		3. Write header
		4. Write data, iterating through new pages as required
		5. Update header of last item to point to new item
	*/

	key := s.NextKey()
	dataWritten := 0

	curPageNum := s.StoreTailPage()
	curPage, err := s.pager.Page(curPageNum)
	if err != nil {
		return 0, err
	}

	curNode := NewNode(curPage)

	itemHeader := NewStoreItem(key, uint32(len(payload)), 0, 0)
	serialisedHeader := Serialise(itemHeader)

	// Do we have enough room for the header?
	// If not, block off remianing bytes and get a new tail
	if curNode.SpaceRemaining() < uint16(len(serialisedHeader)) {
		curNode.CloseNode()
		curPageNum, curNode, err = s.makeNewTailNode(curPageNum, curNode)
		if err != nil {
			return 0, err
		}
	}

	startPageNum := curPageNum
	startingOffset := curNode.NextFreePosition()

	// Write the header
	curNode.Write(serialisedHeader)

	for dataWritten < len(payload) {
		bytesAvailable := curNode.SpaceRemaining()
		if bytesAvailable <= 0 {
			curPageNum, curNode, err = s.makeNewTailNode(curPageNum, curNode)
			if err != nil {
				return 0, err
			}

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

	// Update the Next Item details of the last Item
	lastItemPageNum := s.LastValueWrittenPage()
	lastItemPage, _ := s.pager.Page(lastItemPageNum)
	lastItemPos := s.LastValueWrittenPos()
	lastItemHeader := ReadHeader(lastItemPage, lastItemPos)

	lastItemHeader.NextItemPageNum = startPageNum
	lastItemHeader.NextItemOffset = startingOffset

	WriteHeader(lastItemPage, lastItemHeader, lastItemPos)

	// Update the last item details of the stream with this item
	s.setLastValueWrittenPage(startPageNum)
	s.setLastValueWrittenPos(startingOffset)

	// Add to index
	s.index.Insert(key, data.NewIndexItem(startPageNum, startingOffset, uint32(len(payload))))
	s.setNextKey(key + 1)
	return key, nil
}

func (s *Stream) Get(key uint32) ([]byte, error) {

	indexItem := s.index.Get(key)

	if cmp.Equal(indexItem, data.IndexItem{}) {
		return nil, fmt.Errorf("key not found")
	}

	curOffset := indexItem.Offset
	curPageNum := indexItem.PageNum

	curPage, err := s.pager.Page(curPageNum)
	if err != nil {
		return nil, err
	}

	header := Deserialise((*curPage)[curOffset : curOffset+14])

	if header.Key != key {
		return nil, fmt.Errorf("incorrect key found in header")
	}
	if header.Length != indexItem.Length {
		return nil, fmt.Errorf("length mismatch")
	}

	curOffset += 14

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

func (s *Stream) makeNewTailNode(curPageNum uint32, curNode *Node) (uint32, *Node, error) {
	newPageNum := s.pager.GetNextUnusedPageNum()
	newPage, err := s.pager.Page(newPageNum)
	if err != nil {
		return 0, nil, err
	}
	newNode := InititaliseNode(newPage)

	curNode.SetNext(newPageNum)
	newNode.SetPrevious(curPageNum)
	s.SetStoreTailPage(newPageNum)

	return newPageNum, newNode, nil
}
