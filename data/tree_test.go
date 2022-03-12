package data

import "testing"

func TestNewRoot(t *testing.T) {
	page := Page(make([]byte, PageSize))
	node := NewInternal(&page)
	node.SetChildPointer(0, 0)
	node.SetNumKeys(1)
	node.SetChildPointer(1, 1)
	node.SetInternalKey(0, 1)

	node.SetType(InternalNode)

	if node.RightChild() != 1 {
		t.Errorf("unexpected value for node.RightChild(), expected %d, got %d", 1, node.RightChild())
	}
}

func TestNewInternal(t *testing.T) {
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

	tree := Tree{}
	tree.rootPageNum = 0
	node.SetInternalKey(0, 0)
	node.SetChildPointer(0, 0)
	node.SetNumKeys(1)

	for i := uint32(1); i < uint32(InternalNodeMaxCells); i++ {
		node.SetInternalKey(uint16(i), i)
		node.SetChildPointer(uint16(i), i)
		node.SetNumKeys(uint16(i + 1))
	}

	node.SetChildPointer(InternalNodeMaxCells, uint32(InternalNodeMaxCells))
	if node.RightChild() != uint32(InternalNodeMaxCells) {
		t.Errorf("unexpected value for right child, expected %d, got %d", uint32(InternalNodeMaxCells), node.RightChild())
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
	copy(page[6:28], []byte{0x2, 0x0, 0x3, 0x0, 0x0, 0x0, 0x1, 0x0, 0x0, 0x0, 0x9, 0x0, 0x0, 0x0, 0x2, 0x0, 0x0, 0x0, 0x10, 0x0, 0x0, 0x0})

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

func TestSplitInternal(t *testing.T) {
	p := MemoryPager{}
	tree := Tree{pager: &p, rootPageNum: 0}

	rootPage, _ := p.Page(0)
	root := NewInternal(rootPage)
	root.SetIsRoot(true)

	for i := uint16(1); i <= InternalNodeMaxCells; i++ {
		p.Page(uint32(i))
		root.SetInternalKey(i-1, uint32(i))
		root.SetNumKeys(i)
		root.SetChildPointer(i-1, uint32(i))
	}

	if root.RightChild() != uint32(0) {
		t.Errorf("unexpected value for root.RightChild, expected %d, got %d", uint32(0), root.RightChild())
	}
	root.SetNumKeys(InternalNodeMaxCells)
	p.Page(uint32(InternalNodeMaxCells + 1))
	root.SetChildPointer(InternalNodeMaxCells, uint32(InternalNodeMaxCells+1))

	expectedRightPageNum := p.GetNextUnusedPageNum()
	expectedLeftPageNum := expectedRightPageNum + 1

	// Everything in indexes greater than or equal to InternalNodeLeftSplitCount should be in right node.
	// So the expected ket will be the max key in Left Node, which should be the key in
	// InternalNodeLeftSplitCount-1
	//expectedMaxLeftKey := root.InternalKey(InternalNodeLeftSplitCount - 1)

	if root.RightChild() != uint32(InternalNodeMaxCells+1) {
		t.Errorf("unexpected value for root.RightCHild, expected %d, got %d", uint32(InternalNodeMaxCells+1), root.RightChild())
	}

	tree.internalSplitAndInsert(root, uint32(InternalNodeMaxCells+1), uint32(InternalNodeMaxCells+2))

	if tree.rootPageNum != 0 {
		t.Errorf("unexpected value for tree.rootPageNum, expected %d, got %d", 0, tree.rootPageNum)
	}

	if root.ChildPointer(0) != expectedLeftPageNum {
		t.Errorf("unexpected value for root.ChildPointer(0), expected %d, got %d", expectedLeftPageNum, root.ChildPointer(0))
	}
	if root.ChildPointer(1) != expectedRightPageNum {
		t.Errorf("unexpected value for root.RightChild(), expected %d, got %d", expectedRightPageNum, root.ChildPointer(1))
	}

	leftPage, _ := tree.pager.Page(expectedLeftPageNum)
	leftNode := Node{page: leftPage}

	for i := uint16(0); i <= 3; i++ {
		if leftNode.InternalKey(i) != uint32(i+1) {
			t.Errorf("unexpected value for leftNode.InternalKey(), expected %d, got %d", i+1, leftNode.InternalKey(i))
		}
	}
	rightPage, _ := tree.pager.Page(expectedRightPageNum)
	rightNode := Node{page: rightPage}

	for i := uint16(4); i < 7; i++ {
		actualIndex := i % InternalNodeLeftSplitCount
		if rightNode.InternalKey(actualIndex) != uint32(i+1) {
			t.Errorf("unexpected value for rightNode.InternalKey(), expected %d, got %d", i+1, rightNode.InternalKey(actualIndex))
		}
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
		//{4, 5, Record{9, 10}, []byte{5, 0, 0, 0, 9, 0, 0, 0, 10, 0, 0, 0}},
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
	for i := uint32(0); i < uint32(LeafNodeMaxCells)+1; i++ {
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

	// Because we counted up from 0 and the LeafNodeLeftSplitCount = 171, the max key in left node should be 170
	if root.InternalKey(0) != uint32(LeafNodeLeftSplitCount)-1 {
		t.Errorf("unexpected value for key 0 in root node, expected %d, got %d", uint32(LeafNodeLeftSplitCount)-1, root.InternalKey(0))
	}

	leftNodePageNum := root.ChildPointer(0)
	leftPage, _ := tree.pager.Page(leftNodePageNum)
	leftNode := Node{page: leftPage}

	if leftNodePageNum != 3 {
		t.Errorf("unexpected left child page num, expected %d, got %d", 3, leftNodePageNum)
	}

	if leftNode.NumCells() != LeafNodeLeftSplitCount {
		t.Errorf("unexpected number of cells in left child, expected %d, got %d", LeafNodeLeftSplitCount, leftNode.NumCells())
	}

	if leftNode.GetMaxKey() != uint32(LeafNodeLeftSplitCount)-1 {
		t.Errorf("unexpected value for leftNode.GetMaxKey, expected %d, got %d", uint32(LeafNodeLeftSplitCount)-1, leftNode.GetMaxKey())
	}
	if leftNode.GetNodeKey(LeafNodeLeftSplitCount-1) != uint32(LeafNodeLeftSplitCount)-1 {
		t.Errorf("unexpected value for leftNode.cell[170], expected %d, got %d", uint32(LeafNodeLeftSplitCount)-1, leftNode.GetNodeKey(LeafNodeLeftSplitCount))
	}

	if leftNode.ParentPointer() != 1 {
		t.Errorf("unexpected value for left node's parent, expected %d, got %d", 1, leftNode.ParentPointer())
	}

	rightNodePageNum := root.ChildPointer(1)
	if rightNodePageNum != 2 {
		t.Errorf("unexpected right child page num, expected %d, got %d", 2, rightNodePageNum)
	}
	rightPage, _ := tree.pager.Page(rightNodePageNum)
	rightNode := Node{page: rightPage}

	if rightNode.NumCells() != LeafNodeRightSplitCount {
		t.Errorf("unexpected number of cells in right child, expected %d, got %d", LeafNodeRightSplitCount, rightNode.NumCells())
	}
	if rightNode.GetMaxKey() != uint32(LeafNodeMaxCells) {
		t.Errorf("unexpected value for leftNode.GetMaxKey, expected %d, got %d", uint32(LeafNodeMaxCells), rightNode.GetMaxKey())
	}
	if rightNode.GetNodeKey(LeafNodeRightSplitCount-1) != uint32(LeafNodeMaxCells) {
		t.Errorf("unexpected value for last rightNode.cell key, expected %d, got %d", LeafNodeMaxCells, rightNode.GetNodeKey(LeafNodeRightSplitCount-1))
	}
	if rightNode.GetNodeKey(0) != uint32(LeafNodeLeftSplitCount) {
		t.Errorf("unexpected value for rightNode.cell[0], expected %d, got %d", uint32(LeafNodeRightSplitCount), rightNode.GetNodeKey(0))
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

	if root.ChildPointer(1) != rightPageNum {
		t.Errorf("unexpected value for root.RightChild, expected %d, got %d", rightPageNum, root.RightChild())
	}

	if root.ChildPointer(0) != rightPageNum+1 {
		t.Errorf("unexpected value for root.ChildPointer in cell 0, expected %d, got %d", rightPageNum+1, root.InternalKey(0))
	}

	if root.InternalKey(0) != 10 {
		t.Errorf("unexpected value for root.InternalKey in cell 0, expected %d, got %d", 10, root.InternalKey(0))
	}
}

func TestInsertLeaf(t *testing.T) {
	tree := Tree{pager: &MemoryPager{}, rootPageNum: 0}

	tree.Insert(1, Record{2, 3})

	rootPage, _ := tree.pager.Page(0)

	node := Node{page: rootPage}
	node.SetIsRoot(true)

	if node.NumCells() != 1 {
		t.Errorf("unexpected number of cells in node after insert, expected %d, got %d", 1, node.NumCells())
	}

	for i := uint32(2); i < uint32(LeafNodeMaxCells)*2; i++ {
		tree.Insert(i, Record{1 + 1, i + 2})
	}

	if node.Type() != InternalNode {
		t.Errorf("unexpected node type for root after inserting too many records, expected %s, got %s", InternalNode, LeafNode)
	}
}
