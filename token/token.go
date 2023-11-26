package token

type TokenType string
type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	SELECT = "SELECT"
	INSERT = "INSERT"

	SEMICOLON = "SEMICOLON"
	COMMA     = "COMMA"
	LPAREN    = "LPAREN"
	RPAREN    = "RPAREN"

	INT    = "INT"
	STRING = "STRING"

	NIL = "NIL"

	AFTER = "AFTER"
)

var keywords = map[string]TokenType{
	"add":   INSERT,
	"get":   SELECT,
	"after": AFTER,
}

// LookupIdent checks if an identifier is a keyword or a user identifier
func LookupIdent(ident string) TokenType {
	if tok, found := keywords[ident]; found {
		return tok
	}
	return ILLEGAL
}
