package parser

import (
	"github.com/yassinebenaid/nishimia/ast"
	"github.com/yassinebenaid/nishimia/lexer"
	"github.com/yassinebenaid/nishimia/token"
)

type Parser struct {
	lex *lexer.Lexer

	currentToken token.Token
	peekToken    token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{lex: l}

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lex.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	var stat ast.Statement
	var prog = &ast.Program{}

	for ; !p.currentTokenIs(token.EOF); p.nextToken() {
		if stat = p.parseStatement(); stat != nil {
			prog.Statements = append(prog.Statements, stat)
		}
	}

	return prog
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.VAR:
		return p.parseVarBindingStatement()
	default:
		return nil
	}
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}

	return false
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return t == p.peekToken.Type
}

func (p *Parser) currentTokenIs(t token.TokenType) bool {
	return t == p.currentToken.Type
}
