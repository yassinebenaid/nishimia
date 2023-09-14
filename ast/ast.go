package ast

import (
	"bytes"

	"github.com/yassinebenaid/nishimia/token"
)

type (
	Node interface {
		TokenLiteral() string
		String() string
	}

	Expression interface {
		Node
		expressionNode()
	}
	Statement interface {
		Node
		statementNode()
	}

	Program struct {
		Statements []Statement
	}
)

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}

func (p *Program) String() string {
	var s bytes.Buffer
	for _, stat := range p.Statements {
		s.WriteString(stat.String())
	}
	return s.String()
}

type VarStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (v *VarStatement) statementNode()       {}
func (v *VarStatement) TokenLiteral() string { return v.Token.Literal }

func (v *VarStatement) String() string {
	var s bytes.Buffer

	s.WriteString(v.TokenLiteral() + " ")
	s.WriteString(v.Name.String())
	s.WriteString(" = ")

	if v.Value != nil {
		s.WriteString(v.Value.String())
	}

	s.WriteString(";")

	return s.String()
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type ReturnStatement struct {
	Token  token.Token
	Return Expression
}

func (r *ReturnStatement) statementNode()       {}
func (r *ReturnStatement) TokenLiteral() string { return r.Token.Literal }

func (r *ReturnStatement) String() string {
	var s bytes.Buffer

	s.WriteString(r.TokenLiteral() + " ")

	if r.Return != nil {
		s.WriteString(r.Return.String())
	}

	s.WriteString(";")
	return s.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (e *ExpressionStatement) statementNode()       {}
func (e *ExpressionStatement) TokenLiteral() string { return e.Token.Literal }

func (e *ExpressionStatement) String() string {
	if e.Expression != nil {
		return e.Expression.String()
	}

	return ""
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntegerLiteral) expressionNode()      {}
func (i *IntegerLiteral) TokenLiteral() string { return i.Token.Literal }
func (i *IntegerLiteral) String() string       { return i.Token.Literal }
