package parser

import (
// "fmt"
)

type TokenType string

type Token struct {
	TokenValue string
	TokenType  TokenType
}

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"
	// Identifiers + literals
	IDENT TokenType = "Ident" // add, foobar, x, y, ...

	INT    TokenType = "Int" // 1343456
	FLOAT  TokenType = "Float"
	STRING TokenType = "String"
	BOOL   TokenType = "Boolean"
	// Operators
	PLUS     TokenType = "+"
	MINUS    TokenType = "-"
	ASTERISK TokenType = "*"
	SLASH    TokenType = "/"
	MODULO   TokenType = "%"
	POWER    TokenType = "^"

	IF TokenType = "IF"

	EQ     TokenType = "="
	NOT_EQ TokenType = "<>"

	// Built-in Functions
	MIN   TokenType = "MIN"
	MAX   TokenType = "MAX"
	ROUND TokenType = "ROUND"

	// Delimiters
	COMMA  TokenType = ","
	LPAREN TokenType = "("
	RPAREN TokenType = ")"

	LT TokenType = "<"
	GT TokenType = ">"
)

// ParanthesisTokens = []string{"(", ")"}
// 	OperatorTokens = []string{"+", "/", "-", "*"}
// 	AttributesAndConstantTokens = []string{"Ident", "Int", "Float", "String"}

// func tokenizer(formula string) []Token {
// 	var s scanner.Scanner
// 	// s.Mode = scanner.GoTokens
// 	s.Init(strings.NewReader(formula))
// 	// s.Next()
// 	tokens := []Token{}
// 	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
// 		value := Token{
// 			TokenValue: s.TokenText(),
// 			TokenType: scanner.TokenString(tok),
// 		}
// 		tokens = append(tokens, value)
// 	}
// 	return tokens
// }

// func main(){
// 	var tokens = tokenizer(`map_123 + ("xyz" * ahgd) + 123.376`)
// 	tokenValidator(tokens)
// }

// func tokenValidator(tokens []Token){
// 	AcceptedTokens := constants.AcceptedTokens
// 	errorMessages := []string{}
// 	for _, token := range tokens {
// 		if !slices.Contains(AcceptedTokens, token.TokenType){
// 			errorMessages = append(errorMessages, token.TokenValue + ` is not allowed in formula`)
// 		}
// 	}
// }

// func formulaSyntaxValidator(tokens []Token){
// 	//1. No Consecutive Operators

// 	//2. No Consecutive Attributes or Constants

// 	//3. Balanced Paranthesis
// }
