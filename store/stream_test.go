package store

import (
	"testing"

	"github.com/gilmae/btree/data"
)

func TestWriteToStream(t *testing.T) {
	pager := &data.MemoryPager{}
	headPageNum := pager.GetNextUnusedPageNum()
	headPage, _ := pager.Page(headPageNum)

	head := InititaliseNode(headPage)
	indexPageNum := pager.GetNextUnusedPageNum()
	indexPage, _ := pager.Page(indexPageNum)
	indexRootNode := data.NewLeaf(indexPage)

	stream := NewStream(pager, headPageNum, headPageNum, 1, indexPageNum)

	if head.NextFreePosition() != 12 {
		t.Errorf("nextFreePosition is incorrect ,expected %+v, got %+v", 4092, head.NextFreePosition())
	}
	if head.SpaceRemaining() != 4084 {
		t.Errorf("space remaining is incorrect ,expected %+v, got %+v", 4, head.SpaceRemaining())
	}

	stream.Add(make([]byte, 4084))

	if head.NextFreePosition() != 4096 {
		t.Errorf("nextFreePosition is incorrect ,expected %+v, got %+v", 4092, head.NextFreePosition())
	}
	if head.SpaceRemaining() != 0 {
		t.Errorf("space remaining is incorrect ,expected %+v, got %+v", 4, head.SpaceRemaining())
	}

	if indexRootNode.GetMaxKey() != 1 {
		t.Errorf("incorrect number of keys in index, expected %d, got %d", 1, indexRootNode.GetMaxKey())
	}

}

func TestWriteToStreamWithInsufficientSpace(t *testing.T) {
	pager := &data.MemoryPager{}
	headPageNum := pager.GetNextUnusedPageNum()

	headPage, _ := pager.Page(headPageNum)
	head := InititaliseNode(headPage)
	indexPageNum := pager.GetNextUnusedPageNum()

	stream := NewStream(pager, headPageNum, headPageNum, 1, indexPageNum)
	indexPage, _ := pager.Page(indexPageNum)

	indexRootNode := data.NewLeaf(indexPage)

	stream.Add(make([]byte, 4083))
	stream.Add([]byte{0x1, 0x2, 0x3})

	if stream.headPageNum == stream.tailPageNum {
		t.Errorf("headPageNum equals tailPageNum, expected new page.")
	}
	if head.NextFreePosition() != 4096 {
		t.Errorf("nextFreePosition is incorrect ,expected %+v, got %+v", 4092, head.NextFreePosition())
	}
	if head.SpaceRemaining() != 0 {
		t.Errorf("space remaining is incorrect ,expected %+v, got %+v", 4, head.SpaceRemaining())
	}

	if (*headPage)[4095] != 0x1 {
		t.Errorf("incorrect byte at end of headPage, expected %+v, got %+v", 0x1, (*headPage)[4095])
	}

	if (*head).Next() != stream.tailPageNum {
		t.Errorf("incorrect value for next page after headPage, expected %d, got %d", stream.tailPageNum, (*head).Next())

	}

	tailPage, _ := pager.Page(stream.tailPageNum)
	if (*tailPage)[HeaderSize] != 0x2 {
		t.Errorf("incorrect byte after header of tailPage, expected %+v, got %+v", 0x2, (*tailPage)[HeaderSize])
	}
	tail := NewNode(tailPage)
	if (*tail).Previous() != stream.headPageNum {
		t.Errorf("incorrect value for previous page before tailPage, expected %d, got %d", stream.headPageNum, (*tail).Previous())
	}

	if indexRootNode.GetMaxKey() != 2 {
		t.Errorf("incorrect number of keys in index, expected %d, got %d", 2, indexRootNode.GetMaxKey())
	}

}
