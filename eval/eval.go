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
	}
	return NULL
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt)
	}

	return result
}

func nativeBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}

	return FALSE
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
	switch operator {
	case "+":
		return &object.Integer{Value: left.(*object.Integer).Value + right.(*object.Integer).Value}
	case "-":
		return &object.Integer{Value: left.(*object.Integer).Value - right.(*object.Integer).Value}
	case "*":
		return &object.Integer{Value: left.(*object.Integer).Value * right.(*object.Integer).Value}
	case "/":
		return &object.Integer{Value: left.(*object.Integer).Value / right.(*object.Integer).Value}
	case "<":
		return &object.Boolean{Value: left.(*object.Integer).Value < right.(*object.Integer).Value}
	case ">":
		return &object.Boolean{Value: left.(*object.Integer).Value > right.(*object.Integer).Value}
	case ">=":
		return &object.Boolean{Value: left.(*object.Integer).Value >= right.(*object.Integer).Value}
	case "<=":
		return &object.Boolean{Value: left.(*object.Integer).Value <= right.(*object.Integer).Value}
	case "!=":
		return &object.Boolean{Value: left.(*object.Integer).Value != right.(*object.Integer).Value}
	case "==":
		return &object.Boolean{Value: left.(*object.Integer).Value == right.(*object.Integer).Value}
	default:
		return NULL // TODO : throw error
	}
}

func evalBooleanInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch operator {
	case "&&":
		return &object.Boolean{Value: left.(*object.Boolean).Value && right.(*object.Boolean).Value}
	case "||":
		return &object.Boolean{Value: left.(*object.Boolean).Value || right.(*object.Boolean).Value}
	case "==":
		return &object.Boolean{Value: left.(*object.Boolean).Value == right.(*object.Boolean).Value}
	case "!=":
		return &object.Boolean{Value: left.(*object.Boolean).Value != right.(*object.Boolean).Value}
	default:
		return NULL // TODO : throw error
	}
}