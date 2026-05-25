package parser

import (
	"fmt"
	"strconv"
)

// import (
// )

const (
	_ int = iota
	LOWEST
	LESSGREATER
	EQUAL
	SUM
	PRODUCT
	MODULO_PREC
	POWER_PREC
	PREFIX
	// PAREN
)

var precedences = map[TokenType]int{
	EQ:       EQUAL,
	NOT_EQ:   EQUAL,
	LT:       LESSGREATER,
	GT:       LESSGREATER,
	PLUS:     SUM,
	MINUS:    SUM,
	ASTERISK: PRODUCT,
	SLASH:    PRODUCT,
	MODULO:   MODULO_PREC,
	POWER:    POWER_PREC,
	// LPAREN:PAREN,
	// RPAREN:PAREN,
}

type Parser struct {
	l      *Lexer
	errors []string

	currentToken Token
	peekToken    Token

	prefixParseFns map[TokenType]prefixParseFn
	infixParseFns  map[TokenType]infixParseFn
}

func (p *Parser) registerPrefixFunction(tokenType TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfixFunction(tokenType TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func NewParser(l *Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	p.prefixParseFns = make(map[TokenType]prefixParseFn)
	p.infixParseFns = make(map[TokenType]infixParseFn)

	p.registerPrefixFunction(IDENT, p.parseIdentifier)
	p.registerPrefixFunction(INT, p.parseIntegerLiteral)
	p.registerPrefixFunction(MINUS, p.parsePrefixExpression)
	p.registerPrefixFunction(IF, p.parseIfExpression)
	p.registerPrefixFunction(BOOL, p.parseBoolean)
	p.registerPrefixFunction(STRING, p.parseStringLiteral)
	p.registerPrefixFunction(FLOAT, p.parseFloatLiteral)
	p.registerPrefixFunction(MIN, p.parseFunctionCall)
	p.registerPrefixFunction(MAX, p.parseFunctionCall)
	p.registerPrefixFunction(ROUND, p.parseFunctionCall)

	p.registerInfixFunction(EQ, p.parseInfixExpression)
	p.registerInfixFunction(NOT_EQ, p.parseInfixExpression)
	p.registerInfixFunction(LT, p.parseInfixExpression)
	p.registerInfixFunction(GT, p.parseInfixExpression)
	p.registerInfixFunction(PLUS, p.parseInfixExpression)
	p.registerInfixFunction(MINUS, p.parseInfixExpression)
	p.registerInfixFunction(ASTERISK, p.parseInfixExpression)
	p.registerInfixFunction(SLASH, p.parseInfixExpression)
	p.registerInfixFunction(MODULO, p.parseInfixExpression)
	p.registerInfixFunction(POWER, p.parseInfixExpression)

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) parseIdentifier() Expression {
	return &Identifier{Token: p.currentToken, Value: p.currentToken.TokenValue}
}

func (p *Parser) parseIntegerLiteral() Expression {
	lit := &IntegerLiteral{Token: p.currentToken}
	value, err := strconv.Atoi(p.currentToken.TokenValue)
	if err != nil {
		msg := fmt.Sprintf("Cannot parse %q as Integer", p.currentToken.TokenValue)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) parseFloatLiteral() Expression {
	lit := &FloatLiteral{Token: p.currentToken}
	value, err := strconv.ParseFloat(p.currentToken.TokenValue, 64)
	if err != nil {
		msg := fmt.Sprintf("Cannot parse %q as Float", p.currentToken.TokenValue)
		p.errors = append(p.errors, msg)
		return nil
	}
	lit.Value = value
	return lit
}

func (p *Parser) parseStringLiteral() Expression {
	return &StringLiteral{Token: p.currentToken, Value: p.currentToken.TokenValue}
}

func (p *Parser) parseBoolean() Expression {
	lit := &Boolean{Token: p.currentToken}
	switch p.currentToken.TokenValue {
	case "TRUE":
		lit.Value = true
	case "FALSE":
		lit.Value = false
	default:
		msg := fmt.Sprintf("Cannot parse %q as Boolean", p.currentToken.TokenValue)
		p.errors = append(p.errors, msg)
		return nil
	}
	return lit
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *Program {
	program := &Program{}
	program.Statements = []Statement{}

	for p.currentToken.TokenType != EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) currentTokenIs(t TokenType) bool {
	return p.currentToken.TokenType == t
}

func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.TokenType == t
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t TokenType) {
	msg := fmt.Sprintf("Expected next token to be %s, got %s instead", t, p.peekToken.TokenType)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseStatement() Statement {
	switch p.currentToken.TokenType {
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseExpressionStatement() *ExpressionStatement {
	stmt := &ExpressionStatement{Token: p.currentToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(EOF) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpression(precedence int) Expression {
	prefix := p.prefixParseFns[p.currentToken.TokenType]
	if prefix == nil {
		p.noPrefixFnParseError(p.currentToken.TokenType)
		return nil
	}

	leftExp := prefix()
	for !p.peekTokenIs(EOF) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.TokenType]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) noPrefixFnParseError(t TokenType) {
	msg := fmt.Sprintf("No prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) parsePrefixExpression() Expression {
	expression := &PrefixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.TokenValue,
	}
	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)
	return expression
}

func (p *Parser) parseIfExpression() Expression {
	expression := &IfExpression{Token: p.currentToken}

	if !p.expectPeek(LPAREN) {
		return nil
	}

	p.nextToken()

	expression.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(COMMA) {
		return nil
	}

	p.nextToken()

	expression.Consequence = p.parseExpression(LOWEST)

	if !p.peekTokenIs(COMMA) {
		if !p.expectPeek(RPAREN) {
			return nil
		}
		p.nextToken()
		return expression
	}

	p.nextToken()
	p.nextToken()

	expression.Alternative = p.parseExpression(LOWEST)

	p.nextToken()

	return expression

}

// parseFunctionCall parses function calls like MIN(a, b), MAX(x, y, z), ROUND(n, decimals)
func (p *Parser) parseFunctionCall() Expression {
	call := &CallExpression{
		Token:    p.currentToken,
		Function: &Identifier{Token: p.currentToken, Value: p.currentToken.TokenValue},
	}

	if !p.expectPeek(LPAREN) {
		return nil
	}

	call.Arguments = p.parseCallArguments()

	return call
}

// parseCallArguments parses comma-separated function arguments
func (p *Parser) parseCallArguments() []Expression {
	args := []Expression{}

	if p.peekTokenIs(RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(RPAREN) {
		return nil
	}

	return args
}

func (p *Parser) parseInfixExpression(left Expression) Expression {
	expression := &InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.TokenValue,
		Left:     left,
	}
	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}

func (p *Parser) parseGroupedExpression() Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(RPAREN) {
		return nil
	}

	return exp

}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.TokenType]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.currentToken.TokenType]; ok {
		return p
	}
	return LOWEST
}

type (
	prefixParseFn func() Expression
	infixParseFn  func(Expression) Expression
)
