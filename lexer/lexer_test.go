package lexer

import (
	"testing"

	"github.com/gilmae/klite/token"
)

func TestNextToken(t *testing.T) {
	input := `select;
	insert (1, 'a', 'b')`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.SELECT, "select"},
		{token.SEMICOLON, ";"},
		{token.INSERT, "insert"},
		{token.LPAREN, "("},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.STRING, "a"},
		{token.COMMA, ","},
		{token.STRING, "b"},
		{token.RPAREN, ")"},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokenType wrong, expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong, expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
