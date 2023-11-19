package store

import (
	"encoding/binary"
	"fmt"
	"unsafe"

	"github.com/gilmae/klite/data"
)

/*
Node Header
+++++++++++

PrevPage: Holds page of the previous node
NextPage: Holds page of the next node
*/

const (
	PrevPageOffset         = 0
	PrevPageSize           = uint16(unsafe.Sizeof(uint32(0)))
	NextPageOffset         = PrevPageSize
	NextPageSize           = uint16(unsafe.Sizeof(uint32(0)))
	NextFreePositionOffset = NextPageOffset + NextPageSize
	NextFreePositionSize   = uint16(unsafe.Sizeof(uint32(0)))
	//LastValueWritenOffset  = NextFreePositionOffset + NextFreePositionSize
	//LastValueWrittenSize   = uint16(unsafe.Sizeof(uint32(0)))
	HeaderSize = NextFreePositionOffset + NextFreePositionSize //LastValueWritenOffset + LastValueWrittenSize
)

type Node struct {
	page *data.Page
}

func NewNode(page *data.Page) *Node {
	return &Node{page: page}
}

func InititaliseNode(p *data.Page) *Node {
	n := Node{p}
	n.SetNextFreePosition(HeaderSize)
	return &n
}

func (n *Node) Previous() uint32 {
	return binary.LittleEndian.Uint32((*n.page)[PrevPageOffset : PrevPageOffset+PrevPageSize])
}

func (n *Node) SetPrevious(previous uint32) {
	binary.LittleEndian.PutUint32((*n.page)[PrevPageOffset:PrevPageOffset+PrevPageSize], previous)
}

func (n *Node) Next() uint32 {
	return binary.LittleEndian.Uint32((*n.page)[NextPageOffset : NextPageOffset+NextPageSize])
}

func (n *Node) SetNext(next uint32) {
	binary.LittleEndian.PutUint32((*n.page)[NextPageOffset:NextPageOffset+NextPageSize], next)
}

func (n *Node) NextFreePosition() uint16 {
	return binary.LittleEndian.Uint16((*n.page)[NextFreePositionOffset : NextFreePositionOffset+NextFreePositionSize])
}

func (n *Node) SetNextFreePosition(nextFreePos uint16) {
	binary.LittleEndian.PutUint16((*n.page)[NextFreePositionOffset:NextFreePositionOffset+NextFreePositionSize], nextFreePos)
}

func (n *Node) CloseNode() {
	n.SetNextFreePosition(data.PageSize)
}

// func (n *Node) LastValueWrittenPosition() uint16 {
// 	return binary.LittleEndian.Uint16((*n.page)[LastValueWritenOffset : LastValueWritenOffset+LastValueWrittenSize])
// }

// func (n *Node) SetLastValueWrittenPosition(lastValueWritten uint16) {
// 	binary.LittleEndian.PutUint16((*n.page)[LastValueWritenOffset:LastValueWritenOffset+LastValueWrittenSize], lastValueWritten)
// }

func (n *Node) SpaceRemaining() uint16 {
	return data.PageSize - n.NextFreePosition()
}

func (n *Node) Write(data []byte) (uint16, error) {
	nextFree := int(n.NextFreePosition())
	if n.SpaceRemaining() < uint16(len(data)) {
		return 0, fmt.Errorf("insufficent space remaining, %d bytes free", n.SpaceRemaining())
	}
	copy((*n.page)[nextFree:nextFree+len(data)], data)
	n.SetNextFreePosition(uint16(nextFree + len(data)))
	return uint16(len(data)), nil
}

func (n *Node) Read(offset uint16, length uint32, buffer []byte) (uint32, error) {
	if uint32(len(buffer)) < length {
		return 0, fmt.Errorf("buffer too small")
	}

	numBytesToRead := length

	if uint32(offset)+length > data.PageSize {
		numBytesToRead = uint32(data.PageSize - offset)
	}

	copy(buffer, (*n.page)[uint32(offset):uint32(offset)+numBytesToRead])
	return numBytesToRead, nil
}
