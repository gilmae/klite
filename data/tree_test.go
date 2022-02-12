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
	copy(page[0:], []byte{0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0})
	leaf := NewNode(&page)
	tree := Tree{}
	for _, test := range tests {
		c, found := tree.leafNodeFind(leaf, test.key)
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

	for _, test := range tests {
		cellOffset := LeafNodeHeaderSize + test.cell*LeafNodeCellSize
		c, _ := tree.leafNodeFind(leaf, test.key)
		tree.leafInsert(c, test.cell, test.key, test.value)
		bytes := (*leaf.page)[cellOffset : cellOffset+LeafNodeCellSize]
		if !bytesMatch(bytes, test.expectedBytes) {
			t.Errorf("incorrect bytes found at cell %d, expected %+v, got %+v", test.key, test.expectedBytes, bytes)
		}
	}
}

func TestCreateRoot(t *testing.T) {
	tree := Tree{pager: &MemoryPager{}, rootPageNum: 0}

	rootPage, _ := tree.pager.Page(tree.rootPageNum)
	root := Node{page: rootPage}

	root.SetType(LeafNode)
	root.SetNumKeys(1)
	root.SetNodeKey(0, 10)
	root.SetNodeValue(0, Record{2, 3})

	rightPageNum := tree.pager.GetNextUnusedPageNum()
	rightPage, _ := tree.pager.Page(rightPageNum)
	rightNode := Node{page: rightPage}

	rightNode.SetType(LeafNode)
	rightNode.SetNumKeys(1)
	rightNode.SetNodeKey(0, 40)
	rightNode.SetNodeValue(0, Record{5, 6})

	tree.CreateNewRoot(rightPageNum)

	if root.Type() != InternalNode {
		t.Errorf("unexpected node type for root, expected %s, got %s", InternalNode, root.Type())
	}

	if root.NumKeys() != 1 {
		t.Errorf("unexpected value for root.NumKeys, expected %d, got %d", 1, root.NumKeys())
	}

	if root.RightChild() != rightPageNum {
		t.Errorf("unexpected value for root.RightChild, expected %d, got %d", rightPageNum, root.RightChild())
	}

	if root.ChildPointer(0) != rightPageNum+1 {
		t.Errorf("unexpected value for root.ChildPointer in cell 0, expected %d, got %d", rightPageNum+1, root.InternalKey(0))
	}

	if root.InternalKey(0) != 10 {
		t.Errorf("unexpected value for root.InternalKey in cell 0, expected %d, got %d", 10, root.InternalKey(0))
	}
}
