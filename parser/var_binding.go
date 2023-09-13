package parser

import (
	"github.com/yassinebenaid/nishimia/ast"
	"github.com/yassinebenaid/nishimia/token"
)

func (p *Parser) parseVarBindingStatement() ast.Statement {
	stat := &ast.VarStatement{Token: p.currentToken}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stat.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO: parse expression here
	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stat
}
