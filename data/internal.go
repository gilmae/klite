package data

import (
	"encoding/binary"
	"unsafe"
)

const (
	NumKeysSize            = uint16(unsafe.Sizeof(uint16(0)))
	NumKeysOffset          = GenericHeaderSize
	RightChildSize         = uint16(unsafe.Sizeof(uint32(0)))
	RightChildOffset       = NumKeysOffset + NumKeysSize
	InternalNodeHeaderSize = GenericHeaderSize + NumKeysSize + RightChildSize

	InternalNodeKeySize       = uint16(unsafe.Sizeof(uint32(0)))
	InternalNodeChildSize     = uint16(unsafe.Sizeof(uint32(0)))
	InternalNodeCellSize      = InternalNodeKeySize + InternalNodeChildSize
	InternalNodeSpaceForCells = PageSize - InternalNodeHeaderSize
	InternalNodeMaxCells      = InternalNodeSpaceForCells / InternalNodeCellSize
)

func NewInternal(data *Page) *Node {
	n := Node{data}
	n.SetType(InternalNode)
	n.SetIsRoot(false)
	n.SetNumKeys(0)
	return &n
}

func (n *Node) NumKeys() uint16 {
	return binary.LittleEndian.Uint16((*n.page)[NumKeysOffset : NumKeysOffset+NumKeysSize])
}

func (n *Node) SetNumKeys(c uint16) {
	binary.LittleEndian.PutUint16((*n.page)[NumKeysOffset:NumKeysOffset+NumKeysSize], c)
}

func (n *Node) RightChild() uint32 {
	return binary.LittleEndian.Uint32((*n.page)[RightChildOffset : RightChildOffset+RightChildSize])
}

func (n *Node) SetRightChild(c uint32) {
	binary.LittleEndian.PutUint32((*n.page)[RightChildOffset:RightChildOffset+RightChildSize], c)
}

func (n *Node) internalCell(cell uint16) []byte {
	cellOffset := InternalNodeHeaderSize + InternalNodeCellSize*cell
	return (*n.page)[cellOffset : cellOffset+InternalNodeCellSize]
}

func (n *Node) moveInternalCell(cell uint16, newCell uint16) {
	cellOffset := InternalNodeHeaderSize + InternalNodeCellSize*cell
	newCellOffset := InternalNodeHeaderSize + InternalNodeCellSize*newCell
	copy((*n.page)[newCellOffset:newCellOffset+InternalNodeCellSize], (*n.page)[cellOffset:cellOffset+InternalNodeCellSize])
}

func (n *Node) ChildPointer(cell uint16) uint32 {
	return binary.LittleEndian.Uint32(n.internalCell(cell)[0:InternalNodeChildSize])
}

func (n *Node) SetChildPointer(cell uint16, ptr uint32) {
	c := n.internalCell(cell)
	binary.LittleEndian.PutUint32(c[0:InternalNodeChildSize], ptr)
}

func (n *Node) InternalKey(cell uint16) uint32 {
	return binary.LittleEndian.Uint32(n.internalCell(cell)[InternalNodeChildSize : InternalNodeChildSize+InternalNodeKeySize])
}

func (n *Node) SetInternalKey(cell uint16, key uint32) {
	c := n.internalCell(cell)
	binary.LittleEndian.PutUint32(c[InternalNodeChildSize:InternalNodeChildSize+InternalNodeKeySize], key)
}
