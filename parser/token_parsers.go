package parser

import (
	"fmt"
	"strconv"

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

	p.nextToken()

	stat.Value = p.parseExpression(LOWEST)

	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}

	return stat
}

func (p *Parser) parseReturnStatement() ast.Statement {
	stat := &ast.ReturnStatement{Token: p.currentToken}

	p.nextToken()

	stat.Return = p.parseExpression(LOWEST)

	for !p.expectPeek(token.SEMICOLON) {
		return nil
	}

	return stat
}

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

func (p *Parser) parseString() ast.Expression {
	return &ast.StringLiteral{Token: p.currentToken, Value: p.currentToken.Literal}
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

func (p *Parser) parseIfElseExpression() ast.Expression {
	exp := &ast.IfElseExpression{
		Token: p.currentToken,
	}

	p.nextToken()

	exp.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	exp.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		exp.Alternative = p.parseBlockStatement()
	}

	return exp
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.currentToken}
	block.Statements = make([]ast.Statement, 0)

	p.nextToken()

	var stat ast.Statement

	for !p.currentTokenIs(token.RBRACE) && !p.currentTokenIs(token.EOF) {
		if stat = p.parseStatement(); stat != nil {
			block.Statements = append(block.Statements, stat)
		}

		p.nextToken()
	}

	return block
}

func (p *Parser) parseFunctionExpression() ast.Expression {
	exp := &ast.FunctionLiteral{Token: p.currentToken}

	if !p.expectPeek(token.LPARENT) {
		return nil
	}

	exp.Params = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	exp.Body = p.parseBlockStatement()

	return exp
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	var args []*ast.Identifier

	if p.peekTokenIs(token.RPARENT) {
		p.nextToken()
		return args
	}

	p.nextToken()

	args = append(args, &ast.Identifier{
		Token: p.currentToken,
		Value: p.currentToken.Literal,
	})

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, &ast.Identifier{
			Token: p.currentToken,
			Value: p.currentToken.Literal,
		})
	}

	if !p.expectPeek(token.RPARENT) {
		return nil
	}

	return args
}

func (p *Parser) parseFunctionCallExpression(function ast.Expression) ast.Expression {
	fnExp := &ast.CallExpression{Token: p.currentToken, Function: function}
	fnExp.Arguments = p.parseCallArguments()

	return fnExp
}

func (p *Parser) parseCallArguments() []ast.Expression {
	var args []ast.Expression

	if p.peekTokenIs(token.RPARENT) {
		p.nextToken()
		return args
	}

	p.nextToken()

	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPARENT) {
		return nil
	}

	return args
}

func (p *Parser) parseArrayExpression() ast.Expression {
	var items []ast.Expression

	if p.peekTokenIs(token.RBRACKET) {
		p.nextToken()
		return &ast.ArrayLiteral{Items: items}
	}

	p.nextToken()

	items = append(items, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		items = append(items, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return &ast.ArrayLiteral{Items: items}
}

func (p *Parser) parseHashExpression() ast.Expression {
	items := map[ast.Expression]ast.Expression{}

	if p.peekTokenIs(token.RBRACE) {
		p.nextToken()
		return &ast.HashLiteral{Items: items}
	}

	p.nextToken()

	key := p.parseExpression(LOWEST)

	if !p.expectPeek(token.COLON) {
		return nil
	}

	p.nextToken()

	value := p.parseExpression(LOWEST)

	items[key] = value

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()

		key := p.parseExpression(LOWEST)

		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.nextToken()
		value := p.parseExpression(LOWEST)
		items[key] = value
	}

	if !p.expectPeek(token.RBRACE) {
		return nil
	}

	return &ast.HashLiteral{Items: items}
}

func (p *Parser) parseArrayIndexExpression(left ast.Expression) ast.Expression {
	var exp ast.IndexExpression

	exp.Left = left
	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return &exp
}
