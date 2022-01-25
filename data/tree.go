package data

type Tree struct {
	pager Pager
}

func NewTree(pager Pager) *Tree {
	return &Tree{pager: pager}
}
