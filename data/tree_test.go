package data

import "testing"

func TestLeafFind(t *testing.T) {
	tests := []struct {
		key              uint32
		expectedPosition uint16
		expectedFound    bool
	}{
		{1, 0, true},
		{2, 1, true},
		{5, 3, true},
		{3, 2, true},
		{6, 4, false}, // key is not present but would be in index 4 if it were inserted
		{4, 3, false}, // key is not present but would be in index 3 if it were inserted
	}

	page := Page(make([]byte, PageSize))
	copy(page[0:], []byte{0, 0, 0, 0, 0, 0, 4, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0})
	leaf := NewNode(&page)
	for _, test := range tests {
		c, found := leafNodeFind(leaf, test.key)
		if c.Index != test.expectedPosition {
			t.Errorf("unexpected position for key %d, expected %d, got %d", test.key, test.expectedPosition, c.Index)
		}

		if found != test.expectedFound {
			t.Errorf("unexpected found flag for key %d, expected %t, got %t", test.key, test.expectedFound, found)
		}
	}
}

func TestLeafInsert(t *testing.T) {
	page := Page(make([]byte, PageSize))
	copy(page[0:], []byte{0, 0, 0, 0, 0, 0, 4, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0})
	leaf := NewNode(&page)

	tests := []struct {
		cell          uint16
		key           uint32
		value         Record
		expectedBytes []byte
	}{
		{0, 1, Record{1, 2}, []byte{1, 0, 0, 0, 1, 0, 0, 0, 2, 0, 0, 0}},
		{1, 2, Record{3, 4}, []byte{2, 0, 0, 0, 3, 0, 0, 0, 4, 0, 0, 0}},
		{2, 3, Record{5, 6}, []byte{3, 0, 0, 0, 5, 0, 0, 0, 6, 0, 0, 0}},
		{3, 4, Record{7, 8}, []byte{4, 0, 0, 0, 7, 0, 0, 0, 8, 0, 0, 0}},
		{4, 5, Record{9, 10}, []byte{5, 0, 0, 0, 9, 0, 0, 0, 10, 0, 0, 0}},
	}

	tree := Tree{}

	// tree.leafInsert(leaf, 2, 3, Record{5, 6})
	// tree.leafInsert(leaf, 4, 5, Record{9, 20})

	for _, test := range tests {
		cellOffset := LeafNodeHeaderSize + test.cell*LeafNodeCellSize
		tree.leafInsert(leaf, test.cell, test.key, test.value)
		bytes := (*leaf.page)[cellOffset : cellOffset+LeafNodeCellSize]
		if !bytesMatch(bytes, test.expectedBytes) {
			t.Errorf("incorrect bytes found at cell %d, expected %+v, got %+v", test.key, test.expectedBytes, bytes)
		}
	}
}
