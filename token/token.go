package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLIGAL TokenType = "ILLIGAL"
	EOF     TokenType = "EOF"

	// literals + identifiers
	IDENT  TokenType = "IDENT"
	INT    TokenType = "INT"
	STRING TokenType = "STRING"

	// keywords
	FUNCTION TokenType = "FUNCTION"
	VAR      TokenType = "VAR"
	RETURN   TokenType = "RETURN"
	IF       TokenType = "IF"
	ELSE     TokenType = "ELSE"
	TRUE     TokenType = "TRUE"
	FALSE    TokenType = "FALSE"

	// operators
	ASSIGN   TokenType = "="
	PLUS     TokenType = "+"
	MINUS    TokenType = "-"
	ASTERISK TokenType = "*"
	SLASH    TokenType = "/"
	BANG     TokenType = "!"
	GT       TokenType = ">"
	LT       TokenType = "<"
	EQUAL    TokenType = "=="
	GTEQUAL  TokenType = ">="
	LTEQUAL  TokenType = "<="
	NOTEQU   TokenType = "!="
	AND      TokenType = "&&"
	OR       TokenType = "||"

	// delimiters
	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"
	LPARENT   TokenType = "("
	RPARENT   TokenType = ")"
	LBRACE    TokenType = "{"
	RBRACE    TokenType = "}"
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
