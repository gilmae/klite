package store

import (
	"fmt"
	"math"

	"github.com/gilmae/klite/data"
	"github.com/google/go-cmp/cmp"
)

type Stream struct {
	pager       data.Pager
	headPageNum uint32
	tailPageNum uint32
	nextKey     uint32
	index       data.Tree
}

func NewStream(p data.Pager, headPageNum uint32, tailPageNum uint32, nextKey uint32, indexPageNum uint32) *Stream {
	t := data.NewTree(p, indexPageNum)
	return &Stream{pager: p, headPageNum: headPageNum, tailPageNum: tailPageNum, nextKey: nextKey, index: *t}
}

func (s *Stream) Add(payload []byte) (uint32, error) {
	key := s.nextKey
	dataWritten := 0

	curPageNum := s.tailPageNum
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
			s.tailPageNum = nextNodePageNum

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
	s.nextKey += 1
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
