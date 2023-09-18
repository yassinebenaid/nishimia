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

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch v := node.(type) {
	case *ast.Program:
		return evalProgram(v.Statements, env)
	case *ast.ExpressionStatement:
		return Eval(v.Expression, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: v.Value}
	case *ast.BooleanLiteral:
		return nativeBooleanObject(v.Value)
	case *ast.IfElseExpression:
		return evalConditionalExpression(v, env)
	case *ast.VarStatement:
		return evalVariableInitializationExpression(v, env)
	case *ast.BlockStatement:
		return evalBlockStatements(v, env)
	case *ast.FunctionLiteral:
		return &object.Function{Params: v.Params, Body: v.Body, Env: env}
	case *ast.CallExpression:
		return evalCallExpression(v, env)
	case *ast.PrefixExpression:
		val := Eval(v.Right, env)
		if isError(val) {
			return val
		}

		return evalPrefixExpression(v.Operator, val)
	case *ast.Identifier:
		val, ok := env.Get(v.Value)
		if !ok {
			return newError("undefined identifier : %s", v.Value)
		}

		return val
	case *ast.InfixExpression:
		l := Eval(v.Left, env)
		if isError(l) {
			return l
		}

		r := Eval(v.Right, env)
		if isError(r) {
			return r
		}

		return evalInfixExpression(v.Operator, l, r)

	case *ast.ReturnStatement:
		val := Eval(v.Return, env)
		if isError(val) {
			return val
		}

		return &object.ReturnValue{Value: val}
	}

	return NULL
}

func evalProgram(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatements(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range block.Statements {
		result = Eval(stmt, env)

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
	case "&&", "||":
		return newError(
			"invalid operation: %d %s %d , operator %s can only be used with BOOLEAN. got INTEGER",
			leftValue,
			operator,
			rightValue,
			operator,
		)
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

func evalConditionalExpression(node *ast.IfElseExpression, env *object.Environment) object.Object {
	cond := Eval(node.Condition, env)
	if isError(cond) {
		return cond
	}

	if cond.Type() != object.BOOLEAN_OBJ {
		return newError("non-boolean value in if-statement , ( got=%s, want=BOOLEAN )", cond.Type())
	}

	if cond.Inspect() == "true" {
		return Eval(node.Consequence, env)
	}

	if node.Alternative != nil {
		return Eval(node.Alternative, env)
	}

	return NULL
}

func evalVariableInitializationExpression(node *ast.VarStatement, env *object.Environment) object.Object {
	val := Eval(node.Value, env)
	if isError(val) {
		return val
	}

	if env.Has(node.Name.Value) {
		return newError("variable %s already defined", node.Name.Value)
	}

	env.Set(node.Name.Value, val)

	return NULL
}

func evalCallExpression(node *ast.CallExpression, env *object.Environment) object.Object {
	function := Eval(node.Function, env)
	if isError(function) {
		return function
	}

	fn, ok := function.(*object.Function)
	if !ok {
		return newError(
			"invalid identifier in function call : %s is not a valid identifier or function literal",
			function.Inspect(),
		)
	}

	if len(fn.Params) != len(node.Arguments) {
		return newError("invalid arguments count in function call, expected %d argumets, got %d ",
			len(fn.Params),
			len(node.Arguments),
		)
	}

	args := evalExpressions(node.Arguments, env)
	if len(args) == 1 && isError(args[0]) {
		return args[0]
	}

	newEnv := object.NewEnclosedEnvironment(fn.Env)
	for k, v := range fn.Params {
		newEnv.Set(v.Value, args[k])
	}

	result := Eval(fn.Body, newEnv)

	if result.Type() == object.RETURN_VALUE_OBJ {
		return result.(*object.ReturnValue).Value
	}

	if result.Type() == object.ERROR_OBJ {
		return result
	}

	return NULL
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, exp := range exps {
		evaluated := Eval(exp, env)

		if isError(evaluated) {
			return []object.Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
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
