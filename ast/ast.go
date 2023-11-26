package ast

import (
	"bytes"

	"github.com/gilmae/klite/token"
)

type Node interface {
	String() string
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

type SelectStatement struct {
	Token  token.Token
	Key    Expression
	Num    Expression
	Stream Expression
}

func (ss *SelectStatement) statementNode()       {}
func (ss *SelectStatement) TokenLiteral() string { return ss.Token.Literal }
func (ss *SelectStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ss.TokenLiteral() + " ")
	if ss.Num != nil {
		out.WriteString(ss.Num.String() + " after")
	}

	out.WriteString(ss.Key.String())

	if ss.Stream != nil {
		out.WriteString(" from " + ss.Stream.String())
	}
	return out.String()
}

type InsertStatement struct {
	Token    token.Token
	Argument Expression
}

func (is *InsertStatement) statementNode()       {}
func (is *InsertStatement) TokenLiteral() string { return is.Token.Literal }
func (is *InsertStatement) String() string {
	var out bytes.Buffer

	out.WriteString("insert ")
	out.WriteString(is.Argument.String())
	return out.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) String() string {
	return i.Value
}

func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }
