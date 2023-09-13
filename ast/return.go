package ast

import "github.com/yassinebenaid/nishimia/token"

type ReturnStatement struct {
	Token  token.Token
	Return Expression
}

func (r *ReturnStatement) statementNode()       {}
func (r *ReturnStatement) TokenLiteral() string { return r.Token.Literal }
