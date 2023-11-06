package data

import "encoding/binary"

type IndexItem struct {
	PageNum uint32
	Offset  uint16
	Length  uint32
}

func NewIndexItem(pageNum uint32, offset uint16, length uint32) IndexItem {
	return IndexItem{PageNum: pageNum, Offset: offset, Length: length}
}

func Deserialise(enc []byte) IndexItem {
	r := IndexItem{}
	r.PageNum = binary.LittleEndian.Uint32(enc[0:4])
	r.Offset = binary.LittleEndian.Uint16(enc[4:6])
	r.Length = binary.LittleEndian.Uint32(enc[6:10])

	return r
}

func Serialise(r IndexItem) []byte {
	enc := make([]byte, 10)
	binary.LittleEndian.PutUint32(enc[0:4], r.PageNum)
	binary.LittleEndian.PutUint16(enc[4:6], r.Offset)
	binary.LittleEndian.PutUint32(enc[6:10], r.Length)
	return enc
}
