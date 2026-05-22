package parser

import (
	// "fmt"
	"strings"
	"text/scanner"
)

type Lexer struct {
	input                string
	position             int
	readPosition         int
	ch                   string
	scannerTokenType     rune
	s                    scanner.Scanner
	peekScannerTokenType rune
}

func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	var s scanner.Scanner
	s.Init(strings.NewReader(l.input))
	l.s = s
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	tok := l.s.Scan()
	l.ch = l.s.TokenText()
	l.scannerTokenType = tok //scanner.TokenString(tok)
	l.position = l.s.Position.Offset
	l.readPosition = l.position + 1
	l.peekScannerTokenType = l.s.Peek()
}

func newToken(tokenType TokenType, ch string) Token {
	if tokenType == STRING {
		ch = strings.Trim(ch, `"`)
	}
	return Token{TokenType: tokenType, TokenValue: ch}
}

func (l *Lexer) NextToken() Token {
	var tok Token
	switch l.scannerTokenType {
	case '(':
		tok = newToken(LPAREN, l.ch)
	case ')':
		tok = newToken(RPAREN, l.ch)
	case ',':
		tok = newToken(COMMA, l.ch)
	case scanner.Ident:
		switch l.ch {
		case "IF":
			tok = newToken(IF, l.ch)
		case "TRUE":
			tok = newToken(BOOL, l.ch)
		case "FALSE":
			tok = newToken(BOOL, l.ch)
		default:
			tok = newToken(IDENT, l.ch)
		}
	case '+':
		tok = newToken(PLUS, l.ch)
	case '-':
		tok = newToken(MINUS, l.ch)
	case '*':
		tok = newToken(ASTERISK, l.ch)
	case '/':
		tok = newToken(SLASH, l.ch)
	case scanner.EOF:
		tok = newToken(EOF, l.ch)
	case scanner.Int:
		tok = newToken(INT, l.ch)
	case scanner.Float:
		tok = newToken(FLOAT, l.ch)
	case scanner.String:
		tok = newToken(STRING, l.ch)
	case '<':
		switch l.peekScannerTokenType {
		case '>':
			tok = newToken(NOT_EQ, "<>")
			l.readChar()
		default:
			tok = newToken(LT, l.ch)
		}
	case '>':
		tok = newToken(GT, l.ch)
	case '=':
		tok = newToken(EQ, l.ch)
	default:
		tok = newToken(ILLEGAL, l.ch)
	}
	l.readChar()
	return tok
}
