package object

import "fmt"

/// Constants / Variables

const typeDelim = "@"

// Object types

var (
	O_NULL ObjectType = typeString("null")

	O_INTEGER ObjectType = typeString("int")
	O_BOOLEAN ObjectType = typeString("bool")

	O_RETURN_VALUE ObjectType = typeString("return_value")

	O_ERROR = typeString("error")
)

// Object string formats
const (
	F_NULL = "null"

	F_BOOLEAN = "%v"
	F_INTEGER = "%d"

	F_RETURN_VALUE = "%v"

	F_ERROR = "ERROR: %s"
)

/// Functions

func typeString(typ string) ObjectType {
	return ObjectType(fmt.Sprintf("%s%s%s", typeDelim, typ, typeDelim))
}

/// Types

type ObjectType string

// Object is the base interface for all values in Monkey
type Object interface {
	// The Type of the Monkey object
	Type() ObjectType
	// Inspect returns the string representation of the Monkey object
	Inspect() string
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return O_BOOLEAN }
func (b *Boolean) Inspect() string  { return fmt.Sprintf(F_BOOLEAN, b.Value) }

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return O_INTEGER }
func (i *Integer) Inspect() string  { return fmt.Sprintf(F_INTEGER, i.Value) }

type Null struct{}

func (n *Null) Type() ObjectType { return O_NULL }
func (n *Null) Inspect() string  { return F_NULL }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return O_RETURN_VALUE }
func (rv *ReturnValue) Inspect() string  { return fmt.Sprintf(F_RETURN_VALUE, rv.Value.Inspect()) }

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return O_ERROR }
func (e *Error) Inspect() string  { return fmt.Sprintf(F_ERROR, e.Message) }
