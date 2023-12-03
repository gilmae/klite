package main

import (
	"fmt"
	"os"

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
}
