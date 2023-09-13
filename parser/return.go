package parser

import (
	"github.com/yassinebenaid/nishimia/ast"
	"github.com/yassinebenaid/nishimia/token"
)

func (p *Parser) parseReturnStatement() ast.Statement {
	stat := &ast.ReturnStatement{Token: p.currentToken}

	p.nextToken()

	// TODO : parse experession here .
	for !p.currentTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stat
}
