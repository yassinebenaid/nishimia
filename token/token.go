package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLIGAL = "ILLIGAL"
	EOF     = "EOF"

	// literals + identifiers
	IDENT = "IDENT"
	INT   = "INT"

	// keywords
	FUNCTION = "FUNCTION"
	VAR      = "FUNCTION"

	// operators
	ASSIGN = "="
	PLUS   = "+"
	MINUS  = "-"
	MULTP  = "*"
	DEVIDE = "/"

	// delimiters
	COMMA     = ","
	SEMICOLON = ";"
	LPARENT   = "("
	RPARENT   = ")"
	LBRACE    = "{"
	RBRACE    = "}"
)

var keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"var": VAR,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENT
}
