package ast

import "github.com/yassinebenaid/nishimia/token"

type VarStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (v *VarStatement) statementNode()       {}
func (v *VarStatement) TokenLiteral() string { return v.Token.Literal }

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
