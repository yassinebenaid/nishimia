package lexer

import (
	"testing"

	"github.com/yassinebenaid/nishimia/token"
)

func TestNextToken(t *testing.T) {
	input := `
	var five = 5;
	var ten = 10;
	
	var add = fn(x, y) {
		x + y;
	};
	
	var result = add(five, ten);`

	cases := []struct {
		tokenType    token.TokenType
		tokenLiteral string
	}{
		{token.VAR, "var"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.VAR, "var"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		{token.VAR, "var"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPARENT, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPARENT, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},

		{token.VAR, "var"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPARENT, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPARENT, ")"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, cas := range cases {
		tok := l.NextToken()

		if tok.Type != cas.tokenType {
			t.Fatalf("test #%d failed, expected type [%s] but got [%s]", i, cas.tokenType, tok.Type)
		}

		if tok.Literal != cas.tokenLiteral {
			t.Fatalf("test #%d failed, expected literal [%s] but got [%s]", i, cas.tokenLiteral, tok.Literal)
		}
	}
}
