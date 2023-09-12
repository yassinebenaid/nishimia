package parser

import (
	"testing"

	"github.com/yassinebenaid/nishimia/ast"
	"github.com/yassinebenaid/nishimia/lexer"
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
