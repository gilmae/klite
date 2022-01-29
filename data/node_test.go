package data

import "testing"

func TestNodeType(t *testing.T) {
	tests := []struct {
		page         Page
		expectedType NodeType
	}{
		{Page{0, 1}, LeafNode},
		{Page{1, 0}, InternalNode},
	}

	for _, test := range tests {
		n := NewNode(&test.page)
		if n.Type() != test.expectedType {
			t.Errorf("node is wrong type ,expected %s, got %s", test.expectedType, n.Type())
		}
	}
}

func TestSetNodeType(t *testing.T) {
	tests := []struct {
		page         Page
		expectedType NodeType
	}{
		{Page{0x1, 0x0}, LeafNode},
		{Page{0x0, 0x0}, InternalNode},
	}

	for _, test := range tests {
		n := NewNode(&test.page)
		n.SetType(test.expectedType)

		if n.Type() != test.expectedType {
			t.Errorf("node is wrong type ,expected %s, got %s", test.expectedType, n.Type())

		}
	}

}
