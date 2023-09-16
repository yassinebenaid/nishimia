package parser

import (
	"fmt"
	"testing"

	"github.com/yassinebenaid/nishimia/ast"
	"github.com/yassinebenaid/nishimia/lexer"
	"github.com/yassinebenaid/nishimia/token"
)

func TestVarStatement(t *testing.T) {
	input := `
	var year = 2023;
	var month = 9;
	var day = 12;
	`
	lex := lexer.New(input)
	par := New(lex)

	program := par.ParseProgram()
	checkParserErrors(t, par)

	if program == nil {
		t.Fatal("ParseProgram() returned nil !")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("failed to parse program, expected statements count of 3 , got : %d", len(program.Statements))
	}

	cases := []struct {
		expectedIdentifier string
	}{
		{"year"},
		{"month"},
		{"day"},
	}

	for i, cas := range cases {
		stat := program.Statements[i]

		if !testVarStatement(t, stat, cas.expectedIdentifier) {
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
	input := `
	return 5;
	return 5 == 7;
	return add(4,5);
	`

	lex := lexer.New(input)
	par := New(lex)

	program := par.ParseProgram()
	checkParserErrors(t, par)

	if len(program.Statements) != 3 {
		t.Fatalf("expected 3 statemnts, got %d", len(program.Statements))
	}

	for _, stat := range program.Statements {
		ret, ok := stat.(*ast.ReturnStatement)

		if !ok {
			t.Errorf("the statement type is not ast.ReturnStatement, got %T", stat)
			continue
		}

		if ret.Token.Literal != "return" {
			t.Errorf("statement.Token.Literal is not 'return', got %s", ret.Token.Literal)
		}

		if ret.TokenLiteral() != "return" {
			t.Errorf("statement.Token.Literal is not 'return', got %s", ret.TokenLiteral())
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

		{"(a + b) * 2", "((a + b) * 2)"},
		{"(a + b) / 2", "((a + b) / 2)"},
		{"!(false == true)", "(!(false == true))"},
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
			x + y;
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

	// ret, ok := funct.Block.Statements[0].(*ast.ReturnStatement)
	// if !ok {
	// 	t.Fatalf("expected first statement type of ReturnStatement, got=%T", funct.Block.Statements[0])
	// }

	// if !testInfixExpression(t, ret.Return, "x", "+", "y") {
	// 	return
	// }

	ret, ok := funct.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected first statement type of ExpressionStatement, got=%T", funct.Body.Statements[0])
	}

	if !testInfixExpression(t, ret.Expression, "x", "+", "y") {
		return
	}
}

func TestFunctionLiteralWithOneParamParsing(t *testing.T) {
	input := `
		func(x){
			x * 2;
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

	// ret, ok := funct.Block.Statements[0].(*ast.ReturnStatement)
	// if !ok {
	// 	t.Fatalf("expected first statement type of ReturnStatement, got=%T", funct.Block.Statements[0])
	// }

	// if !testInfixExpression(t, ret.Return, "x", "+", "y") {
	// 	return
	// }

	ret, ok := funct.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected first statement type of ExpressionStatement, got=%T", funct.Body.Statements[0])
	}

	if !testInfixExpression(t, ret.Expression, "x", "*", 2) {
		return
	}
}

func TestFunctionLiteralWithManyParamsParsing(t *testing.T) {
	input := `
		func(x,y){
			x + y;
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

	// ret, ok := funct.Block.Statements[0].(*ast.ReturnStatement)
	// if !ok {
	// 	t.Fatalf("expected first statement type of ReturnStatement, got=%T", funct.Block.Statements[0])
	// }

	// if !testInfixExpression(t, ret.Return, "x", "+", "y") {
	// 	return
	// }

	ret, ok := funct.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("expected first statement type of ExpressionStatement, got=%T", funct.Body.Statements[0])
	}

	if !testInfixExpression(t, ret.Expression, "x", "+", "y") {
		return
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
