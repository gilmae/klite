package store

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSerialise(t *testing.T) {
	r := StoreItem{Key: 1, NextItemPageNum: 3, NextItemOffset: 2, Length: 513}
	enc := Serialise(r)
	expectedValue := []byte{1, 0, 0, 0, 1, 2, 0, 0, 3, 0, 0, 0, 2, 0}
	if !cmp.Equal(enc, expectedValue) {
		t.Errorf("incorrect serialised value, expected %+v, got %+v", expectedValue, enc)
	}
}

func TestDeserialise(t *testing.T) {
	bytes := []byte{7, 0, 0, 0, 3, 1, 0, 0, 1, 0, 0, 0, 4, 0}
	r := Deserialise(bytes)

	expectedKey := uint32(7)
	expectedLength := uint32(259)
	expectedNextItemPageNum := uint32(1)
	expectedNextItemOffset := uint16(4)

	if r.Key != expectedKey {
		t.Errorf("incorrect Key, expected %d, got %d", expectedKey, r.Key)
	}

	if r.Length != expectedLength {
		t.Errorf("incorrect r.length, expected %d, got %d", expectedLength, r.Length)
	}

	if r.NextItemPageNum != expectedNextItemPageNum {
		t.Errorf("incorrect NextItemPageNum, expected %d, got %d", expectedNextItemPageNum, r.NextItemPageNum)
	}

	if uint32(r.NextItemOffset) != uint32(expectedNextItemOffset) {
		t.Errorf("incorrect NextItemOffset, expected %d, got %d", expectedNextItemOffset, r.NextItemOffset)
	}
}
