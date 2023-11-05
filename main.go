package main

import (
	"fmt"

	"github.com/gilmae/data/data"
)

func main() {
	var pager data.Pager
	pager, _ = data.NewFilePager(".test.db")
	tree := data.NewTree(pager)
	for i := uint32(0); i < 1000; i++ {
		tree.Insert(i+1, data.NewIndexItem(i, i+1))
	}
	(pager).Flush()
	(pager).Close()

	record := tree.Get(2)
	fmt.Println(record)
}
