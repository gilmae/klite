package data

type Tree struct {
	pager       Pager
	rootPageNum uint32
}

func NewTree(pager Pager, rootPageNum uint32) *Tree {
	return &Tree{pager: pager, rootPageNum: rootPageNum}
}

func (t *Tree) Get(key uint32) IndexItem {
	rootPage, _ := t.pager.Page(t.rootPageNum)
	root := &Node{page: rootPage}

	switch root.Type() {
	case LeafNode:
		c, found := t.leafNodeFind(root, key)
		if found {
			return c.Node.GetNodeValue(c.Index)
		}
	case InternalNode:
		c, found := t.internalNodeFind(root, key)
		if found {
			return c.Node.GetNodeValue(c.Index)
		}
	}
	return IndexItem{}

}

func (t *Tree) Insert(key uint32, data IndexItem) {
	// Add to tree
	rootPage, _ := t.pager.Page(t.rootPageNum)
	root := &Node{page: rootPage}
	var c Cursor
	var found bool
	switch root.Type() {
	case LeafNode:
		c, found = t.leafNodeFind(root, key)

	case InternalNode:
		c, found = t.internalNodeFind(root, key)

	}
	if !found {
		t.leafInsert(c, key, data)
	}

}

func (t *Tree) internalNodeFind(n *Node, key uint32) (Cursor, bool) {
	numKeys := n.NumKeys()
	minIndex, maxIndex := uint16(0), numKeys
	for minIndex != maxIndex {
		index := (minIndex + maxIndex) / 2

		keyToRight := n.InternalKey(index)
		if keyToRight >= key {
			maxIndex = index
		} else {
			minIndex += 1
		}
	}

	childPageNum := n.ChildPointer(minIndex)
	childPage, _ := t.pager.Page(childPageNum)
	child := Node{page: childPage}

	switch child.Type() {
	case LeafNode:
		return t.leafNodeFind(&child, key)
	case InternalNode:
		return t.internalNodeFind(&child, key)
	}
	return Cursor{}, false
}

func (t *Tree) internalInsert(n *Node, key uint32, rightChildPageNum uint32) {
	// If node is already full, need to call internalSplitAndInsert
	numKeys := n.NumKeys()
	if numKeys >= InternalNodeMaxCells {
		t.internalSplitAndInsert(n, key, rightChildPageNum)
		return
	}

	// Find where it should go
	minIndex, maxIndex := uint16(0), numKeys
	for minIndex != maxIndex {
		index := (minIndex + maxIndex) / 2

		keyToRight := n.InternalKey(index)
		if keyToRight >= key {
			maxIndex = index
		} else {
			minIndex += 1
		}
	}

	// Find position of first key larger than it.
	// Shuffle all keys and child pointers from that position one to the right
	// Add new key and child pointer
	n.SetNumKeys(numKeys + 1)
	// If not at the end, move cells over to make room
	for idx := numKeys; idx > minIndex; idx-- {
		n.SetInternalKey(idx, n.InternalKey(idx-1))
		n.SetChildPointer(idx, n.ChildPointer(idx-1))
	}
	n.SetInternalKey(minIndex, key)

	n.SetChildPointer(minIndex, rightChildPageNum)

}

func (t *Tree) internalSplitAndInsert(n *Node, key uint32, rightChildPageNum uint32) {
	nextPageNum := t.pager.GetNextUnusedPageNum()
	newPage, _ := t.pager.Page(nextPageNum) // TODO handle error
	newInternal := NewInternal(newPage)

	// Find where the key should go
	numKeys := n.NumKeys()
	minIndex, maxIndex := uint16(0), numKeys
	for minIndex != maxIndex {
		index := (minIndex + maxIndex) / 2

		keyToRight := n.InternalKey(index)
		if keyToRight >= key {
			maxIndex = index
		} else {
			minIndex += 1
		}
	}

	for i := int16(InternalNodeMaxCells); i >= 0; i-- {
		index := uint16(i)
		var destination *Node

		if index >= InternalNodeLeftSplitCount {
			destination = newInternal

		} else {
			destination = n

		}

		newIndex := index % InternalNodeLeftSplitCount
		if index == minIndex {
			destination.SetInternalKey(newIndex, key)
			destination.SetChildPointer(newIndex, rightChildPageNum)
		} else {
			if index > minIndex {
				destination.setInternalCell(newIndex, n.internalCell(index-1))
			} else {
				destination.setInternalCell(newIndex, n.internalCell(index))
			}
		}

	}

	n.SetNumKeys(InternalNodeLeftSplitCount - 1)
	newInternal.SetNumKeys(numKeys + 1 - (InternalNodeLeftSplitCount - 1))

	// Update Parent
	if n.IsRoot() {
		t.CreateNewRoot(nextPageNum)
	} else {
		// Update parent
		parentPage, _ := t.pager.Page(n.ParentPointer())
		parent := Node{page: parentPage}
		maxKey, _ := n.GetMaxKey()
		t.internalInsert(&parent, maxKey, nextPageNum)
		newInternal.SetParentPointer(n.ParentPointer())
	}
}

// leafNodeFind returns the position in the node the key should be in. The key may not actually be present
func (t *Tree) leafNodeFind(n *Node, key uint32) (Cursor, bool) {
	numCells := n.NumCells()
	minIndex := uint16(0)
	onePastMaxIndex := numCells

	c := Cursor{n, 0}
	for onePastMaxIndex != minIndex {

		idx := (minIndex + onePastMaxIndex) / 2
		keyAtIndex := n.GetNodeKey(idx)

		if key == keyAtIndex {
			c.Index = idx
			return c, true
		}
		if key < keyAtIndex {
			onePastMaxIndex = idx
		} else {
			minIndex = idx + 1
		}
	}
	c.Index = minIndex
	return c, false
}

func (t *Tree) leafInsert(c Cursor, key uint32, data IndexItem) {
	n := c.Node
	// If leaf is already full, need to call leafSplitAndInsert
	numCells := n.NumCells()
	if numCells >= LeafNodeMaxCells {
		t.leafSplitAndInsert(c, key, data)
		return
	}

	// Make space for new cell if it is not at the right hand end
	if c.Index < numCells {
		for i := numCells; i > c.Index; i-- {
			n.setNodeCell(i, n.getNodeCell(i-1))
		}
	}

	// Insert at the requested cell
	n.SetNumCells(numCells + 1)
	n.SetNodeKey(c.Index, key)
	n.SetNodeValue(c.Index, data)
}

func (t *Tree) leafSplitAndInsert(c Cursor, key uint32, value IndexItem) {
	nextPageNum := t.pager.GetNextUnusedPageNum()
	newPage, _ := t.pager.Page(nextPageNum) // TODO handle error
	newLeaf := NewLeaf(newPage)

	// Divide cells between nodes
	for i := int16(LeafNodeMaxCells); i >= 0; i-- {
		index := uint16(i)
		var destinationLeaf *Node
		if index >= LeafNodeLeftSplitCount {
			destinationLeaf = newLeaf
		} else {
			destinationLeaf = c.Node
		}

		newIndex := index % LeafNodeLeftSplitCount
		if index == c.Index {
			destinationLeaf.SetNodeKey(newIndex, key)
			destinationLeaf.SetNodeValue(newIndex, value)
		} else {
			if index > c.Index {
				destinationLeaf.setNodeCell(newIndex, c.Node.getNodeCell(index-1))
			} else {
				destinationLeaf.setNodeCell(newIndex, c.Node.getNodeCell(index))
			}
		}

	}

	// Update Number of cells
	c.Node.SetNumCells(LeafNodeLeftSplitCount)
	newLeaf.SetNumCells(LeafNodeRightSplitCount)

	// Update Parent
	if c.Node.IsRoot() {
		t.CreateNewRoot(nextPageNum)
	} else {
		// Update parent
		parentPage, _ := t.pager.Page(c.Node.ParentPointer())
		parent := Node{page: parentPage}
		maxKey, _ := c.Node.GetMaxKey()
		t.internalInsert(&parent, maxKey, nextPageNum)
		newLeaf.SetParentPointer(c.Node.ParentPointer())
	}
}

func (t *Tree) CreateNewRoot(rightChildPageNum uint32) {
	currentRootPage, _ := t.pager.Page(t.rootPageNum)

	leftChildPageNum := t.pager.GetNextUnusedPageNum()

	leftChildPage, _ := t.pager.Page(leftChildPageNum)

	// TODO Deal with err

	copy(*leftChildPage, *currentRootPage)
	leftChild := Node{page: leftChildPage}
	leftChild.SetIsRoot(false)
	leftChild.SetParentPointer(t.rootPageNum)

	root := NewInternal(currentRootPage)
	root.SetIsRoot(true)
	root.SetNumKeys(1)
	root.SetChildPointer(0, leftChildPageNum)
	leftChildMaxKey, _ := leftChild.GetMaxKey()
	root.SetInternalKey(0, leftChildMaxKey)
	root.SetChildPointer(1, rightChildPageNum)

	rightChildPage, _ := t.pager.Page(rightChildPageNum)
	rightChild := Node{page: rightChildPage}
	rightChild.SetParentPointer(t.rootPageNum)
}
