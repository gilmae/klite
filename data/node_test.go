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

func TestIsRoot(t *testing.T) {
	tests := []struct {
		page          Page
		expectedValue bool
	}{
		{Page{0, 1}, true},
		{Page{1, 0}, false},
	}

	for _, test := range tests {
		n := NewNode(&test.page)
		if n.IsRoot() != test.expectedValue {
			t.Errorf("node is wrong type ,expected %+v, got %+v", test.expectedValue, n.IsRoot())
		}
	}
}

func TestSetIsRoot(t *testing.T) {
	tests := []struct {
		isRoot        bool
		expectedValue byte
	}{
		{false, 0x0},
		{true, 0x1},
	}

	for _, test := range tests {
		page := Page(make([]byte, PageSize))
		n := NewNode(&page)
		n.SetIsRoot(test.isRoot)

		if page[1:2][0] != test.expectedValue {
			t.Errorf("node is wrong type ,expected %+v, got %+v", test.expectedValue, page[1:2][0])

		}
	}
}
