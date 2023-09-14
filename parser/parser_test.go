package parser

import (
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

	ident, ok := stat.Expression.(*ast.Identifier)

	if ident.Value != "foobar" {
		t.Fatalf("expected identifir value to be [foobar] , got=%v", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Fatalf("expected identifir value to be [foobar] , got=%v", ident.TokenLiteral())
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

	integ, ok := stat.Expression.(*ast.IntegerLiteral)

	if !ok {
		t.Fatalf("statement type is incorrect, expected IntegerLiteral, got=%T", program.Statements[0])
	}

	if integ.Value != 5 {
		t.Fatalf("expected identifir value to be [foobar] , got=%v", integ.Value)
	}

	if integ.TokenLiteral() != "5" {
		t.Fatalf("expected identifir value to be [foobar] , got=%v", integ.TokenLiteral())
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
