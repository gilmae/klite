package data

import "encoding/binary"

type IndexItem struct {
	pageNum uint32
	length  uint32
}

func NewIndexItem(pageNum uint32, length uint32) IndexItem {
	return IndexItem{pageNum: pageNum, length: length}
}

func Deserialise(enc []byte) IndexItem {
	r := IndexItem{}
	r.pageNum = binary.LittleEndian.Uint32(enc[0:4])
	r.length = binary.LittleEndian.Uint32(enc[4:8])

	return r
}

func Serialise(r IndexItem) []byte {
	enc := make([]byte, 8)
	binary.LittleEndian.PutUint32(enc[0:4], r.pageNum)
	binary.LittleEndian.PutUint32(enc[4:8], r.length)
	return enc
}
