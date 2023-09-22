package parser

import (
	"fmt"
	"testing"

	"github.com/yassinebenaid/nishimia/ast"
	"github.com/yassinebenaid/nishimia/lexer"
	"github.com/yassinebenaid/nishimia/token"
)

func TestVarStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      any
	}{
		{"var x = 5;", "x", 5},
		{"var y = true;", "y", true},
		{"var foobar = y;", "foobar", "y"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testVarStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.VarStatement).Value

		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func testVarStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "var" {
		t.Errorf("s.TokenLiteral isn't var, got=%s", s.TokenLiteral())
		return false
	}

	varStat, ok := s.(*ast.VarStatement)

	if !ok {
		t.Errorf("s is not VarStatement , got=%T", s)
		return false
	}

	if varStat.Name.Value != name {
		t.Errorf("varStat.Name.Value is not '%s' , got=%s", name, varStat.Name.Value)
		return false
	}

	if varStat.Name.TokenLiteral() != name {
		t.Errorf("varStat.Name is not '%s' , got=%s", name, varStat.Name.TokenLiteral())
		return false
	}

	return true
}

func TestReturnStatementParser(t *testing.T) {
	tests := []struct {
		input              string
		expectedExpression string
	}{
		{`return 5;`, "5"},
		{`return 5 * 5;`, "(5 * 5)"},
		{`return 5 == 7;`, "(5 == 7)"},
		{`return add(4,5);`, "add(4, 5)"},
		{`return add(4,min(5,6));`, "add(4, min(5, 6))"},
	}

	for _, test := range tests {
		lex := lexer.New(test.input)
		par := New(lex)
		program := par.ParseProgram()
		checkParserErrors(t, par)

		if len(program.Statements) != 1 {
			t.Fatalf("expected 1 statemnts, got %d", len(program.Statements))
		}

		ret, ok := program.Statements[0].(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("the program.Statements[0] type is not ast.ReturnStatement, got %T", program.Statements[0])
		}

		if ret.Token.Literal != "return" {
			t.Fatalf("statement.Token.Literal is not 'return', got %s", ret.Token.Literal)
		}

		if ret.TokenLiteral() != "return" {
			t.Fatalf("statement.Token.Literal is not 'return', got %s", ret.TokenLiteral())
		}

		if test.expectedExpression != ret.Return.String() {
			t.Fatalf("expected expression to be %s , got %s", test.expectedExpression, ret.Return.String())
		}
	}

}

func TestStringOnProgram(t *testing.T) {
	proram := ast.Program{
		Statements: []ast.Statement{
			&ast.VarStatement{
				Token: token.Token{Type: token.VAR, Literal: "var"},
				Name: &ast.Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "year"},
					Value: "year",
				},
				Value: &ast.Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "2023"},
					Value: "2023",
				},
			},
		},
	}

	if proram.String() != "var year = 2023;" {
		t.Fatalf("expected program string to be %s, got %s", "var year = 2023;", proram.String())
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := `foobar;`
	lex := lexer.New(input)
	par := New(lex)
	program := par.ParseProgram()

	checkParserErrors(t, par)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}

	stat, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("statement type is incorrect, expected ExpressionStatement, got=%T", program.Statements[0])
	}

	if !testIdentifierLiteral(t, stat.Expression, "foobar") {
		return
	}
}

func TestBooleanExpression(t *testing.T) {
	input := `
	true;
	false;
	`
	lex := lexer.New(input)
	par := New(lex)
	program := par.ParseProgram()

	checkParserErrors(t, par)

	if len(program.Statements) != 2 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}

	stat, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("statement type is incorrect, expected ExpressionStatement, got=%T", program.Statements[0])
	}

	if !testLiteralExpression(t, stat.Expression, true) {
		return
	}

	stat, ok = program.Statements[1].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("statement type is incorrect, expected ExpressionStatement, got=%T", program.Statements[1])
	}

	if !testLiteralExpression(t, stat.Expression, false) {
		return
	}

}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	lex := lexer.New(input)
	par := New(lex)
	program := par.ParseProgram()

	checkParserErrors(t, par)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}

	stat, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("statement type is incorrect, expected ExpressionStatement, got=%T", program.Statements[0])
	}

	if !testLiteralExpression(t, stat.Expression, 5) {
		return
	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`
	lex := lexer.New(input)
	par := New(lex)
	program := par.ParseProgram()

	checkParserErrors(t, par)

	if len(program.Statements) != 1 {
		t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
	}

	stat, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("statement type is incorrect, expected ExpressionStatement, got=%T", program.Statements[0])
	}

	testStringLiteral(t, stat.Expression, "hello world")
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input      string
		operator   string
		rightValue any
	}{
		{"!5;", "!", 5},
		{"-5;", "-", 5},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		par := New(l)
		program := par.ParseProgram()
		checkParserErrors(t, par)

		if len(program.Statements) != 1 {
			t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
		}

		stat, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("statement type is incorrect, expected ExpressionStatement, got=%T", program.Statements[0])
		}

		if !testPrefixExpression(t, stat.Expression, tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  any
		operator   string
		rightValue any
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"5 >= 5;", 5, ">=", 5},
		{"5 <= 5;", 5, "<=", 5},
		{"true == false;", true, "==", false},
		{"false != true;", false, "!=", true},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		par := New(l)
		program := par.ParseProgram()
		checkParserErrors(t, par)

		if len(program.Statements) != 1 {
			t.Fatalf("expected 1 statement, got=%d", len(program.Statements))
		}

		stat, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("statement type is incorrect, expected ExpressionStatement, got=%T", program.Statements[0])
		}

		if !testInfixExpression(t, stat.Expression, tt.leftValue, tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"a + b + c + d * 5", "(((a + b) + c) + (d * 5))"},
		{"-a * b", "((-a) * b)"},
		{"+a * b", "((+a) * b)"},
		{"1 + +a * b", "(1 + ((+a) * b))"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a <= b + c", "(a <= (b + c))"},
		{"a >= b + c", "(a >= (b + c))"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a / b / 2", "((a / b) / 2)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 < 4 != 3 > 4", "((5 < 4) != (3 > 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},

		{"true == !false", "(true == (!false))"},
		{"!true == !false", "((!true) == (!false))"},
		{"!true != false", "((!true) != false)"},
		{"!true != !false", "((!true) != (!false))"},
		{"!true != !false", "((!true) != (!false))"},

		{"true && false", "(true && false)"},
		{"true || false", "(true || false)"},
		{"true || false && true", "((true || false) && true)"},
		{"true && false && true || false", "(((true && false) && true) || false)"},

		{"(a + b) * 2", "((a + b) * 2)"},
		{"(a + b) / 2", "((a + b) / 2)"},
		{"!(false == true)", "(!(false == true))"},

		{"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))", "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))"},
		{"add(a + b + c * d / f + g)", "add((((a + b) + ((c * d) / f)) + g))"},

		{"add(a[0] + b[2] + c[3] * d[4] / f[5] + g[6])", "add((((a[0] + b[2]) + ((c[3] * d[4]) / f[5])) + g[6]))"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		par := New(l)
		program := par.ParseProgram()
		checkParserErrors(t, par)

		actual := program.String()

		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestIfExpressionParsing(t *testing.T) {
	input := `if x < y {x}`

	l := lexer.New(input)
	par := New(l)
	program := par.ParseProgram()
	checkParserErrors(t, par)

	if len(program.Statements) != 1 {
		t.Fatalf("expected statements count to be 1, got=%d", len(program.Statements))
	}

	stat, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected statement type of ExpressionStatement, got=%T", program.Statements[0])
	}

	ifelse, ok := stat.Expression.(*ast.IfElseExpression)
	if !ok {
		t.Fatalf("expected statement type of IfElseExpression, got=%T", stat)
	}

	if !testInfixExpression(t, ifelse.Condition, "x", "<", "y") {
		return
	}

	if len(ifelse.Consequence.Statements) != 1 {
		t.Fatalf("expected ifelse.Consequence.Statements count to be 1, got=%d", len(ifelse.Consequence.Statements))
	}

	seque, ok := ifelse.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected sequence first statement type of ExpressionStatement, got=%T", ifelse.Consequence.Statements[0])
	}

	if !testIdentifierLiteral(t, seque.Expression, "x") {
		return
	}

	if ifelse.Alternative != nil {
		t.Fatalf("ifels.Alternative is not nil, got %+v", ifelse.Alternative)
	}
}

func TestIfElseExpressionParsing(t *testing.T) {
	input := `if x < y { x } else { y }`

	l := lexer.New(input)
	par := New(l)
	program := par.ParseProgram()
	checkParserErrors(t, par)

	if len(program.Statements) != 1 {
		t.Fatalf("expected statements count to be 1, got=%d", len(program.Statements))
	}

	stat, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected statement type of ExpressionStatement, got=%T", program.Statements[0])
	}

	ifelse, ok := stat.Expression.(*ast.IfElseExpression)
	if !ok {
		t.Fatalf("expected statement type of IfElseExpression, got=%T", stat)
	}

	if !testInfixExpression(t, ifelse.Condition, "x", "<", "y") {
		return
	}

	if len(ifelse.Consequence.Statements) != 1 {
		t.Fatalf("expected ifelse.Consequence.Statements count to be 1, got=%d", len(ifelse.Consequence.Statements))
	}

	seque, ok := ifelse.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected sequence first statement type of ExpressionStatement, got=%T", ifelse.Consequence.Statements[0])
	}

	if !testIdentifierLiteral(t, seque.Expression, "x") {
		return
	}

	if len(ifelse.Alternative.Statements) != 1 {
		t.Fatalf("expected ifelse.Alternative.Statements count to be 1, got=%d", len(ifelse.Alternative.Statements))
	}

	alt, ok := ifelse.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected sequence first statement type of ExpressionStatement, got=%T", ifelse.Alternative.Statements[0])
	}

	if !testIdentifierLiteral(t, alt.Expression, "y") {
		return
	}
}

func TestFunctionLiteralWithoutParamsParsing(t *testing.T) {
	input := `
		func(){
			return x + y;
		}`

	l := lexer.New(input)
	par := New(l)
	program := par.ParseProgram()
	checkParserErrors(t, par)

	if len(program.Statements) != 1 {
		t.Fatalf("expected statements count to be 1, got=%d", len(program.Statements))
	}

	stat, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected statement type of ExpressionStatement, got=%T", program.Statements[0])
	}

	funct, ok := stat.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("expected statement type of FunctionLiteral, got=%T", stat)
	}

	if len(funct.Params) != 0 {
		t.Fatalf("expected arguments count to be 0, got=%d", len(funct.Params))
	}

	if len(funct.Body.Statements) != 1 {
		t.Fatalf("expected funct.Body.Statements count to be 1, got=%d", len(funct.Body.Statements))
	}

	ret, ok := funct.Body.Statements[0].(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("expected first statement type of ReturnStatement, got=%T", funct.Body.Statements[0])
	}

	if !testInfixExpression(t, ret.Return, "x", "+", "y") {
		return
	}

}

func TestFunctionLiteralWithOneParamParsing(t *testing.T) {
	input := `
		func(x){
			return x * 2;
		}`

	l := lexer.New(input)
	par := New(l)
	program := par.ParseProgram()
	checkParserErrors(t, par)

	if len(program.Statements) != 1 {
		t.Fatalf("expected statements count to be 1, got=%d", len(program.Statements))
	}

	stat, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected statement type of ExpressionStatement, got=%T", program.Statements[0])
	}

	funct, ok := stat.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("expected statement type of FunctionLiteral, got=%T", stat)
	}

	if len(funct.Params) != 1 {
		t.Fatalf("expected arguments count to be 2, got=%d", len(funct.Params))
	}

	if !testIdentifierLiteral(t, funct.Params[0], "x") {
		return
	}

	if len(funct.Body.Statements) != 1 {
		t.Fatalf("expected funct.Body.Statements count to be 1, got=%d", len(funct.Body.Statements))
	}

	ret, ok := funct.Body.Statements[0].(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("expected first statement type of ReturnStatement, got=%T", funct.Body.Statements[0])
	}

	if !testInfixExpression(t, ret.Return, "x", "*", 2) {
		return
	}

}

func TestFunctionLiteralWithManyParamsParsing(t *testing.T) {
	input := `
		func(x,y){
			return x + y;
		}`

	l := lexer.New(input)
	par := New(l)
	program := par.ParseProgram()
	checkParserErrors(t, par)

	if len(program.Statements) != 1 {
		t.Fatalf("expected statements count to be 1, got=%d", len(program.Statements))
	}

	stat, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected statement type of ExpressionStatement, got=%T", program.Statements[0])
	}

	funct, ok := stat.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("expected statement type of FunctionLiteral, got=%T", stat)
	}

	if len(funct.Params) != 2 {
		t.Fatalf("expected arguments count to be 2, got=%d", len(funct.Params))
	}

	if !testIdentifierLiteral(t, funct.Params[0], "x") && !testIdentifierLiteral(t, funct.Params[0], "y") {
		return
	}

	if len(funct.Body.Statements) != 1 {
		t.Fatalf("expected funct.Body.Statements count to be 1, got=%d", len(funct.Body.Statements))
	}

	ret, ok := funct.Body.Statements[0].(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("expected first statement type of ReturnStatement, got=%T", funct.Body.Statements[0])
	}

	if !testInfixExpression(t, ret.Return, "x", "+", "y") {
		return
	}
}

func TestFunctionCallExpressionParsing(t *testing.T) {
	input := `add(2 ,4 + y, 4 * 2);`

	l := lexer.New(input)
	par := New(l)
	program := par.ParseProgram()
	checkParserErrors(t, par)

	if len(program.Statements) != 1 {
		t.Fatalf("expected statements count to be 1, got=%d", len(program.Statements))
	}

	stat, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected statement type of ExpressionStatement, got=%T", program.Statements[0])
	}

	funct, ok := stat.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("expected statement type of CallExpression, got=%T", stat)
	}

	if !testLiteralExpression(t, funct.Function, "add") {
		return
	}

	if len(funct.Arguments) != 3 {
		t.Fatalf("expected arguments count to be 2, got=%d", len(funct.Arguments))
	}

	testLiteralExpression(t, funct.Arguments[0], 2)
	testInfixExpression(t, funct.Arguments[1], 4, "+", "y")
	testInfixExpression(t, funct.Arguments[2], 4, "*", 2)

}

func TestFunctionCallArgumentParsing(t *testing.T) {
	input := `add(-2 ,4 + y, 4 * -2 / -4, !true);`

	tests := []string{
		"(-2)",
		"(4 + y)",
		"((4 * (-2)) / (-4))",
		"(!true)",
	}

	l := lexer.New(input)
	par := New(l)
	program := par.ParseProgram()
	checkParserErrors(t, par)

	if len(program.Statements) != 1 {
		t.Fatalf("expected statements count to be 1, got=%d", len(program.Statements))
	}

	stat, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected statement type of ExpressionStatement, got=%T", program.Statements[0])
	}

	funct, ok := stat.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("expected statement type of CallExpression, got=%T", stat)
	}

	if !testLiteralExpression(t, funct.Function, "add") {
		return
	}

	if len(funct.Arguments) != 4 {
		t.Fatalf("expected arguments count to be 2, got=%d", len(funct.Arguments))
	}

	for i, test := range tests {
		if test != funct.Arguments[i].String() {
			t.Fatalf("expected argument to be %s, got %s", test, funct.Arguments[i].String())
		}
	}
}

func TestArrayLiteralParsing(t *testing.T) {
	input := `[10, 2*5, "yassinebenaid",func(){}];`

	tests := []string{
		"10",
		"(2 * 5)",
		"yassinebenaid",
		"func(){}",
	}

	l := lexer.New(input)
	par := New(l)
	program := par.ParseProgram()
	checkParserErrors(t, par)

	if len(program.Statements) != 1 {
		t.Fatalf("expected statements count to be 1, got=%d", len(program.Statements))
	}

	stat, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected statement type of ExpressionStatement, got=%T", program.Statements[0])
	}

	arr, ok := stat.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("expected statement type of ArrayLiteral, got=%T", stat)
	}

	if len(arr.Items) != 4 {
		t.Fatalf("expected array items count to be 4, got=%d", len(arr.Items))
	}

	for i, test := range tests {
		if test != arr.Items[i].String() {
			t.Fatalf("expected argument to be %s, got %s", test, arr.Items[i].String())
		}
	}
}

func TestArrayIndexParsing(t *testing.T) {
	tests := []struct {
		input string
		left  string
		index string
	}{
		{"[1,2,3][0]", "[1, 2, 3]", "0"},
		{"someFunc()[1]", "someFunc()", "1"},
		{"funcName[0]", "funcName", "0"},
		{"[1,2,3][3 - 2]", "[1, 2, 3]", "(3 - 2)"},
	}

	for _, test := range tests {
		l := lexer.New(test.input)
		par := New(l)
		program := par.ParseProgram()
		checkParserErrors(t, par)

		if len(program.Statements) != 1 {
			t.Fatalf("expected statements count to be 1, got=%d", len(program.Statements))
		}

		stat, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("expected statement type of ExpressionStatement, got=%T", program.Statements[0])
		}

		arr, ok := stat.Expression.(*ast.ArrayIndexExpression)
		if !ok {
			t.Fatalf("expected statement type of ArrayIndexExpression, got=%T", stat)
		}

		if test.left != arr.Left.String() {
			t.Fatalf("expected arr.Left be %s, got=%s", test.left, arr.Left.String())
		}

		if test.index != arr.Index.String() {
			t.Fatalf("expected arr.Left be %s, got=%s", test.index, arr.Index.String())
		}
	}
}

func TestHashLiteralParsing(t *testing.T) {
	input := `{name: "yassinebenaid",age: 10};`

	tests := []struct {
		expectedKey   string
		expectedValue any
	}{
		{"name", "yassinebenaid"},
		{"age", 10},
	}

	l := lexer.New(input)
	par := New(l)
	program := par.ParseProgram()
	checkParserErrors(t, par)

	if len(program.Statements) != 1 {
		t.Fatalf("expected statements count to be 1, got=%d", len(program.Statements))
	}

	stat, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected statement type of ExpressionStatement, got=%T", program.Statements[0])
	}

	hash, ok := stat.Expression.(*ast.HashLiteral)
	if !ok {
		t.Fatalf("expected statement type of HashLiteral, got=%T", stat)
	}

	if len(hash.Items) != 2 {
		t.Fatalf("expected hash items count to be 2, got=%d", len(hash.Items))
	}

	for _, tt := range tests {

		item, ok := hash.Items[tt.expectedKey]
		if !ok {
			t.Fatalf("failed to assert that the hash has item : %s", tt.expectedKey)
		}

		if item.String() != tt.expectedValue {
			t.Fatalf("failed to assert that the %s is to %s", tt.expectedValue, item)
		}
	}

}

func checkParserErrors(t *testing.T, p *Parser) {
	if len(p.errors) == 0 {
		return
	}

	t.Errorf("parser got %d errors\n", len(p.errors))

	for i, err := range p.errors {
		t.Errorf("#%d - parser error : %v\n", i, err)
	}

	t.FailNow()
}

func testIntegerLiteral(t *testing.T, exp ast.Expression, expected int64) bool {
	integ, ok := exp.(*ast.IntegerLiteral)

	if !ok {
		t.Errorf("expected expression of type IntegerLiteral, got %T", exp)
		return false
	}

	if integ.Value != expected {
		t.Errorf("failed to assert that the expected value of %d is equal to %d", expected, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", expected) {
		t.Errorf("failed to assert that the expected token literal of %s is equal to %s", fmt.Sprintf("%d", expected), integ.TokenLiteral())
		return false
	}

	return true
}

func testStringLiteral(t *testing.T, exp ast.Expression, expected string) bool {
	str, ok := exp.(*ast.StringLiteral)

	if !ok {
		t.Errorf("expected expression of type StringLiteral, got %T", exp)
		return false
	}

	if str.Value != expected {
		t.Errorf("failed to assert that the expected value of %s is equal to %s", expected, str.Value)
		return false
	}

	if str.TokenLiteral() != expected {
		t.Errorf("failed to assert that the expected token literal of %s is equal to %s", expected, str.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, expected bool) bool {
	boolExp, ok := exp.(*ast.BooleanLiteral)

	if !ok {
		t.Errorf("expected expression of type BooleanLiteral, got %T", exp)
		return false
	}

	if boolExp.Value != expected {
		t.Errorf("failed to assert that the expected value of %v is equal to %v", expected, boolExp.Value)
		return false
	}

	if boolExp.TokenLiteral() != fmt.Sprintf("%v", expected) {
		t.Errorf("failed to assert that the expected token literal of %v is equal to %v", expected, boolExp.TokenLiteral())
		return false
	}

	return true
}

func testIdentifierLiteral(t *testing.T, exp ast.Expression, expected string) bool {
	integ, ok := exp.(*ast.Identifier)

	if !ok {
		t.Errorf("expected expression of type Identifier, got %T", exp)
		return false
	}

	if integ.Value != expected {
		t.Errorf("failed to assert that the expected value of %s is equal to %s", expected, integ.Value)
		return false
	}

	if integ.TokenLiteral() != expected {
		t.Errorf("failed to assert that the expected token literal of %s is equal to %s", expected, integ.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected any) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifierLiteral(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}

	t.Errorf("type of %T is not handled", expected)
	return true
}

func testInfixExpression(t *testing.T, exp ast.Expression, left any, operator string, right any) bool {
	infix, ok := exp.(*ast.InfixExpression)

	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, infix.Left, left) {
		return false
	}

	if infix.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, infix.Operator)
		return false
	}

	if !testLiteralExpression(t, infix.Right, right) {
		return false
	}

	return true
}

func testPrefixExpression(t *testing.T, exp ast.Expression, operator string, right any) bool {
	infix, ok := exp.(*ast.PrefixExpression)

	if !ok {
		t.Errorf("exp is not ast.PrefixExpression. got=%T(%s)", exp, exp)
		return false
	}

	if infix.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, infix.Operator)
		return false
	}

	if !testLiteralExpression(t, infix.Right, right) {
		return false
	}

	return true
}
