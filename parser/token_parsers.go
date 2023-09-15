package parser

import (
	"fmt"
	"strconv"

	"github.com/yassinebenaid/nishimia/ast"
	"github.com/yassinebenaid/nishimia/token"
)

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parseBoolean() ast.Expression {
	if p.currentTokenIs(token.TRUE) {
		return &ast.BooleanLiteral{Token: p.currentToken, Value: true}
	}

	return &ast.BooleanLiteral{Token: p.currentToken, Value: false}
}

func (p *Parser) parseInteger() ast.Expression {
	exp := &ast.IntegerLiteral{Token: p.currentToken}

	i, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)

	if err != nil {
		p.errors = append(p.errors, fmt.Sprintf("couldn't parse %v as integer, %v", p.currentToken.Literal, err))
		return nil
	}

	exp.Value = i

	return exp
}

func (p *Parser) parsePrefixExpressions() ast.Expression {
	exp := &ast.PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
	}

	p.nextToken()

	exp.Right = p.parseExpression(PREFIX)

	return exp
}

func (p *Parser) parseInfixExpressions(left ast.Expression) ast.Expression {
	exp := &ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
		Left:     left,
	}

	precedence := p.currentPrecedence()
	p.nextToken()

	exp.Right = p.parseExpression(precedence)

	return exp
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPARENT) {
		return nil
	}

	return exp
}
