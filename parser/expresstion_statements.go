package parser

import (
	"github.com/yassinebenaid/nishimia/ast"
	"github.com/yassinebenaid/nishimia/token"
)

func (p *Parser) parseExpressionStatement() ast.Statement {
	stat := &ast.ExpressionStatement{Token: p.currentToken}

	stat.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stat
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixPareseFns[p.currentToken.Type]

	if prefix == nil {
		p.noPrefixParseFnError(p.currentToken.Type)
		return nil
	}

	leftExp := prefix()

	return leftExp
}
