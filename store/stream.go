package store

import (
	"math"

	"github.com/gilmae/btree/data"
)

type Stream struct {
	pager       data.Pager
	headPageNum uint32
	tailPageNum uint32
	nextKey     uint32
	index       data.Tree
}

func NewStream(p data.Pager, headPageNum uint32, tailPageNum uint32, nextKey uint32, indexPageNum uint32) *Stream {
	t := data.NewTree(p)
	return &Stream{pager: p, headPageNum: headPageNum, tailPageNum: tailPageNum, nextKey: nextKey, index: *t}
}

func (s *Stream) Add(payload []byte) {
	//nextKey := s.nextKey
	dataWritten := 0 //len(data)

	curPageNum := s.tailPageNum
	curPage, _ := s.pager.Page(curPageNum)
	//startPageNum := curPageNum
	curNode := NewNode(curPage)

	for dataWritten < len(payload) {
		bytesAvailable := curNode.SpaceRemaining()
		if bytesAvailable <= 0 {
			nextNodePageNum := s.pager.GetNextUnusedPageNum()
			nextNodePage, _ := s.pager.Page(nextNodePageNum)
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
	//s.index.Insert(nextKey, data.NewIndexItem(startPageNum, uint32(len(payload))))
	s.nextKey += 1
}
