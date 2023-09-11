package lexer

import (
	"github.com/yassinebenaid/nishimia/token"
)

type Lexer struct {
	input        string // the source code
	position     int    // the current position , points to the index of ch
	readPosition int    // the current read position, the next character after
	ch           byte   // the haracter under examination
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhiteSpace()

	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, '=')
	case '+':
		tok = newToken(token.PLUS, '+')
	case '-':
		tok = newToken(token.MINUS, '-')
	case '*':
		tok = newToken(token.MULTP, '*')
	case '(':
		tok = newToken(token.LPARENT, '(')
	case ')':
		tok = newToken(token.RPARENT, ')')
	case '{':
		tok = newToken(token.LBRACE, '{')
	case '}':
		tok = newToken(token.RBRACE, '}')
	case ',':
		tok = newToken(token.COMMA, ',')
	case ';':
		tok = newToken(token.SEMICOLON, ';')
	case '/':
		tok = newToken(token.SLASH, '/')
	case '!':
		tok = newToken(token.BANG, '!')
	case '<':
		tok = newToken(token.LT, '<')
	case '>':
		tok = newToken(token.GT, '>')
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLIGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	l.position = l.readPosition
	l.readPosition++
}

func (l *Lexer) readIdentifier() string {
	position := l.position

	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position

	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[position:l.position]
}

func (l *Lexer) skipWhiteSpace() {
	for isWhiteSpace(l.ch) {
		l.readChar()
	}
}

func newToken(t token.TokenType, v byte) token.Token {
	return token.Token{Type: t, Literal: string(v)}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(v byte) bool {
	return '0' <= v && v <= '9'
}

func isWhiteSpace(v byte) bool {
	return v == ' ' || v == '\t' || v == '\n' || v == '\r'
}
