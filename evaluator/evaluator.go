package evaluator

import (
	"github.com/smalldevshima/go-monkey/ast"
	"github.com/smalldevshima/go-monkey/object"
)

// Constants / Variables

var (
	NULL = &object.Null{}

	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// Functions

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// * Statements:
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	// * Expressions:
	case *ast.BooleanLiteral:
		return nativeBooleanToObject(node.Value)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	}

	return nil
}

func evalStatements(statements []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range statements {
		result = Eval(stmt)
	}

	return result
}

func nativeBooleanToObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}
