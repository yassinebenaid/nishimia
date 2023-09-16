package object

import "fmt"

type ObjectType string

const (
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	NULL_OBJ    = "NULL"
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

func (i *Integer) Type() string {
	return INTEGER_OBJ
}

type Boolean struct {
	Value bool
}

func (i *Boolean) Inspect() string {
	return fmt.Sprintf("%t", i.Value)
}

func (i *Boolean) Type() string {
	return BOOLEAN_OBJ
}

type Null struct{}

func (i *Null) Inspect() string {
	return NULL_OBJ
}

func (i *Null) Type() string {
	return "null"
}
