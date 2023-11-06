package store

import (
	"testing"

	"github.com/gilmae/klite/data"
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

func TestReadFromNode(t *testing.T) {
	tests := []struct {
		offset            uint16
		length            uint32
		bytes             []byte
		expectedBuffer    []byte
		expectedBytesRead uint32
		expectedError     error
	}{
		{0, 2, []byte{2, 3}, []byte{2, 3}, 2, nil},
		{4094, 2, []byte{3, 1}, []byte{3, 1}, 2, nil},
		{4095, 2, []byte{2}, []byte{2, 0}, 1, nil},
	}

	for _, test := range tests {

		page := data.Page(make([]byte, data.PageSize))
		if test.expectedError == nil {
			copy(page[test.offset:uint32(test.offset)+test.expectedBytesRead], test.expectedBuffer)
		}

		node := InititaliseNode(&page)

		buffer := make([]byte, test.length)

		bytesRead, err := node.Read(test.offset, test.length, buffer)

		if bytesRead < test.expectedBytesRead {
			t.Errorf("incorrect bytesRead, expected %d, got %d", test.expectedBytesRead, bytesRead)
		}

		if err == nil && test.expectedError != nil {
			t.Errorf("expected an error, %+v", test.expectedError)
		}

		if err != nil && test.expectedError == nil {
			t.Errorf("unexpected error, got %+v", err)
		}

		if !bytesMatch(buffer, test.expectedBuffer) {
			t.Errorf("incorrect buffer, expected %+v, got %+v", test.expectedBuffer, buffer)
		}
	}
}

func bytesMatch(x, y []byte) bool {
	if len(x) != len(y) {
		return false
	}

	for i, b := range x {
		if b != y[i] {
			return false
		}
	}
	return true
}
