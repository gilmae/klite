package data

import "testing"

func TestSerialise(t *testing.T) {
	r := IndexItem{pageNum: 3, length: 513}
	enc := Serialise(r)
	expectedValue := []byte{3, 0, 0, 0, 0, 0, 1, 2, 0, 0}
	if !bytesMatch(enc, expectedValue) {
		t.Errorf("incorrect serialised value, expected %+v, got %+v", expectedValue, enc)
	}
}

func TestDeserialise(t *testing.T) {
	bytes := []byte{7, 0, 0, 0, 0, 0, 3, 1, 0, 0}
	r := Deserialise(bytes)

	expectedPageNum := uint32(7)
	expectedLength := uint32(259)

	if r.pageNum != expectedPageNum {
		t.Errorf("incorrect r.pageNum, expected %d, got %d", expectedPageNum, r.pageNum)
	}

	if r.length != expectedLength {
		t.Errorf("incorrect r.length, expected %d, got %d", expectedLength, r.length)
	}
}
