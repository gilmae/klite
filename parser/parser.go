package parser

import (
	"fmt"
	"strconv"

	"github.com/gilmae/klite/token"

	"github.com/gilmae/klite/ast"
	"github.com/gilmae/klite/lexer"
)

const (
	_ int = iota
	LOWEST
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
}

type (
	prefixParseFn func() ast.Expression
)

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	p.nextToken()
	p.nextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)

	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseSelectStatement() *ast.SelectStatement {
	stmt := &ast.SelectStatement{Token: p.curToken}

	for p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseInsertStatement() *ast.InsertStatement {
	stmt := &ast.InsertStatement{Token: p.curToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	stmt.Arguments = p.parseExpressionList(token.RPAREN)

	return stmt
}

func (p *Parser) parseExpression(precendence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]

	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()

	// for !p.peekTokenIs(token.SEMICOLON) && precendence < p.peekPrecedence() {
	// 	infix := p.infixParseFns[p.peekToken.Type]
	// 	if infix == nil {
	// 		return leftExp
	// 	}

	// 	p.nextToken()

	// 	leftExp = infix(leftExp)
	// }

	return leftExp
}

func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.SELECT:
		return p.parseSelectStatement()
	case token.INSERT:
		return p.parseInsertStatement()
	default:
		return nil
	}
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t,
		p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}
