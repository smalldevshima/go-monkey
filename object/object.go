package object

import "fmt"

/// Constants / Variables

// Object types
const (
	O_NULL = "NULL"

	O_INTEGER = "INTEGER"
	O_BOOLEAN = "BOOLEAN"
)

// Object string formats
const (
	F_NULL = "null"

	F_BOOLEAN = "%v"
	F_INTEGER = "%d"
)

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
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%v", b.Value) }

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return O_INTEGER }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type Null struct{}

func (n *Null) Type() ObjectType { return O_NULL }
func (n *Null) Inspect() string  { return "null" }
