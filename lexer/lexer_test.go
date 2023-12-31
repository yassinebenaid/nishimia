package lexer

import (
	"testing"

	"github.com/yassinebenaid/nishimia/token"
)

func TestNextToken(t *testing.T) {
	input := `
	var five = 5;
	var ten = 10;
	
	var add = func(x, y) {
		return x + y;
	};

	var multiply = func(x, y) {
		return x * y;
	};

	var devide = func(x, y) {
		if(y > 0){
			return x / y;
		}else{
			return 0;
		}
	};

	var isPositive = func(x) {
		if(x >= 0){
			return true;
		}else{
			return false;
		}
	};

	var isZero = func(x) {
		return x == 0;
	};

	var isNotZero = func(x) {
		return x != 0;
	};

	var isNegativeOrZero = func(x) {
		return x <= 0;
	};

	var max = func(x, y) {
		if(x > y){
			return x ;
		}

		if(x < y){
			return y ;
		}

		return x;
	};
	
	var addition = add(five, ten);
	var multiplication = multiply(five, ten);
	var devision = devide(five, ten);
	var maximum = max(five, ten);
	var fiveIsPositive = isPositive(five);
	var fiveIsZero = isZero(five);
	var fiveIsNotZero = isNotZero(five);
	var tenIsNegativeOrZero = isNegativeOrZero(five);

	var name = "yassine benaid";
	var skipped = "yassine\" benaid";

	var myArray = ["yassinebenaid",10,-113];

	var myHash = {name: "yassine benaid", age: 21};
`

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
		{token.FUNCTION, "func"},
		{token.LPARENT, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPARENT, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},

		{token.VAR, "var"},
		{token.IDENT, "multiply"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "func"},
		{token.LPARENT, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPARENT, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENT, "x"},
		{token.ASTERISK, "*"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},

		{token.VAR, "var"},
		{token.IDENT, "devide"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "func"},
		{token.LPARENT, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPARENT, ")"},
		{token.LBRACE, "{"},
		{token.IF, "if"},
		{token.LPARENT, "("},
		{token.IDENT, "y"},
		{token.GT, ">"},
		{token.INT, "0"},
		{token.RPARENT, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENT, "x"},
		{token.SLASH, "/"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.INT, "0"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},

		{token.VAR, "var"},
		{token.IDENT, "isPositive"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "func"},
		{token.LPARENT, "("},
		{token.IDENT, "x"},
		{token.RPARENT, ")"},
		{token.LBRACE, "{"},
		{token.IF, "if"},
		{token.LPARENT, "("},
		{token.IDENT, "x"},
		{token.GTEQUAL, ">="},
		{token.INT, "0"},
		{token.RPARENT, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},

		{token.VAR, "var"},
		{token.IDENT, "isZero"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "func"},
		{token.LPARENT, "("},
		{token.IDENT, "x"},
		{token.RPARENT, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENT, "x"},
		{token.EQUAL, "=="},
		{token.INT, "0"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},

		{token.VAR, "var"},
		{token.IDENT, "isNotZero"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "func"},
		{token.LPARENT, "("},
		{token.IDENT, "x"},
		{token.RPARENT, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENT, "x"},
		{token.NOTEQU, "!="},
		{token.INT, "0"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},

		{token.VAR, "var"},
		{token.IDENT, "isNegativeOrZero"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "func"},
		{token.LPARENT, "("},
		{token.IDENT, "x"},
		{token.RPARENT, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENT, "x"},
		{token.LTEQUAL, "<="},
		{token.INT, "0"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},

		{token.VAR, "var"},
		{token.IDENT, "max"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "func"},
		{token.LPARENT, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPARENT, ")"},
		{token.LBRACE, "{"},
		{token.IF, "if"},
		{token.LPARENT, "("},
		{token.IDENT, "x"},
		{token.GT, ">"},
		{token.IDENT, "y"},
		{token.RPARENT, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENT, "x"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.IF, "if"},
		{token.LPARENT, "("},
		{token.IDENT, "x"},
		{token.LT, "<"},
		{token.IDENT, "y"},
		{token.RPARENT, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.RETURN, "return"},
		{token.IDENT, "x"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},

		{token.VAR, "var"},
		{token.IDENT, "addition"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPARENT, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPARENT, ")"},
		{token.SEMICOLON, ";"},

		{token.VAR, "var"},
		{token.IDENT, "multiplication"},
		{token.ASSIGN, "="},
		{token.IDENT, "multiply"},
		{token.LPARENT, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPARENT, ")"},
		{token.SEMICOLON, ";"},

		{token.VAR, "var"},
		{token.IDENT, "devision"},
		{token.ASSIGN, "="},
		{token.IDENT, "devide"},
		{token.LPARENT, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPARENT, ")"},
		{token.SEMICOLON, ";"},

		{token.VAR, "var"},
		{token.IDENT, "maximum"},
		{token.ASSIGN, "="},
		{token.IDENT, "max"},
		{token.LPARENT, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPARENT, ")"},
		{token.SEMICOLON, ";"},

		{token.VAR, "var"},
		{token.IDENT, "fiveIsPositive"},
		{token.ASSIGN, "="},
		{token.IDENT, "isPositive"},
		{token.LPARENT, "("},
		{token.IDENT, "five"},
		{token.RPARENT, ")"},
		{token.SEMICOLON, ";"},

		{token.VAR, "var"},
		{token.IDENT, "fiveIsZero"},
		{token.ASSIGN, "="},
		{token.IDENT, "isZero"},
		{token.LPARENT, "("},
		{token.IDENT, "five"},
		{token.RPARENT, ")"},
		{token.SEMICOLON, ";"},

		{token.VAR, "var"},
		{token.IDENT, "fiveIsNotZero"},
		{token.ASSIGN, "="},
		{token.IDENT, "isNotZero"},
		{token.LPARENT, "("},
		{token.IDENT, "five"},
		{token.RPARENT, ")"},
		{token.SEMICOLON, ";"},

		{token.VAR, "var"},
		{token.IDENT, "tenIsNegativeOrZero"},
		{token.ASSIGN, "="},
		{token.IDENT, "isNegativeOrZero"},
		{token.LPARENT, "("},
		{token.IDENT, "five"},
		{token.RPARENT, ")"},
		{token.SEMICOLON, ";"},

		{token.VAR, "var"},
		{token.IDENT, "name"},
		{token.ASSIGN, "="},
		{token.STRING, "yassine benaid"},
		{token.SEMICOLON, ";"},

		{token.VAR, "var"},
		{token.IDENT, "skipped"},
		{token.ASSIGN, "="},
		{token.STRING, "yassine\" benaid"},
		{token.SEMICOLON, ";"},

		{token.VAR, "var"},
		{token.IDENT, "myArray"},
		{token.ASSIGN, "="},
		{token.LBRACKET, "["},
		{token.STRING, "yassinebenaid"},
		{token.COMMA, ","},
		{token.INT, "10"},
		{token.COMMA, ","},
		{token.MINUS, "-"},
		{token.INT, "113"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},

		{token.VAR, "var"},
		{token.IDENT, "myHash"},
		{token.ASSIGN, "="},
		{token.LBRACE, "{"},
		{token.IDENT, "name"},
		{token.COLON, ":"},
		{token.STRING, "yassine benaid"},
		{token.COMMA, ","},
		{token.IDENT, "age"},
		{token.COLON, ":"},
		{token.INT, "21"},
		{token.RBRACE, "}"},
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
