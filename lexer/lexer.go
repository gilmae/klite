package lexer

import "github.com/gilmae/klite/token"

// Lexer represents the internal state of a lexer parsing through input
type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar() // initialise by reading the first char
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '\'':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		}
		if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		}
		tok = newToken(token.ILLEGAL, l.ch)
	}

	l.readChar()
	return tok
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition

	l.readPosition++
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]

}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '\'' || l.ch == 0 {
			break
		}
	}

	str := l.input[position:l.position]
	if len(str) == 0 {
		return str
	}
	ret := ""
	position = 0
	for {
		ch := str[position]
		if ch == '\\' {
			switch str[position+1] {
			case 't':
				ret = ret + string('\t')
				position += 1
			case 'n':
				ret = ret + string('\n')
				position += 1
			case '\\':
				ret = ret + string('\\')
				position += 1
			default:
				ret = ret + string('\\')
			}
		} else {
			ret = ret + string(ch)
		}
		position = position + 1
		if position >= len(str) {
			break
		}

	}

	return ret
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}
