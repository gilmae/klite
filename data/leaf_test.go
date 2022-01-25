package data

import "testing"

func TestSetKey(t *testing.T) {
	tests := []struct {
		cell           uint16
		key            uint32
		expectedValues []byte
	}{
		{0, 7, []byte{0, 0, 0, 7}},
		{1, 256, []byte{0, 0, 1, 0}},
	}

	for i, test := range tests {
		p := [PageSize]byte{}
		leaf := NewNode(p)
		leaf.SetType(LeafNode)

		leaf.SetNodeKey(test.cell, test.key)
		offset := 8 + test.cell*12
		bytes := leaf.page[offset : offset+4]

		if !bytesMatch(bytes, test.expectedValues) {
			t.Errorf("unexpected key got test %d, expected %+v, got %+v", i, test.expectedValues, bytes)
		}
	}
}

func TestGetKey(t *testing.T) {
	tests := []struct {
		data          [4]byte
		cell          uint16
		expectedValue uint32
	}{
		{[4]byte{0, 0, 0, 7}, 0, 7},
		{[4]byte{0, 0, 1, 0}, 1, 256},
	}

	for i, test := range tests {
		p := [PageSize]byte{}
		offset := 8 + test.cell*12
		copy(p[offset:offset+4], test.data[:])
		leaf := NewNode(p)
		if leaf.GetNodeKey(test.cell) != test.expectedValue {
			t.Errorf("unexpected value for test %d, expected %d, got %d", i, test.expectedValue, leaf.GetNodeKey(test.cell))
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
