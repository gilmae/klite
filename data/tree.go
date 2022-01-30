package data

type Tree struct {
	pager Pager
	root  *Node
}

func (t Tree) Get(key uint32) Record {
	c, found := leafNodeFind(t.root, key)
	if found {
		return t.root.GetNodeValue(c.Index)
	}
	return Record{}
}

// leafNodeFind returns the position in the node the key should be in. The key may not actually be present
func leafNodeFind(n *Node, key uint32) (Cursor, bool) {
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

func (t *Tree) leafInsert(n *Node, cell uint16, key uint32, data Record) {
	// If leaf is already full, need to call leafSplitAndInsert
	numCells := n.NumCells()
	if numCells >= LeafNodeMaxCells {
		return
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

func (l *Node) leafSplitAndInsert(key uint32, value Record) {

}

// func (t *Tree) splitLeafAndInsert(cursor *Cursor, key uint32, value Record) {
// 	nextPageNum := t.pager.GetNextUnusedPageNum()
// 	newPage, _ := t.pager.Page(nextPageNum) // TODO handle error
// 	newLeaf := Node{newPage}

// 	newLeaf
// }
