package ast

import (
	"bytes"
	"strings"

	"github.com/yassinebenaid/nishimia/token"
)

type (
	Node interface {
		TokenLiteral() string
		String() string
	}

	Expression interface {
		Node
		expressionNode()
	}
	Statement interface {
		Node
		statementNode()
	}
)

// This is the base node for every program,
// this turns out that every program in nishimia is just a sequence of statements
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}

func (p *Program) String() string {
	var s bytes.Buffer
	for _, stat := range p.Statements {
		s.WriteString(stat.String())
	}
	return s.String()
}

// This node represents a variable binding statement, typically
// any line that looks like :
//
//	var name = value
type VarStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (v *VarStatement) statementNode()       {}
func (v *VarStatement) TokenLiteral() string { return v.Token.Literal }

func (v *VarStatement) String() string {
	var s bytes.Buffer

	s.WriteString(v.TokenLiteral() + " ")
	s.WriteString(v.Name.String())
	s.WriteString(" = ")

	if v.Value != nil {
		s.WriteString(v.Value.String())
	}

	s.WriteString(";")

	return s.String()
}

// This node represents any identifier in the language, this includes variables, functions, constants etc.
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// This node represents the return statement, typically any statement that looks like :
//
//	return expression
type ReturnStatement struct {
	Token  token.Token
	Return Expression
}

func (r *ReturnStatement) statementNode()       {}
func (r *ReturnStatement) TokenLiteral() string { return r.Token.Literal }

func (r *ReturnStatement) String() string {
	var s bytes.Buffer

	s.WriteString(r.TokenLiteral() + " ")

	if r.Return != nil {
		s.WriteString(r.Return.String())
	}

	s.WriteString(";")
	return s.String()
}

// This node represents the expressions in nishimia,
//
// expressions in nishimia are any statement the produces a value , here is some examples
//
//	1 + 2
//
//	variableName + 4
//
//	functionName(arg1,arg2)
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (e *ExpressionStatement) statementNode()       {}
func (e *ExpressionStatement) TokenLiteral() string { return e.Token.Literal }

func (e *ExpressionStatement) String() string {
	if e.Expression != nil {
		return e.Expression.String()
	}

	return ""
}

// This node represents all the integer literals in the language ,
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntegerLiteral) expressionNode()      {}
func (i *IntegerLiteral) TokenLiteral() string { return i.Token.Literal }
func (i *IntegerLiteral) String() string       { return i.Token.Literal }

// This node represents all the prefix expressions like :
//
//	-1 // prefix is -
//
//	!true // prefix is !
//
//	!variableName // prefix is !
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (p *PrefixExpression) expressionNode()      {}
func (p *PrefixExpression) TokenLiteral() string { return p.Token.Literal }
func (p *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(p.Operator)
	out.WriteString(p.Right.String())
	out.WriteString(")")

	return out.String()
}

// This node represents all the infix expressions like :
//
//	1 + 1 // operator is + (infix)
type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (p *InfixExpression) expressionNode()      {}
func (p *InfixExpression) TokenLiteral() string { return p.Token.Literal }
func (p *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(p.Left.String())
	out.WriteString(" ")
	out.WriteString(p.Operator)
	out.WriteString(" ")
	out.WriteString(p.Right.String())
	out.WriteString(")")

	return out.String()
}

// This node represents the boolean literals in the language :
//
//	true | false
type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (i *BooleanLiteral) expressionNode()      {}
func (i *BooleanLiteral) TokenLiteral() string { return i.Token.Literal }
func (i *BooleanLiteral) String() string       { return i.Token.Literal }

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

// This node represents the if-else expression
type IfElseExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (i *IfElseExpression) expressionNode()      {}
func (i *IfElseExpression) TokenLiteral() string { return i.Token.Literal }
func (i *IfElseExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(i.Condition.String())
	out.WriteString(" ")
	out.WriteString(i.Consequence.String())

	if i.Alternative != nil {
		out.WriteString(" else ")
		out.WriteString(i.Alternative.String())
	}

	return out.String()
}

// This node represents the blocks of statements, typically any block between braces like
// inside if-else statements, or function definitions etc.
type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (b *BlockStatement) expressionNode()      {}
func (b *BlockStatement) TokenLiteral() string { return b.Token.Literal }
func (b *BlockStatement) String() string {
	var out bytes.Buffer

	out.WriteString("{")

	for _, stat := range b.Statements {
		out.WriteString(stat.String())
	}

	out.WriteString("}")

	return out.String()
}

// This node represents the function literal,
//
//	fn(){}
type FunctionLiteral struct {
	Token  token.Token
	Params []*Identifier
	Body   *BlockStatement
}

func (fn *FunctionLiteral) expressionNode()      {}
func (fn *FunctionLiteral) TokenLiteral() string { return fn.Token.Literal }
func (fn *FunctionLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("func(")

	var params []string
	for _, param := range fn.Params {
		params = append(params, param.String())
	}

	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(fn.Body.String())

	return out.String()
}

// This node represents the function calls ,.
//
//	functionName() // using identifier
//	fn(){}() //using function literal
type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (c *CallExpression) expressionNode()      {}
func (c *CallExpression) TokenLiteral() string { return c.Token.Literal }
func (c *CallExpression) String() string {
	var out bytes.Buffer

	var args []string
	for _, arg := range c.Arguments {
		args = append(args, arg.String())
	}

	out.WriteString(c.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

// This node represents the array ,.
type ArrayLiteral struct {
	Token token.Token
	Items []Expression
}

func (a *ArrayLiteral) expressionNode()      {}
func (a *ArrayLiteral) TokenLiteral() string { return a.Token.Literal }
func (a *ArrayLiteral) String() string {
	var out bytes.Buffer

	var items []string
	for _, item := range a.Items {
		items = append(items, item.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(items, ", "))
	out.WriteString("]")

	return out.String()
}

// This node represents the array index like array[index] ,.
type ArrayIndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (a *ArrayIndexExpression) expressionNode()      {}
func (a *ArrayIndexExpression) TokenLiteral() string { return a.Token.Literal }
func (a *ArrayIndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString(a.Left.String())
	out.WriteString("[")
	out.WriteString(a.Index.String())
	out.WriteString("]")

	return out.String()
}
