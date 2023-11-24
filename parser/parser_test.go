package parser

import (
	"fmt"
	"testing"

	"github.com/gilmae/sqlike/ast"
	"github.com/gilmae/sqlike/lexer"
)

func TestSelectStatement(t *testing.T) {
	input := `
	select
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 expressions, got %d",
			len(program.Statements))

	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.SelectStatement)
		if !ok {
			t.Errorf("stmt not *ast.SelectStatement. got %T", stmt)
		}

		if returnStmt.TokenLiteral() != "select" {
			t.Errorf("returnStmt.TokenLiteral() not 'select', got %q", returnStmt.TokenLiteral())
		}
	}
}

func TestInsertStatement(t *testing.T) {
	input := "insert (1, 'a');"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Errorf("program does not have enough statements. Expected %d, got %d", 1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.InsertStatement)
	if !ok {
		t.Errorf(" program.Statements[0] not *ast.InsertStatement. got %T", stmt)
	}

	if len(stmt.Arguments) != 2 {
		t.Fatalf("stmt.Arguments is wrong, expected %d, got %d",
			2,
			len(stmt.Arguments))
	}

	testLiteralExpression(t, stmt.Arguments[0], 1)
	testLiteralExpression(t, stmt.Arguments[1], "a")
}

// func TestInsertStatement(t *testing.T) {
// 	input := []string{"insert 1 a b into topic1"}

// 	for _, txt := range input {
// 		s, err := PrepareStatement(txt)
// 		if err != nil {
// 			t.Fatalf("statement could not be parsed as a statement: '%s'", txt)
// 		}
// 		ss, ok := s.(*InsertStatement)
// 		if !ok {
// 			t.Errorf("expected InsertStatement, got %T (%+v)", ss, ss)
// 		}

// 		if ss.RowToInsert.Id != 1 {
// 			t.Errorf("expected RowToInsert.Id to be %d, got %d", 1, ss.RowToInsert.Id)
// 		}

// 		if ss.RowToInsert.Username != "a" {
// 			t.Errorf("expected RowToInsert.Username to be %s, got %s", "a", ss.RowToInsert.Username)
// 		}

// 		if ss.RowToInsert.Email != "b" {
// 			t.Errorf("expected RowToInsert.Email to be %s, got %s", "b", ss.RowToInsert.Email)
// 		}
// 	}
// }
// func TestSelectStatement(t *testing.T) {
// 	input := []string{"select", "select 1", "select 1 from topic1"}

// 	for _, txt := range input {
// 		s, err := PrepareStatement(txt)
// 		if err != nil {
// 			t.Fatalf("statement could not be parsed as a statement: '%s'", txt)
// 		}
// 		ss, ok := s.(*SelectStatement)
// 		if !ok {
// 			t.Errorf("expected SelectStatement, got %T (%+v)", ss, ss)
// 		}
// 	}
// }

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral, got %T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d, got %d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d, got %s", value, integ.TokenLiteral())
		return false
	}

	return true
}

func testStringLiteral(t *testing.T, il ast.Expression, value string) bool {
	str, ok := il.(*ast.StringLiteral)
	if !ok {
		t.Errorf("il not *ast.StringLiteral, got %T", il)
		return false
	}

	if str.Value != value {
		t.Errorf("str.Value not %s, got %s", value, str.Value)
		return false
	}

	if str.TokenLiteral() != fmt.Sprintf("%s", value) {
		t.Errorf("str.TokenLiteral not %s, got %s", value, str.TokenLiteral())
		return false
	}

	return true
}
func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value,
			ident.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(
	t *testing.T,
	exp ast.Expression,
	expected interface{},
) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testStringLiteral(t, exp, v)
		// case bool:
		// 	return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser as %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("Parser error: %q", msg)
	}

	t.FailNow()
}
