package evaluator

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gilmae/klite/ast"
	"github.com/gilmae/klite/environment"
	"github.com/gilmae/klite/object"
)

func Eval(node ast.Node, env *environment.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.SelectStatement:

		stream := env.GetStream()
		key, err := strconv.Atoi(node.Key.String())
		if err != nil {
			return &object.Error{Message: fmt.Sprintf("%s", err)}
		}
		if node.Num == nil {
			value, err := stream.Get(uint32(key))
			if err != nil {
				return &object.Error{Message: fmt.Sprintf("%s", err)}
			}

			return &object.String{Value: fmt.Sprintf("%d:\t%s", value.Key, string(value.Data))}
		} else {
			num, err := strconv.Atoi(node.Num.String())
			if err != nil {
				return &object.Error{Message: fmt.Sprintf("%s", err)}
			}
			values, err := stream.GetFrom(uint32(key), uint16(num))
			if err != nil {
				return &object.Error{Message: fmt.Sprintf("%s", err)}
			}
			lines := make([]string, len(values))
			for idx, v := range values {
				lines[idx] = fmt.Sprintf("%d:\t%s", v.Key, string(v.Data))
			}
			return &object.String{Value: strings.Join(lines, "\n")}
		}

		return nil
	case *ast.InsertStatement:
		stream := env.GetStream()
		key, err := stream.Add([]byte(node.Argument.String()))
		if err != nil {
			return &object.Error{Message: fmt.Sprintf("%s", err)}
		}
		return &object.Integer{Value: int64(key)} // TODO should return number of rows
	}
	return &object.Null{}
}

func evalProgram(program *ast.Program, env *environment.Environment) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement, env)
	}

	return result
}
