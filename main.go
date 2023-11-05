package main

import (
	"fmt"

	"github.com/gilmae/klite/data"
)

func main() {
	var pager data.Pager
	_ = pager.GetNextUnusedPageNum()
	pager, _ = data.NewFilePager(".test.db")
	tree := data.NewTree(pager, pager.GetNextUnusedPageNum())
	for i := uint32(0); i < 1000; i++ {
		tree.Insert(i+1, data.NewIndexItem(i, 0, i+1))
	}
	(pager).Flush()
	(pager).Close()

	record := tree.Get(2)
	fmt.Println(record)
}
