package parser

import (
	"bytes"
	// "fmt"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type ExpressionStatement struct {
	Token      Token
	Expression Expression
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

func (es *ExpressionStatement) statementNode() {}

func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.TokenValue
}

type Identifier struct {
	Token Token
	Value string
}

func (i *Identifier) expressionNode() {}

func (i *Identifier) String() string {
	return i.Value
}

func (i *Identifier) TokenLiteral() string {
	return i.Token.TokenValue
}

type IntegerLiteral struct {
	Token Token
	Value int
}

func (il *IntegerLiteral) expressionNode() {}

func (il *IntegerLiteral) String() string {
	return il.Token.TokenValue
}

func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.TokenValue
}

type FloatLiteral struct {
	Token Token
	Value float64
}

func (il *FloatLiteral) expressionNode() {}

func (il *FloatLiteral) String() string {
	return il.Token.TokenValue
}

func (il *FloatLiteral) TokenLiteral() string {
	return il.Token.TokenValue
}

type StringLiteral struct {
	Token Token
	Value string
}

func (sl *StringLiteral) expressionNode() {}

func (sl *StringLiteral) String() string {
	return sl.Token.TokenValue
}

func (sl *StringLiteral) TokenLiteral() string {
	return sl.Token.TokenValue
}

type Boolean struct {
	Token Token
	Value bool
}

func (il *Boolean) expressionNode() {}

func (il *Boolean) String() string {
	return il.Token.TokenValue
}

func (il *Boolean) TokenLiteral() string {
	return il.Token.TokenValue
}

type PrefixExpression struct {
	Token    Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}

func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.TokenValue
}

type InfixExpression struct {
	Token    Token
	Right    Expression
	Operator string
	Left     Expression
}

func (ie *InfixExpression) expressionNode() {}

func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.TokenValue
}

type IfExpression struct {
	Token       Token
	Condition   Expression
	Consequence Expression
	Alternative Expression
}

func (ie *IfExpression) expressionNode() {}

func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("IF (")
	out.WriteString(ie.Condition.String())
	out.WriteString(",")
	out.WriteString(ie.Consequence.String())
	if ie.Alternative != nil {
		out.WriteString(",")
		out.WriteString(ie.Alternative.String())
	}
	out.WriteString(")")

	return out.String()
}

func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.TokenValue
}

// CallExpression represents a function call like MIN(a, b) or MAX(x, y, z)
type CallExpression struct {
	Token     Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}

func (ce *CallExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	for i, arg := range ce.Arguments {
		if i > 0 {
			out.WriteString(", ")
		}
		out.WriteString(arg.String())
	}
	out.WriteString(")")
	return out.String()
}

func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.TokenValue
}
