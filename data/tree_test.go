package data

import "testing"

func TestNewINternal(t *testing.T) {
	page := Page(make([]byte, PageSize))

	node := NewInternal(&page)

	if node.NumKeys() != 0 {
		t.Errorf("unexpected num keys, expected %d, got %d", 0, node.NumKeys())
	}

	if node.Type() != InternalNode {
		t.Errorf("unexpected node type. expected %s, got %s", InternalNode, node.Type())
	}

	if node.IsRoot() {
		t.Errorf("unexpected IsRoot, expected %+v, got %+v", false, node.IsRoot())
	}
}

func TestInternalInsertAtEnd(t *testing.T) {
	page := Page(make([]byte, PageSize))

	node := NewInternal(&page)
	// Num Keys = 1, Right Child = 2, cell0: key=9, child=1
	copy(page[6:20], []byte{0x1, 0x0, 0x2, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x9, 0x0, 0x0, 0x0})

	tree := Tree{}
	tree.rootPageNum = 0

	tree.internalInsert(node, 16, 3)

	if node.InternalKey(1) != 16 {
		t.Errorf("unexpected value for first key, expected %d, got %d", 16, node.InternalKey(1))
	}

	if node.RightChild() != 3 {
		t.Errorf("unexpected value for right child, expected %d, got %d", 3, node.RightChild())
	}

	if node.ChildPointer(1) != 2 {
		t.Errorf("unexpected value for child pointer[1], expected %d, got %d", 2, node.ChildPointer(1))
	}

}

func TestInternalInsertAtBeginning(t *testing.T) {
	page := Page(make([]byte, PageSize))

	node := NewInternal(&page)
	// Num Keys = 1, Right Child = 2, cell0: key=16, child=1
	copy(page[6:20], []byte{0x1, 0x0, 0x2, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x9, 0x0, 0x0, 0x0})

	tree := Tree{}
	tree.rootPageNum = 0

	tree.internalInsert(node, 9, 3)

	if node.InternalKey(0) != 9 {
		t.Errorf("unexpected value for first key, expected %d, got %d", 9, node.InternalKey(0))
	}

	if node.ChildPointer(0) != 3 {
		t.Errorf("unexpected value for child pointer[1], expected %d, got %d", 3, node.ChildPointer(0))
	}
}

func TestInternalInsertInMiddle(t *testing.T) {
	page := Page(make([]byte, PageSize))

	node := NewInternal(&page)
	// Num Keys = 2, Right Child = 3, cell0: key=9, child=1, cell1: key=16, child=2
	copy(page[6:28], []byte{0x2, 0x0, 0x2, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x9, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x10, 0x0, 0x0, 0x0})

	tree := Tree{}
	tree.rootPageNum = 0

	tree.internalInsert(node, 13, 4)

	if node.InternalKey(0) != 9 {
		t.Errorf("unexpected value for first key, expected %d, got %d", 9, node.InternalKey(0))
	}

	if node.InternalKey(1) != 13 {
		t.Errorf("unexpected value for second key, expected %d, got %d", 13, node.InternalKey(1))
	}

	if node.InternalKey(2) != 16 {
		t.Errorf("unexpected value for third key, expected %d, got %d", 16, node.InternalKey(2))
	}

	if node.ChildPointer(0) != 1 {
		t.Errorf("unexpected value for child pointer[0], expected %d, got %d", 1, node.ChildPointer(4))
	}
	if node.ChildPointer(1) != 4 {
		t.Errorf("unexpected value for child pointer[1], expected %d, got %d", 1, node.ChildPointer(4))
	}
	if node.ChildPointer(2) != 2 {
		t.Errorf("unexpected value for child pointer[2], expected %d, got %d", 1, node.ChildPointer(4))
	}
}
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
	leaf := NewNode(&page)
	leaf.SetType(LeafNode)
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
		tree.leafInsert(c, test.key, test.value)
		bytes := (*leaf.page)[cellOffset : cellOffset+LeafNodeCellSize]
		if !bytesMatch(bytes, test.expectedBytes) {
			t.Errorf("incorrect bytes found at cell %d (%d), expected %+v, got %+v", test.key, c.Index, test.expectedBytes, bytes)
		}
	}
}

func TestLeafSplit(t *testing.T) {
	tree := Tree{pager: &MemoryPager{}}
	tree.pager.Page(0) // unused page to ensure parent pointers are not 0

	leafPage, _ := tree.pager.Page(1)
	tree.rootPageNum = 1
	leaf := NewLeaf(leafPage)
	leaf.SetIsRoot(true)
	for i := uint32(0); i < 341; i++ {
		c, _ := tree.leafNodeFind(leaf, i)
		tree.leafInsert(c, i, Record{i, i})
	}

	if tree.pager.GetNextUnusedPageNum() != 4 {
		t.Errorf("Unexpected next unused page, expected %d, got %d", 4, tree.pager.GetNextUnusedPageNum())
	}

	if tree.rootPageNum != 1 {
		t.Errorf("unexpected root page num, expected %d, got %d", 1, tree.rootPageNum)
	}

	rootPage, _ := tree.pager.Page(1)
	root := Node{page: rootPage}

	if root.Type() != InternalNode {
		t.Errorf("unexpected type for root, expected %s, got %s", InternalNode, root.Type())
	}

	if root.NumKeys() != 1 {
		t.Errorf("unexpected number of keys in root node, expected %d, got %d", 1, root.NumKeys())
	}

	// Because we counted up from 0 and the LeafNodeLEftSplitCount = 171, the max key in left node should be 170
	if root.InternalKey(0) != 170 {
		t.Errorf("unexpected value for key 0 in root node, expected %d, got %d", 170, root.InternalKey(0))
	}

	leftNodePageNum := root.ChildPointer(0)
	leftPage, _ := tree.pager.Page(leftNodePageNum)
	leftNode := Node{page: leftPage}

	if leftNodePageNum != 3 {
		t.Errorf("unexpected left child page num, expected %d, got %d", 3, leftNodePageNum)
	}

	if leftNode.NumCells() != 171 {
		t.Errorf("unexpected number of cells in left child, expected %d, got %d", 171, leftNode.NumCells())
	}

	if leftNode.GetMaxKey() != 170 {
		t.Errorf("unexpected value for leftNode.GetMaxKey, expected %d, got %d", 170, leftNode.GetMaxKey())
	}
	if leftNode.GetNodeKey(170) != 170 {
		t.Errorf("unexpected value for leftNode.cell[170], expected %d, got %d", 170, leftNode.GetNodeKey(171))
	}

	if leftNode.ParentPointer() != 1 {
		t.Errorf("unexpected value for left node's parent, expected %d, got %d", 1, leftNode.ParentPointer())
	}

	rightNodePageNum := root.RightChild()
	if rightNodePageNum != 2 {
		t.Errorf("unexpected right child page num, expected %d, got %d", 2, rightNodePageNum)
	}
	rightPage, _ := tree.pager.Page(rightNodePageNum)
	rightNode := Node{page: rightPage}

	if rightNode.NumCells() != 170 {
		t.Errorf("unexpected number of cells in right child, expected %d, got %d", 170, rightNode.NumCells())
	}
	if rightNode.GetMaxKey() != 340 {
		t.Errorf("unexpected value for leftNode.GetMaxKey, expected %d, got %d", 170, rightNode.GetMaxKey())
	}
	if rightNode.GetNodeKey(169) != 340 {
		t.Errorf("unexpected value for rightNode.cell[169], expected %d, got %d", 340, rightNode.GetNodeKey(169))
	}
	if rightNode.GetNodeKey(0) != 171 {
		t.Errorf("unexpected value for rightNOde.cell[0], expected %d, got %d", 171, rightNode.GetNodeKey(0))
	}
	if rightNode.ParentPointer() != 1 {
		t.Errorf("unexpected value for right node's parent, expected %d, got %d", 1, rightNode.ParentPointer())
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
