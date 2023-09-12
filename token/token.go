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
	VAR      = "VAR"
	RETURN   = "RETURN"
	IF       = "IF"
	ELSE     = "ELSE"
	TRUE     = "TRUE"
	FALSE    = "FALSE"

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
	"func":   FUNCTION,
	"var":    VAR,
	"return": RETURN,
	"if":     IF,
	"else":   ELSE,
	"true":   TRUE,
	"false":  FALSE,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENT
}
