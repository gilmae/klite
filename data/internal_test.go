package data

import "testing"

func TestGetNumKeys(t *testing.T) {

	p := Page{0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x2, 0x0}
	expectedNumKeys := 2
	node := Node{page: &p}

	actualNumKeys := node.NumKeys()

	if expectedNumKeys != int(actualNumKeys) {
		t.Errorf("incorrect number of keys, expected %d, got %d", expectedNumKeys, actualNumKeys)
	}

}

func TestSetNumKeys(t *testing.T) {
	p := Page(make([]byte, PageSize))
	node := NewInternal(&p)

	node.SetNumKeys(6)

	expectedData := []byte{0x6, 0x0}
	actualData := p[6:8]

	if !bytesMatch(expectedData, actualData) {
		t.Errorf("incorrect bytes for NumKeys set, expected %+v, got %+v", expectedData, actualData)
	}
}

func TestRightChild(t *testing.T) {
	p := Page{0x1, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1, 0x0, 0x1, 0x1, 0x0, 0x0}
	node := Node{page: &p}

	expectedResult := uint32(257)
	actualResult := node.RightChild()

	if expectedResult != actualResult {
		t.Errorf("incorrect value for node.RightChild, expected %d, got %d", expectedResult, actualResult)
	}

}

func TestSetRightChild(t *testing.T) {
	p := Page(make([]byte, PageSize))
	node := NewInternal(&p)

	node.SetRightChild(264)

	expectedData := []byte{0x8, 0x1, 0x0, 0x0}
	actualData := p[8:12]

	if !bytesMatch(expectedData, actualData) {
		t.Errorf("incorrect bytes for NumKeys set, expected %+v, got %+v", expectedData, actualData)
	}
}

func TestChildPointer(t *testing.T) {
	tests := []struct {
		cell           uint16
		data           []byte
		expectedResult uint32
	}{
		{0, []byte{0x0, 0x1, 0x0, 0x0}, 256},
		{1, []byte{0x1, 0x1, 0x1, 0x1}, 16843009},
	}

	for _, test := range tests {
		p := Page(make([]byte, PageSize))
		node := NewInternal(&p)
		//copy(p[NumCellsOffset:NumCellsOffset+1], []byte{byte(test.cell + 1)})
		offset := 12 + 8*int(test.cell)
		copy(p[offset:offset+4], test.data)
		actualResult := node.ChildPointer(test.cell)

		if test.expectedResult != actualResult {
			t.Errorf("Incorrect value for node.ChildPointer, exected %d, got %d", test.expectedResult, actualResult)
		}
	}

}

func TestSetChildPointer(t *testing.T) {
	tests := []struct {
		cell           uint16
		data           uint32
		expectedResult []byte
	}{
		{0, 256, []byte{0x0, 0x1, 0x0, 0x0}},
		{1, 16843008, []byte{0x0, 0x1, 0x1, 0x1}},
	}

	for _, test := range tests {
		p := Page(make([]byte, PageSize))
		node := NewInternal(&p)
		//copy(p[NumCellsOffset:NumCellsOffset+1], []byte{byte(test.cell + 1)})

		node.SetChildPointer(test.cell, test.data)
		offset := 12 + 8*int(test.cell)
		actualResult := p[offset : offset+4]

		if !bytesMatch(test.expectedResult, actualResult) {
			t.Errorf("Incorrect data set by node.SetChildPointer, exected %+v, got %+v", test.expectedResult, actualResult)
		}
	}

}

func TestInternalKey(t *testing.T) {
	tests := []struct {
		cell           uint16
		data           []byte
		expectedResult uint32
	}{
		{0, []byte{0x0, 0x1, 0x0, 0x0}, 256},
		{1, []byte{0x1, 0x1, 0x1, 0x1}, 16843009},
	}

	for _, test := range tests {
		p := Page(make([]byte, PageSize))
		node := NewInternal(&p)
		offset := 16 + 8*int(test.cell)
		copy(p[offset:offset+4], test.data)
		actualResult := node.InternalKey(test.cell)

		if test.expectedResult != actualResult {
			t.Errorf("Incorrect value for node.ChildPointer, exected %d, got %d", test.expectedResult, actualResult)
		}
	}

}

func TestSetTestInternalKey(t *testing.T) {
	tests := []struct {
		cell           uint16
		data           uint32
		expectedResult []byte
	}{
		{0, 257, []byte{0x1, 0x1, 0x0, 0x0}},
		{1, 16843010, []byte{0x2, 0x1, 0x1, 0x1}},
	}

	for _, test := range tests {
		p := Page(make([]byte, PageSize))
		node := NewInternal(&p)
		node.SetInternalKey(test.cell, test.data)
		offset := 16 + 8*int(test.cell)
		actualResult := p[offset : offset+4]

		if !bytesMatch(test.expectedResult, actualResult) {
			t.Errorf("Incorrect data set by node.SetChildPointer, exected %+v, got %+v", test.expectedResult, actualResult)
		}
	}

}
