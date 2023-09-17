package eval

import (
	"fmt"

	"github.com/yassinebenaid/nishimia/ast"
	"github.com/yassinebenaid/nishimia/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch v := node.(type) {
	case *ast.Program:
		return evalProgram(v.Statements)
	case *ast.ExpressionStatement:
		return Eval(v.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: v.Value}
	case *ast.BooleanLiteral:
		return nativeBooleanObject(v.Value)
	case *ast.PrefixExpression:
		val := Eval(v.Right)
		if isError(val) {
			return val
		}

		return evalPrefixExpression(v.Operator, val)
	case *ast.InfixExpression:
		l := Eval(v.Left)
		if isError(l) {
			return l
		}

		r := Eval(v.Right)
		if isError(r) {
			return r
		}

		return evalInfixExpression(v.Operator, l, r)
	case *ast.IfElseExpression:
		return evalConditionalExpression(v)
	case *ast.BlockStatement:
		return evalBlockStatements(v)
	case *ast.ReturnStatement:
		val := Eval(v.Return)
		if isError(val) {
			return val
		}

		return &object.ReturnValue{Value: val}
	}
	return NULL
}

func evalProgram(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatements(block *ast.BlockStatement) object.Object {
	var result object.Object

	for _, stmt := range block.Statements {
		result = Eval(stmt)

		if result != nil {
			if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	case "+":
		return evalPlusPrefixOperatorExpression(right)
	default:
		return newError(
			"unknown operator: %s%s",
			operator,
			right.Type(),
		)
	}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {

	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		return evalIntegerInfixExpression(operator, left, right)
	}

	if left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ {
		return evalBooleanInfixExpression(operator, left, right)
	}

	return newError(
		"invalid operation: %s %s %s (mismatched types %s and %s)",
		left.Inspect(),
		operator,
		right.Inspect(),
		left.Type(),
		right.Type(),
	)
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return newError(
			"invalid operation: !%s (operator \"!\" not defined on %s)",
			right.Inspect(),
			right.Type(),
		)
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {

	if right.Type() != object.INTEGER_OBJ {
		return newError(
			"invalid operation: -%s (operator \"-\" not defined on %s)",
			right.Inspect(),
			right.Type(),
		)
	}

	value := right.(*object.Integer).Value

	return &object.Integer{Value: value * -1}
}

func evalPlusPrefixOperatorExpression(right object.Object) object.Object {

	if right.Type() != object.INTEGER_OBJ {
		return newError(
			"invalid operation: +%s (operator \"+\" not defined on %s)",
			right.Inspect(),
			right.Type(),
		)
	}

	value := right.(*object.Integer).Value

	return &object.Integer{Value: value}
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftValue + rightValue}
	case "-":
		return &object.Integer{Value: leftValue - rightValue}
	case "*":
		return &object.Integer{Value: leftValue * rightValue}
	case "/":
		return &object.Integer{Value: leftValue / rightValue}
	case "<":
		return nativeBooleanObject(leftValue < rightValue)
	case ">":
		return nativeBooleanObject(leftValue > rightValue)
	case ">=":
		return nativeBooleanObject(leftValue >= rightValue)
	case "<=":
		return nativeBooleanObject(leftValue <= rightValue)
	case "!=":
		return nativeBooleanObject(leftValue != rightValue)
	case "==":
		return nativeBooleanObject(leftValue == rightValue)
	default:
		return newError(
			"invalid operation: %d %s %d",
			leftValue,
			operator,
			rightValue,
		)
	}
}

func evalBooleanInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftValue := left.(*object.Boolean).Value
	rightValue := right.(*object.Boolean).Value

	switch operator {
	case "&&":
		return nativeBooleanObject(leftValue && rightValue)
	case "||":
		return nativeBooleanObject(leftValue || rightValue)
	case "==":
		return nativeBooleanObject(leftValue == rightValue)
	case "!=":
		return nativeBooleanObject(leftValue != rightValue)
	default:
		return newError(
			"invalid operation: %t %s %t",
			leftValue,
			operator,
			rightValue,
		)
	}
}

func evalConditionalExpression(node *ast.IfElseExpression) object.Object {
	cond := Eval(node.Condition)
	if isError(cond) {
		return cond
	}

	if cond.Type() != object.BOOLEAN_OBJ {
		return newError("non-boolean value in if-statement , ( got=%s, want=BOOLEAN )", cond.Type())
	}

	if cond.Inspect() == "true" {
		return Eval(node.Consequence)
	}

	if node.Alternative != nil {
		return Eval(node.Alternative)
	}

	return NULL
}

func nativeBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}

	return FALSE
}

func newError(msg string, args ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(msg, args...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}
