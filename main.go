package main

import (
	"fmt"

	"github.com/gilmae/klite/data"
	"github.com/gilmae/klite/store"
)

func main() {
	// var pager data.Pager
	// _ = pager.GetNextUnusedPageNum()
	// pager, _ = data.NewFilePager(".test.db")
	// tree := data.NewTree(pager, pager.GetNextUnusedPageNum())
	// for i := uint32(0); i < 1000; i++ {
	// 	tree.Insert(i+1, data.NewIndexItem(i, 0, i+1))
	// }
	// (pager).Flush()
	// (pager).Close()

	// record := tree.Get(2)
	// fmt.Println(record)

	pager, _ := data.NewFilePager(".test.db")
	indexPageNum := pager.GetNextUnusedPageNum()
	indexPage, _ := pager.Page(indexPageNum)
	data.NewLeaf(indexPage)

	stream := store.NewStream(pager, 0, 0, 1, indexPageNum)

	for i := uint32(0); i < 1; i++ {
		stream.Add([]byte(fmt.Sprintf("Node %d", i)))
	}
	(pager).Flush()
	(pager).Close()
}
