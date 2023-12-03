package store

import (
	"bytes"
	"testing"

	"github.com/gilmae/klite/data"
)

func TestWriteToStream(t *testing.T) {
	pager := &data.MemoryPager{}

	stream, _ := InitialiseStream(pager)

	headPage, _ := pager.Page(stream.StoreHeadPage())
	head := InititaliseNode(headPage)
	indexPage, _ := pager.Page(stream.IndexPage())
	indexRootNode := data.NewNode(indexPage)

	if head.NextFreePosition() != 12 {
		t.Errorf("nextFreePosition is incorrect ,expected %+v, got %+v", 4092, head.NextFreePosition())
	}
	if head.SpaceRemaining() != 4084 {
		t.Errorf("space remaining is incorrect ,expected %+v, got %+v", 4, head.SpaceRemaining())
	}

	key, _ := stream.Add(make([]byte, 4070))

	if head.NextFreePosition() != 4096 {
		t.Errorf("nextFreePosition is incorrect ,expected %+v, got %+v", 4092, head.NextFreePosition())
	}
	if head.SpaceRemaining() != 0 {
		t.Errorf("space remaining is incorrect ,expected %+v, got %+v", 4, head.SpaceRemaining())
	}

	maxKey, _ := indexRootNode.GetMaxKey()
	if maxKey != 0 {
		t.Errorf("incorrect max key in index, expected %d, got %d", 0, maxKey)
	}

	if key != 0 {
		t.Errorf("incorrect key returned, expected %d, got %d", 0, key)
	}

	expectedHeaderBytes := []byte{0x0, 0, 0, 0, 0xE6, 0xF, 0, 0, 0x2, 0, 0, 0, 0xC, 0}
	actualHeaderBytes := (*headPage)[12:26]

	if !bytes.Equal(expectedHeaderBytes, actualHeaderBytes) {
		t.Errorf("data header incorrect, expected %+v, got %+v", expectedHeaderBytes, actualHeaderBytes)
	}

	if stream.LastValueWrittenPage() != stream.StoreHeadPage() {
		t.Errorf("lastWrittenPage incorrect, expected %d, got %d", stream.StoreHeadPage(), stream.LastValueWrittenPage())
	}

	if stream.LastValueWrittenPos() != 12 {
		t.Errorf("lastWrittenPos incorrect, expected %d, got %d", stream.LastValueWrittenPos(), 12)
	}
}

func TestValueHeaderIsAssignedNextValueDetails(t *testing.T) {
	pager := &data.MemoryPager{}

	stream, _ := InitialiseStream(pager)

	headPage, _ := pager.Page(stream.StoreHeadPage())
	_ = InititaliseNode(headPage)
	indexPage, _ := pager.Page(stream.IndexPage())
	_ = data.NewNode(indexPage)

	stream.Add(make([]byte, 4055))
	stream.Add([]byte{0x1, 0x2, 0x3})

	valueHeader := ReadHeader(headPage, 12)

	if valueHeader.NextItemPageNum != stream.StoreHeadPage() {
		t.Errorf("newItemPageNum of first value header incorrect, expected %d, got %d", stream.StoreHeadPage(), valueHeader.NextItemPageNum)
	}

	if valueHeader.NextItemOffset != 4081 { // Start from 12 + 4055 bytes + 14 for the header
		t.Errorf("newItemOffset of first value header incorrect, expected %d, got %d", 4081, valueHeader.NextItemOffset)
	}

	stream.Add([]byte{0x4, 0x5, 0x6})

	valueHeader = ReadHeader(headPage, 4081)

	if valueHeader.NextItemPageNum != stream.StoreTailPage() {
		t.Errorf("newItemPageNum of second value header incorrect, expected %d, got %d", stream.StoreTailPage(), valueHeader.NextItemPageNum)
	}

	if valueHeader.NextItemOffset != 14 { // Start from 12 + 4055 bytes + 14 for the header
		t.Errorf("newItemOffset of second value header incorrect, expected %d, got %d", 14, valueHeader.NextItemOffset)
	}
}

func TestWriteToStreamWithInsufficientSpace(t *testing.T) {
	pager := &data.MemoryPager{}
	stream, _ := InitialiseStream(pager)

	headPage, _ := pager.Page(stream.StoreHeadPage())
	head := NewNode(headPage)
	indexPage, _ := pager.Page(stream.IndexPage())
	indexRootNode := data.NewNode(indexPage)

	stream.Add(make([]byte, 4055))
	stream.Add([]byte{0x1, 0x2, 0x3})

	if stream.StoreHeadPage() == stream.StoreTailPage() {
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

	if (*head).Next() != stream.StoreTailPage() {
		t.Errorf("incorrect value for next page after headPage, expected %d, got %d", stream.StoreHeadPage(), (*head).Next())

	}

	tailPage, _ := pager.Page(stream.StoreTailPage())
	if (*tailPage)[HeaderSize] != 0x2 {
		t.Errorf("incorrect byte after header of tailPage, expected %+v, got %+v", 0x2, (*tailPage)[HeaderSize])
	}
	tail := NewNode(tailPage)
	if (*tail).Previous() != stream.StoreHeadPage() {
		t.Errorf("incorrect value for previous page before tailPage, expected %d, got %d", stream.StoreHeadPage(), (*tail).Previous())
	}
	maxKey, _ := indexRootNode.GetMaxKey()
	if maxKey != 1 {
		t.Errorf("incorrect max key in index, expected %d, got %d", 1, maxKey)
	}

}

func TestReadFromStream(t *testing.T) {

	//headerBuffer := []byte{0, 0, 0, 0, 1, 0, 0, 0}
	expectedBuffer := []byte{0x1, 0x2, 0x3, 0x4}

	pager := &data.MemoryPager{}
	stream, _ := InitialiseStream(pager)

	headPage, _ := pager.Page(stream.StoreHeadPage())

	indexPage, _ := pager.Page(stream.IndexPage())

	// Add the actual value to the node
	copy((*headPage)[34:38], expectedBuffer)

	// Add the value header to the node
	copy((*headPage)[20:28], []byte{0, 0, 0, 0, 4, 0, 0, 0})
	indexRootNode := data.NewNode(indexPage)

	indexRootNode.SetNodeKey(0, 0)
	indexRootNode.SetNodeValue(0, data.IndexItem{PageNum: stream.StoreHeadPage(), Offset: 20, Length: 4})
	indexRootNode.SetNumKeys(1)

	actualBuffer, err := stream.Get(0)

	if err != nil {
		t.Errorf("unexpected error, got %+v", err)
	}
	if !bytes.Equal(actualBuffer, expectedBuffer) {
		t.Errorf("incorrect buffer returned, expected %+v, got %+v", expectedBuffer, actualBuffer)
	}
}

func TestReadFromStreamUsingInvalidKey(t *testing.T) {
	expectedBuffer := []byte{0x1, 0x2, 0x3, 0x4}

	pager := &data.MemoryPager{}
	stream, _ := InitialiseStream(pager)

	headPage, _ := pager.Page(stream.StoreHeadPage())

	indexPage, _ := pager.Page(stream.IndexPage())
	copy((*headPage)[20:24], expectedBuffer)
	indexRootNode := data.NewLeaf(indexPage)

	indexRootNode.SetNodeKey(0, 0)
	indexRootNode.SetNodeValue(0, data.IndexItem{PageNum: stream.StoreHeadPage(), Offset: 20, Length: 4})
	indexRootNode.SetNumKeys(1)

	actualBuffer, err := stream.Get(2)
	if err == nil {
		t.Errorf("expected an error for invalid key")
	}

	if actualBuffer != nil {
		t.Errorf("incorrect buffer returned, got %+v", actualBuffer)
	}
}

func TestReadMultiple(t *testing.T) {
	pager := &data.MemoryPager{}
	stream, _ := InitialiseStream(pager)

	expectedItem1 := []byte{0x1, 0x2, 0x3}
	expectedItem2 := []byte{0x4, 0x5, 0x6}
	stream.Add(expectedItem1)
	stream.Add(expectedItem2)

	items, err := stream.GetFrom(0, 2)

	if err != nil {
		t.Errorf("Unexpected error, got %s", err)
	}

	if len(items) < 2 {
		t.Errorf("incorrect number of items recived, expected %d, got %d", 2, len(items))
	}

	if !bytes.Equal(expectedItem1, items[0]) {
		t.Errorf("incorrect value received, expected %+v, got %+v", expectedItem1, items[0])
	}

	if !bytes.Equal(expectedItem2, items[1]) {
		t.Errorf("incorrect value received, expected %+v, got %+v", expectedItem2, items[1])
	}
}
