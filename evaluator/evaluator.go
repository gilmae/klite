package evaluator

import (
	"github.com/gilmae/klite/ast"
	"github.com/gilmae/klite/environment"
	"github.com/gilmae/klite/object"
)

func Eval(node ast.Node, env *environment.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.SelectStatement:
		// c := cursor.Start(env.Table)
		// rs := object.RecordSet{
		// 	Value:   [][]object.Object{},
		// 	Columns: []object.String{{Value: "Id"}, {Value: "Username"}, {Value: "Email"}},
		// }
		// for !c.Eof() {
		// 	fmt.Printf("%+v\n", c)
		// 	row, err := c.Row()
		// 	if err != nil {
		// 		return &object.Error{Message: fmt.Sprintf("%s", err)}
		// 	}

		// 	data, err := data.DeserialiseRow(*row)
		// 	if err != nil {
		// 		return &object.Error{Message: fmt.Sprintf("%s", err)}
		// 	}
		// 	r := make([]object.Object, 3)
		// 	r[0] = &object.Integer{Value: int64(data.Id)}
		// 	r[1] = &object.String{Value: string(data.Username[:])}
		// 	r[2] = &object.String{Value: string(data.Email[:])}
		// 	rs.Value = append(rs.Value, r)
		// 	c.Next()
		// }
		return nil
	case *ast.InsertStatement:
		// 	page, err := env.Table.Pager.Page(env.Table.RootPageNum)
		// 	if err != nil {
		// 		return &object.Error{Message: fmt.Sprintf("%s", err)}
		// 	}
		//leaf := data.Leaf{page}
		//numCells := leaf.GetNumCells()
		// if numCells >= uint32(data.LEAF_NODE_MAX_CELLS) {
		// 	return &object.Error{Message: "Maximum cells reached"}
		// }

		// rowToInsert := parser.Row{}
		// id, _ := strconv.Atoi(node.Arguments[0].String())
		// rowToInsert.Id = uint32(id)
		// username := []byte(node.Arguments[1].String())
		// copy(rowToInsert.Username[0:len(username)], username[:])

		// email := []byte(node.Arguments[2].String())
		// copy(rowToInsert.Email[0:len(email)], email[:])

		// rd, err := data.SerialiseRow(&rowToInsert)
		// if err != nil {
		// 	return &object.Error{Message: fmt.Sprintf("%s", err)}
		// }

		// c, err := cursor.Find(env.Table, rowToInsert.Id)
		// if err != nil {
		// 	return &object.Error{Message: fmt.Sprintf("%s", err)}
		// }

		// if !c.Eof() {
		// 	keyAtIndex, err := c.CurrentKey()
		// 	if err != nil {
		// 		return &object.Error{Message: fmt.Sprintf("%s", err)}
		// 	}

		// 	if keyAtIndex == rowToInsert.Id {
		// 		return &object.Error{Message: fmt.Sprintf("Duplicate key: %d", keyAtIndex)}
		// 	}
		// }

		// c.Insert(rowToInsert.Id, rd)

		// return &object.Integer{Value: int64(-1)} // TODO should return number of rows
		return nil

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
