package parser

import (
	"fmt"

	"github.com/yassinebenaid/nishimia/ast"
	"github.com/yassinebenaid/nishimia/lexer"
	"github.com/yassinebenaid/nishimia/token"
)

type Parser struct {
	lex *lexer.Lexer // the lexer to gain tokens

	currentToken token.Token // refers to the current token under examination
	peekToken    token.Token // refers to the next token after currentToken

	errors []string // holds all parsing errors
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
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return nil
	}
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}

	p.peekError(t)

	return false
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return t == p.peekToken.Type
}

func (p *Parser) currentTokenIs(t token.TokenType) bool {
	return t == p.currentToken.Type
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.TokenType) {
	p.errors = append(
		p.errors,
		fmt.Sprintf(`unexpected token  "%s" , expected "%s"`, p.peekToken.Literal, t),
	)
}
