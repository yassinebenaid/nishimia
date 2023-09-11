package token

import "fmt"

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
	RETURN   = "RETURN"
	IF       = "IF"
	ELSE     = "ELSE"

	// operators
	ASSIGN = "="
	PLUS   = "+"
	MINUS  = "-"
	MULTP  = "*"
	SLASH  = "/"
	BANG   = "!"
	GT     = ">"
	LT     = "<"

	// delimiters
	COMMA     = ","
	SEMICOLON = ";"
	LPARENT   = "("
	RPARENT   = ")"
	LBRACE    = "{"
	RBRACE    = "}"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"var":    VAR,
	"return": RETURN,
	"if":     IF,
	"else":   ELSE,
}

func LookupIdent(ident string) TokenType {
	fmt.Println(ident)
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENT
}
