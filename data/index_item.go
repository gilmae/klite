package data

import "encoding/binary"

type IndexItem struct {
	pageNum uint32
	offset  uint16
	length  uint32
}

func NewIndexItem(pageNum uint32, offset uint16, length uint32) IndexItem {
	return IndexItem{pageNum: pageNum, offset: offset, length: length}
}

func Deserialise(enc []byte) IndexItem {
	r := IndexItem{}
	r.pageNum = binary.LittleEndian.Uint32(enc[0:4])
	r.offset = binary.LittleEndian.Uint16(enc[4:6])
	r.length = binary.LittleEndian.Uint32(enc[6:10])

	return r
}

func Serialise(r IndexItem) []byte {
	enc := make([]byte, 10)
	binary.LittleEndian.PutUint32(enc[0:4], r.pageNum)
	binary.LittleEndian.PutUint16(enc[4:6], r.offset)
	binary.LittleEndian.PutUint32(enc[6:10], r.length)
	return enc
}
