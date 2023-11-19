package store

// nodule header
/*
NextRecordPage
NextRecordOffset
PrevRecordPage
PrevRecordOffset


*/

import "encoding/binary"

type StoreItem struct {
	Key             uint32
	Length          uint32
	NextItemPageNum uint32
	NextItemOffset  uint16
}

func NewStoreItem(key uint32, length uint32, nextItemPageNum uint32, nextItemOffset uint16) StoreItem {
	return StoreItem{Key: key, NextItemPageNum: nextItemPageNum, NextItemOffset: nextItemOffset, Length: length}
}

func Deserialise(enc []byte) StoreItem {
	r := StoreItem{}
	r.Key = binary.LittleEndian.Uint32(enc[0:4])
	r.Length = binary.LittleEndian.Uint32(enc[4:8])
	r.NextItemPageNum = binary.LittleEndian.Uint32(enc[8:12])
	r.NextItemOffset = binary.LittleEndian.Uint16(enc[12:14])

	return r
}

func Serialise(r StoreItem) []byte {
	enc := make([]byte, 14)
	binary.LittleEndian.PutUint32(enc[0:4], r.Key)
	binary.LittleEndian.PutUint32(enc[4:8], r.Length)
	binary.LittleEndian.PutUint32(enc[8:12], uint32(r.NextItemPageNum))
	binary.LittleEndian.PutUint16(enc[12:14], uint16(r.NextItemOffset))

	return enc
}
