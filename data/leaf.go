package data

import (
	"encoding/binary"
	"unsafe"
)

// For the sake of argument, assume the values we store are 2 * uint32...a page number + length in pages
// Leaf node details
const (
	NumCellsSize   = uint16(unsafe.Sizeof(uint16(0)))
	NumCellsOffset = GenericHeaderSize

	NextLeafPointerSize   = uint16(unsafe.Sizeof(uint32(0)))
	NextLeafPointerOffset = GenericHeaderSize + NumCellsOffset

	LeafNodeHeaderSize = GenericHeaderSize + NumCellsSize + NextLeafPointerSize

	LeafNodeKeySize   = uint16(unsafe.Sizeof(uint32(0)))
	LeafNodeKeyOffset = 0

	LeafNodeValueSize   = uint16(unsafe.Sizeof(Record{}))
	LeafNodeValueOffset = LeafNodeKeySize + LeafNodeKeyOffset

	LeafNodeCellSize        = LeafNodeKeySize + LeafNodeValueSize
	LeafNodeSpaceForCells   = PageSize - LeafNodeHeaderSize
	LeafNodeMaxCells        = LeafNodeSpaceForCells / LeafNodeCellSize
	LeafNodeRightSplitCount = (LeafNodeMaxCells + 1) / 2
	LeafNodeLeftSplitCount  = LeafNodeMaxCells + 1 - LeafNodeRightSplitCount
)

func NewLeaf(p *Page) *Node {
	n := Node{p}
	n.SetType(LeafNode)
	n.SetIsRoot(false)

	n.SetNumCells(0)
	//n.SetNextLeaf(0)
	return &n
}

func (n *Node) NumCells() uint16 {
	return binary.LittleEndian.Uint16((*n.page)[NumCellsOffset : NumCellsOffset+NumCellsSize])
}

func (n *Node) SetNumCells(c uint16) {
	binary.LittleEndian.PutUint16((*n.page)[NumCellsOffset:NumCellsOffset+NumCellsSize], c)
}

func (n *Node) getNodeCell(cellNum uint16) []byte {
	cellOffset := LeafNodeHeaderSize + cellNum*LeafNodeCellSize
	return (*n.page)[cellOffset : cellOffset+LeafNodeCellSize]
}

func (n *Node) setNodeCell(cellNum uint16, cell []byte) {
	cellOffset := LeafNodeHeaderSize + cellNum*LeafNodeCellSize
	copy((*n.page)[cellOffset:cellOffset+LeafNodeCellSize], cell)
}

func (n *Node) GetNodeKey(cellNum uint16) uint32 {
	return binary.LittleEndian.Uint32(n.getNodeCell(cellNum)[LeafNodeKeyOffset : LeafNodeKeyOffset+LeafNodeKeySize])
}

func (n *Node) SetNodeKey(cellNum uint16, key uint32) {
	cell := n.getNodeCell(cellNum)
	binary.LittleEndian.PutUint32(cell[LeafNodeKeyOffset:LeafNodeKeyOffset+LeafNodeKeySize], key)
}

func (n *Node) GetNodeValue(cellNum uint16) Record {
	cell := n.getNodeCell(cellNum)
	return Deserialise(cell[LeafNodeKeySize:])
}

func (n *Node) SetNodeValue(cellNum uint16, r Record) {
	cell := n.getNodeCell(cellNum)
	copy(cell[LeafNodeKeySize:LeafNodeKeySize+LeafNodeValueSize], Serialise(r))
}
