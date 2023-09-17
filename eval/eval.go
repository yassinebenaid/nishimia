package eval

import (
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
		return evalStatements(v.Statements)
	case *ast.ExpressionStatement:
		return Eval(v.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: v.Value}
	case *ast.BooleanLiteral:
		return nativeBooleanObject(v.Value)
	case *ast.PrefixExpression:
		return evalPrefixExpression(v.Operator, Eval(v.Right))
	case *ast.InfixExpression:
		return evalInfixExpression(v.Operator, Eval(v.Left), Eval(v.Right))
	case *ast.IfElseExpression:
		return evalConditionalExpression(v)
	case *ast.BlockStatement:
		return evalStatements(v.Statements)
	case *ast.ReturnStatement:
		return &object.ReturnValue{Value: Eval(v.Return)}
	}
	return NULL
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt)

		if returned, ok := result.(*object.ReturnValue); ok {
			return returned.Value
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
	default:
		return NULL // TODO : throw error
	}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {

	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		return evalIntegerInfixExpression(operator, left, right)
	}

	if left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ {
		return evalBooleanInfixExpression(operator, left, right)
	}

	return NULL // TODO : throw error

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
		return FALSE // TODO: register evaluation error here
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {

	if right.Type() != object.INTEGER_OBJ {
		return NULL // TODO: register error here
	}

	value := right.(*object.Integer).Value

	return &object.Integer{Value: value * -1}
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
		return NULL // TODO : throw error
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
		return NULL // TODO : throw error
	}
}

func evalConditionalExpression(node *ast.IfElseExpression) object.Object {
	cond := Eval(node.Condition)

	if cond.Type() != object.BOOLEAN_OBJ {
		return NULL // TODO : throw error here
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
