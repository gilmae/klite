package data

import "encoding/binary"

type Record struct {
	pageNum uint32
	length  uint32
}

func Deserialise(enc []byte) Record {
	r := Record{}
	r.pageNum = binary.LittleEndian.Uint32(enc[0:4])
	r.length = binary.LittleEndian.Uint32(enc[4:8])

	return r
}

func Serialise(r Record) []byte {
	enc := make([]byte, 8)
	binary.LittleEndian.PutUint32(enc[0:4], r.pageNum)
	binary.LittleEndian.PutUint32(enc[4:8], r.length)
	return enc
}
