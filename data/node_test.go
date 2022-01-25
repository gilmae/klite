package data

import "testing"

func TestNodeType(t *testing.T) {
	tests := []struct {
		page         [PageSize]byte
		expectedType NodeType
	}{
		{[PageSize]byte{0, 1}, LeafNode},
		{[PageSize]byte{1, 0}, InternalNode},
	}

	for _, test := range tests {
		n := NewNode(test.page)
		if n.Type() != test.expectedType {
			t.Errorf("node is wrong type ,expected %s, got %s", test.expectedType, n.Type())
		}
	}
}

func TestSetNodeType(t *testing.T) {
	tests := []struct {
		page         [PageSize]byte
		expectedType NodeType
	}{
		{[PageSize]byte{0x1, 0x0}, LeafNode},
		{[PageSize]byte{0x0, 0x0}, InternalNode},
	}

	for _, test := range tests {
		n := NewNode(test.page)
		n.SetType(test.expectedType)

		if n.Type() != test.expectedType {
			t.Errorf("node is wrong type ,expected %s, got %s", test.expectedType, n.Type())

		}
	}

}
