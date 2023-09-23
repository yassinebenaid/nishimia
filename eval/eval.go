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
	case *ast.StringLiteral:
		return &object.String{Value: v.Value}
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
	case *ast.ArrayLiteral:
		arr := &object.Array{Items: make([]object.Object, 0, len(v.Items))}

		for _, i := range v.Items {
			arr.Items = append(arr.Items, Eval(i, env))
		}

		return arr
	case *ast.HashLiteral:
		return evalHash(v, env)
	case *ast.IndexExpression:
		return evalIndexExression(v, env)
	case *ast.PrefixExpression:
		val := Eval(v.Right, env)
		if isError(val) {
			return val
		}

		return evalPrefixExpression(v.Operator, val)
	case *ast.Identifier:
		if val, ok := env.Get(v.Value); ok {
			return val
		}

		if val, ok := builtins[v.Value]; ok {
			return val
		}

		return newError("undefined identifier : %s", v.Value)

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

	return newError("unknown expression : %s", node.String())
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

	if left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ {
		return evalStringInfixExpression(operator, left, right)
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

func evalStringInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value

	if operator == "+" {
		return &object.String{Value: leftValue + rightValue}
	}

	return newError(
		"invalid operation: %s %s %s",
		leftValue,
		operator,
		rightValue,
	)

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

	args := evalExpressions(node.Arguments, env)
	if len(args) == 1 && isError(args[0]) {
		return args[0]
	}

	if fn, ok := function.(*object.Builtin); ok {
		return fn.Fn(args...)
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

func evalHash(hash *ast.HashLiteral, env *object.Environment) object.Object {
	h := &object.Hash{Items: make(map[object.HashKey]object.HashPair)}

	for k, v := range hash.Items {
		key := Eval(k, env)
		if isError(key) {
			return key
		}

		val := Eval(v, env)
		if isError(val) {
			return val
		}

		var hashKey object.HashKey
		if b, ok := key.(*object.Boolean); ok {
			hashKey = b.HashKey()
		} else if i, ok := key.(*object.Integer); ok {
			hashKey = i.HashKey()
		} else if s, ok := key.(*object.String); ok {
			hashKey = s.HashKey()
		} else {
			return newError("cannot use value of type %T as hash key", key)
		}

		h.Items[hashKey] = object.HashPair{
			Key:   key,
			Value: val,
		}
	}

	return h
}

func evalIndexExression(arr *ast.IndexExpression, env *object.Environment) object.Object {
	left := Eval(arr.Left, env)
	if isError(left) {
		return left
	}

	switch v := left.(type) {
	case *object.Array:
		return evalArrayIndexExression(v, arr.Index, env)
	case *object.Hash:
		return evalHashIndexExression(v, arr.Index, env)
	default:
		return newError("failed to read index on type %s", v.Type())
	}

}

func evalArrayIndexExression(array *object.Array, i ast.Expression, env *object.Environment) object.Object {

	ind := Eval(i, env)
	if isError(ind) {
		return ind
	}

	index, ok := ind.(*object.Integer)
	if !ok {
		return newError("cannot convert %s of type %s to type %s",
			ind.Inspect(),
			ind.Type(),
			object.INTEGER_OBJ,
		)
	}

	if index.Value >= int64(len(array.Items)) {
		return newError("index out of range [%d] with length %d",
			index.Value,
			len(array.Items),
		)
	}

	return array.Items[index.Value]
}

func evalHashIndexExression(hash *object.Hash, i ast.Expression, env *object.Environment) object.Object {

	ind := Eval(i, env)
	if isError(ind) {
		return ind
	}

	if ind.Type() != object.STRING_OBJ && ind.Type() != object.BOOLEAN_OBJ && ind.Type() != object.INTEGER_OBJ {
		return newError("invalid hash key  %s of type %s ",
			ind.Inspect(),
			ind.Type(),
		)
	}

	var hashKey object.HashKey
	if b, ok := ind.(*object.Boolean); ok {
		hashKey = b.HashKey()
	} else if i, ok := ind.(*object.Integer); ok {
		hashKey = i.HashKey()
	} else if s, ok := ind.(*object.String); ok {
		hashKey = s.HashKey()
	} else {
		return newError("cannot use value of type %T as hash key", ind)
	}

	result, ok := hash.Items[hashKey]
	if !ok {
		return newError("attempts to read undefined hash key [%s]", ind.Inspect())
	}

	return result.Value
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
