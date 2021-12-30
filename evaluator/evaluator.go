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

	// * Literal expressions:
	case *ast.BooleanLiteral:
		return nativeBooleanToObject(node.Value)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

		// * Operator expressions:
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixExpression(node.Operator, right)
	}

	return nil
}
func nativeBooleanToObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalStatements(statements []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range statements {
		result = Eval(stmt)
	}

	return result
}

func evalPrefixExpression(operator string, operand object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(operand)
	case "-":
		return evalDashOperatorExpression(operand)
	}
	return NULL
}

// evalBangOperatorExpression checks all types of falsy values explicitly.
// Otherwise it assumes that operand is truthy and returns the false-object.
func evalBangOperatorExpression(operand object.Object) object.Object {
	switch operand := operand.(type) {
	case *object.Null:
		return TRUE
	case *object.Boolean:
		if operand == FALSE {
			return TRUE
		}
	case *object.Integer:
		if operand.Value == 0 {
			return TRUE
		}
	}

	return FALSE
}

func evalDashOperatorExpression(operand object.Object) object.Object {
	if operand.Type() != object.O_INTEGER {
		return NULL
	}

	value := operand.(*object.Integer).Value
	return &object.Integer{Value: -value}
}
