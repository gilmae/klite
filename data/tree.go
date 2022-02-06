package data

type Tree struct {
	pager       Pager
	rootPageNum uint32
}

func (t Tree) Get(key uint32) Record {
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
	return Record{}

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

func (t *Tree) leafInsert(c Cursor, cell uint16, key uint32, data Record) {
	n := c.Node
	// If leaf is already full, need to call leafSplitAndInsert
	numCells := n.NumCells()
	if numCells >= LeafNodeMaxCells {
		t.leafSplitAndInsert(c, key, data)
	}

	// Make space for new cell if it is not at the right hand end
	if cell < numCells {
		for i := numCells; i > cell; i-- {
			n.setNodeCell(i, n.getNodeCell(i-1))
		}
	}

	// Insert at the requested cell
	n.SetNumCells(numCells + 1)
	n.SetNodeKey(cell, key)
	n.SetNodeValue(cell, data)
}

func (t *Tree) leafSplitAndInsert(c Cursor, key uint32, value Record) {
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
			destinationLeaf.SetNodeKey(c.Index, key)
			destinationLeaf.SetNodeValue(c.Index, value)
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
		panic("Don't know how to update parents yet")
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

	root := NewInternal(currentRootPage)
	root.SetIsRoot(true)
	root.SetNumKeys(1)
	root.SetChildPointer(0, leftChildPageNum)
	root.SetInternalKey(0, leftChild.GetMaxKey())
	root.SetRightChild(rightChildPageNum)
}
