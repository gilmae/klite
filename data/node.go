package data

import (
	"encoding/binary"
	"unsafe"
)

type Cursor struct {
	Node  *Node
	Index uint16
}

type NodeType byte

const (
	LeafNode     NodeType = 0
	InternalNode NodeType = 1
)

func (nt NodeType) String() string {
	if nt == LeafNode {
		return "LEAF"
	} else {
		return "INTERNAL"
	}
}

// Generic Node header
const (
	NodeTypeSize        = uint16(unsafe.Sizeof(LeafNode))
	NodeTypeOffset      = 0
	IsRootSize          = uint16(unsafe.Sizeof(true))
	IsRootOffset        = NodeTypeSize
	ParentPointerSize   = uint16(unsafe.Sizeof(uint32(0)))
	ParentPointerOffset = IsRootOffset + IsRootSize
	GenericHeaderSize   = ParentPointerSize + IsRootSize + NodeTypeSize
)

type Node struct {
	page *Page
}

func NewNode(page *Page) *Node {
	return &Node{page: page}
}

func (n *Node) Type() NodeType {
	return NodeType((*n.page)[NodeTypeOffset : NodeTypeOffset+NodeTypeSize][0])
}

func (n *Node) SetType(t NodeType) {
	copy((*n.page)[NodeTypeOffset:NodeTypeOffset+NodeTypeSize], []byte{byte(t)})
}

func (n *Node) IsRoot() bool {
	return (*n.page)[IsRootOffset : IsRootOffset+IsRootSize][0] == 0x1
}

func (n *Node) SetIsRoot(isRoot bool) {
	if isRoot {
		(*n.page)[IsRootOffset] = 0x1
	} else {
		(*n.page)[IsRootOffset] = 0x0
	}
}

func (n *Node) ParentPointer() uint32 {
	return binary.LittleEndian.Uint32((*n.page)[ParentPointerOffset : ParentPointerOffset+ParentPointerSize])
}

func (n *Node) SetParentPointer(parent uint32) {
	binary.LittleEndian.PutUint32((*n.page)[ParentPointerOffset:ParentPointerOffset+ParentPointerSize], parent)
}

func (n *Node) GetMaxKey() uint32 {
	switch n.Type() {
	case InternalNode:
		return n.InternalKey(n.NumKeys() - 1)
	case LeafNode:
		return n.GetNodeKey(n.NumCells() - 1)
	}
	return 0
}
