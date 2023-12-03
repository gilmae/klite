package repl

import (
	"bufio"
	"fmt"

	"io"
	"os"

	"github.com/gilmae/klite/data"
	"github.com/gilmae/klite/environment"
	"github.com/gilmae/klite/evaluator"
	"github.com/gilmae/klite/lexer"
	"github.com/gilmae/klite/parser"
)

const PROMPT = ">> "

const (
	_ int = iota
	META_COMMAND_SUCCESS
	META_COMMAND_UNRECOGNISED_COMMAND
)

func Start(dbPath string, in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(os.Stdin)
	pager, err := data.NewFilePager(dbPath)

	if err != nil {
		fmt.Printf("Error opening table: %s", err)
	}
	defer pager.Close()

	env, err := environment.NewEnvironment(pager)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if !env.IsInitialised() {
		err = env.Initialise()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		if line == ".exit" {
			break
		} else if line[0] == '.' {
			switch doMetaCommand(line, env) {
			case META_COMMAND_SUCCESS:
				continue
			case META_COMMAND_UNRECOGNISED_COMMAND:
				fmt.Printf("Unrecognised meta command '%s'.\n", line)
				continue
			}
		}

		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()

		result := evaluator.Eval(program, env)
		if result != nil {
			fmt.Printf("%s\n", result.Inspect())
		}
		continue

	}
}

func doMetaCommand(line string, env *environment.Environment) int {
	// .exit is handled outside to make breaking out of the repl easier
	// We'll add to this when there are more meta commands to handle
	switch line {
	case ".peek":
		// for i := uint32(0); i < pager.NumPages; i++ {
		// 	p, err := env.Table.Pager.Page(i)
		// 	if err != nil {
		// 		fmt.Printf("%s", err)
		// 	}
		// 	fmt.Printf("Page %d is a NodeType %s\n", i, data.GetNodeType(p))
		// }
		return 0
	case ".stream":
		s := env.GetStream()
		fmt.Printf("Index Root Page\t\t: %d\n", s.IndexPage())
		fmt.Printf("Store Head Page\t\t: %d\n", s.StoreHeadPage())
		fmt.Printf("Store Tail Page\t\t: %d\n", s.StoreTailPage())
		fmt.Printf("Next Key\t\t: %d\n", s.NextKey())

		indexPage, err := env.Pager().Page(s.IndexPage())
		if err != nil {
			fmt.Println(err)
			return -1
		}
		index := data.NewNode(indexPage)
		maxKey, err := index.GetMaxKey()
		if err != nil {
			fmt.Printf("Max Key in Index\t: %s\n", err)
		} else {
			fmt.Printf("Max Key in Index\t: %d\n", maxKey)
		}

		return 0
	}
	return META_COMMAND_UNRECOGNISED_COMMAND
}
