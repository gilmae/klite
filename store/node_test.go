package store

import (
	"testing"

	"github.com/gilmae/btree/data"
)

func TestNewNode(t *testing.T) {
	page := data.Page(make([]byte, data.PageSize))
	node := InititaliseNode(&page)

	if node.NextFreePosition() != 12 {
		t.Errorf("space remaining is incorrect ,expected %+v, got %+v", 12, node.NextFreePosition())
	}
	if node.SpaceRemaining() != 4084 {
		t.Errorf("space remaining is incorrect ,expected %+v, got %+v", 4084, node.SpaceRemaining())
	}
}

func TestWriteToNode(t *testing.T) {
	page := data.Page(make([]byte, data.PageSize))
	node := InititaliseNode(&page)

	payload := make([]byte, 4080)
	bytesWritten, err := node.Write(payload)

	if err != nil {
		t.Errorf("unexpected error. Got %+v", err)
	}
	if bytesWritten != 4080 {
		t.Errorf("bytesWritten is incorrect ,expected %d, got %d", 4080, bytesWritten)
	}
	if node.NextFreePosition() != 4092 {
		t.Errorf("nextFreePosition is incorrect ,expected %+v, got %+v", 4092, node.NextFreePosition())
	}
	if node.SpaceRemaining() != 4 {
		t.Errorf("space remaining is incorrect ,expected %+v, got %+v", 4, node.SpaceRemaining())
	}
}

func TestWriteToNodeWithInsufficientSpace(t *testing.T) {
	page := data.Page(make([]byte, data.PageSize))
	node := InititaliseNode(&page)

	payload := make([]byte, 4080)
	_, _ = node.Write(payload)

	bytesWritten, err := node.Write(payload)

	if err == nil {
		t.Errorf("error was expected, got nil")
	}

	if bytesWritten != 0 {
		t.Errorf("bytesWritten is incorrect, expected %d, got %d.", 0, bytesWritten)
	}

}
