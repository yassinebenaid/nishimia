package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/yassinebenaid/nishimia/ast"
)

type ObjectType string

const (
	INTEGER_OBJ          ObjectType = "INTEGER"
	STRING_OBJ           ObjectType = "STRING"
	BOOLEAN_OBJ          ObjectType = "BOOLEAN"
	NULL_OBJ             ObjectType = "NULL"
	RETURN_VALUE_OBJ     ObjectType = "RETURN_VALUE"
	ERROR_OBJ            ObjectType = "ERROR"
	FUNCTION_OBJ         ObjectType = "FUNCTION"
	BUILTIN_FUNCTION_OBJ ObjectType = "BUILTIN_FUNCTION"
	ARRAY_OBJ            ObjectType = "ARRAY"
	HASH_OBJ             ObjectType = "HASH"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

type String struct {
	Value string
}

func (str *String) Inspect() string {
	return str.Value
}

func (str *String) Type() ObjectType {
	return STRING_OBJ
}

type Boolean struct {
	Value bool
}

func (i *Boolean) Inspect() string {
	return fmt.Sprintf("%t", i.Value)
}

func (i *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

type Null struct{}

func (i *Null) Inspect() string {
	return "null"
}

func (i *Null) Type() ObjectType {
	return NULL_OBJ
}

type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (r *ReturnValue) Inspect() string {
	return r.Value.Inspect()
}

type Error struct {
	Message string
}

func (*Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}

type Function struct {
	Name   string
	Params []*ast.Identifier
	Body   *ast.BlockStatement
	Env    *Environment
}

func (*Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer
	var params []string

	for _, p := range f.Params {
		params = append(params, p.String())
	}

	out.WriteString("func")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(f.Body.String())

	return out.String()
}

type BuiltinFunction func(...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (*Builtin) Type() ObjectType  { return FUNCTION_OBJ }
func (f *Builtin) Inspect() string { return "builtin function" }

type Array struct {
	Items []Object
}

func (*Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	var items []string

	for _, i := range a.Items {
		items = append(items, i.Inspect())
	}

	return "[" + strings.Join(items, ", ") + "]"
}

type Hash struct {
	Items map[Object]Object
}

func (*Hash) Type() ObjectType { return ARRAY_OBJ }
func (h *Hash) Inspect() string {
	var items []string

	for k, v := range h.Items {
		items = append(items, k.Inspect()+": "+v.Inspect())
	}

	return "{" + strings.Join(items, ", ") + "}"
}
