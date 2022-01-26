package data

import (
	"encoding/binary"
	"fmt"
	"unsafe"
)

// For the sake of argument, assume the values we store are 2 * uint32...a page number + length in pages
// Leaf node details
const (
	NumCellsSize   = uint16(unsafe.Sizeof(uint16(0)))
	NumCellsOffset = GenericHeaderSize

	LeafNodeHeaderSize = GenericHeaderSize + NumCellsSize

	LeafNodeKeySize   = uint16(unsafe.Sizeof(uint32(0)))
	LeafNodeKeyOffset = 0

	LeafNodeValueSize   = uint16(unsafe.Sizeof(uint32(0))) * 2
	LeafNodeValueOffset = LeafNodeKeySize + LeafNodeKeyOffset

	LeafNodeCellSize      = LeafNodeKeySize + LeafNodeValueSize
	LeafNodeSpaceForCells = PageSize - LeafNodeHeaderSize
	LeafNodeMaxCells      = LeafNodeSpaceForCells / LeafNodeCellSize
)

func (n *Node) NumCells() uint16 {
	return binary.LittleEndian.Uint16(n.page[NumCellsOffset : NumCellsOffset+NumCellsSize])
}

func (n *Node) SetNumCells(c uint16) {
	binary.LittleEndian.PutUint16(n.page[NumCellsOffset:NumCellsOffset+NumCellsSize], c)
}

func (n *Node) getNodeCell(cellNum uint16) []byte {
	cellOffset := LeafNodeHeaderSize + cellNum*LeafNodeCellSize
	return n.page[cellOffset : cellOffset+LeafNodeCellSize]
}

func (n *Node) setNodeCell(cellNum uint16, cell []byte) {
	cellOffset := LeafNodeHeaderSize + cellNum*LeafNodeCellSize
	copy(n.page[cellOffset:cellOffset+LeafNodeCellSize], cell)
}

func (n *Node) GetNodeKey(cellNum uint16) uint32 {
	return binary.LittleEndian.Uint32(n.getNodeCell(cellNum)[LeafNodeKeyOffset : LeafNodeKeyOffset+LeafNodeKeySize])
}

func (n *Node) SetNodeKey(cellNum uint16, key uint32) {
	cell := n.getNodeCell(cellNum)
	binary.LittleEndian.PutUint32(cell[LeafNodeKeyOffset:LeafNodeKeyOffset+LeafNodeKeySize], key)
	n.setNodeCell(cellNum, cell)
}

func (n *Node) leafNodeFind(key uint32) Cursor {
	numCells := n.NumCells()
	minIndex := uint16(0)
	onePastMaxIndex := numCells

	c := Cursor{n, 0}
	for onePastMaxIndex != minIndex {

		idx := (minIndex + onePastMaxIndex) / 2
		fmt.Println(onePastMaxIndex, minIndex, idx)
		keyAtIndex := n.GetNodeKey(idx)

		if key == keyAtIndex {
			c.Index = idx
			return c
		}
		if key < keyAtIndex {
			onePastMaxIndex = idx
		} else {
			minIndex = idx + 1
		}
	}
	c.Index = minIndex
	return c
}
