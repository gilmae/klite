package main

import (
	"fmt"
	"os"

	"github.com/gilmae/klite/data"
	"github.com/gilmae/klite/repl"
)

func main() {

	argv := os.Args
	if len(argv) < 2 {
		fmt.Println("Missing database")
		os.Exit(1)
	}
	fmt.Println(argv[1])
	repl.Start(argv[1], os.Stdin, os.Stdout)

	pager, _ := data.NewFilePager(".test.db")

	// env, err := environment.NewEnvironment(pager)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// if !env.IsInitialised() {
	// 	err = env.Initialise()
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		os.Exit(1)
	// 	}
	// }

	// stream := env.GetStream()
	// for i := uint32(0); i < 10; i++ {
	// 	stream.Add([]byte(fmt.Sprintf("Node %d", i)))
	// }

	// item, _ := stream.Get(2)
	// fmt.Println(string(item))

	pager.Flush()
	pager.Close()
}
