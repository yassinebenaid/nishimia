package parser

import (
	"fmt"
	"strconv"

	"github.com/yassinebenaid/nishimia/ast"
)

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
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
