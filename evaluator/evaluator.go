package evaluator

import (
	"github.com/smalldevshima/go-monkey/ast"
	"github.com/smalldevshima/go-monkey/object"
)

// Constants / Variables

// Error format strings
const (
	ERR_PREFIX_UNKNOWN ErrorFormat = "unknown operator: %s%s"
	ERR_INFIX_UNKNOWN ErrorFormat= "unknown operator: %s %s %s"
	ERR_INFIX_MISMATCH ErrorFormat= "type mismatch: %s %s %s"
)

var (
	NULL = &object.Null{}

	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}

	// FALSY_VALUES is a list of all object values considered falsy in Monkey
	FALSY_VALUES = []object.Object{NULL, FALSE}
)

// Functions

func newError(format ErrorFormat, a ...interface)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// * Statements:
	case *ast.Program:
		return evalProgram(node.Statements)
	case *ast.BlockStatement:
		return evalBlockStatement(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue)
		return &object.ReturnValue{Value: val}

	// * Literal expressions:
	case *ast.BooleanLiteral:
		return nativeBooleanToObject(node.Value)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	// * Operator expressions:
	case *ast.PrefixExpression:
		operand := Eval(node.Right)
		return evalPrefixExpression(node.Operator, operand)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpression(node.Operator, left, right)

	// * Control flow expressions:
	case *ast.IfExpression:
		return evalIfExpression(node)

	}

	return nil
}

func nativeBooleanToObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func evalProgram(statements []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range statements {
		result = Eval(stmt)

		// * return early, if result is an object.ReturnValue
		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue.Value
		}
	}

	return result
}

func evalBlockStatement(statements []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range statements {
		result = Eval(stmt)

		// * return early, if result type is object.O_RETURN_VALUE
		if result != nil && result.Type() == object.O_RETURN_VALUE {
			return result
		}
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

// evalBangOperatorExpression returns the opposite object of the isTruthy(operand) result
func evalBangOperatorExpression(operand object.Object) object.Object {
	if isTruthy(operand) {
		return FALSE
	}
	return TRUE
}

func evalDashOperatorExpression(operand object.Object) object.Object {
	if operand.Type() != object.O_INTEGER {
		return NULL
	}

	value := operand.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	// * need to switch on both the type of left and right
	case left.Type() == object.O_INTEGER && right.Type() == object.O_INTEGER:
		return evalIntegerInfixExpression(operator, left, right)

	// * special cases for infix operators '==' and '!='
	// * directly compare pointers, since booleans and null use global objects
	// * all other types are filtered out by preceding cases
	case operator == "==":
		return nativeBooleanToObject(left == right)
	case operator == "!=":
		return nativeBooleanToObject(left != right)
	}

	return NULL
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftInt := left.(*object.Integer).Value
	rightInt := right.(*object.Integer).Value
	var newInt int64
	switch operator {
	case "+":
		newInt = leftInt + rightInt
	case "-":
		newInt = leftInt - rightInt
	case "*":
		newInt = leftInt * rightInt
	case "/":
		newInt = leftInt / rightInt
	case "==":
		return nativeBooleanToObject(leftInt == rightInt)
	case "!=":
		return nativeBooleanToObject(leftInt != rightInt)
	case "<":
		return nativeBooleanToObject(leftInt < rightInt)
	case ">":
		return nativeBooleanToObject(leftInt > rightInt)
	default:
		return NULL
	}

	return &object.Integer{Value: newInt}
}

func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)

	if isTruthy(condition) {
		return Eval(ie.Then)
	} else if ie.Otherwise != nil {
		return Eval(ie.Otherwise)
	}

	return NULL
}

// isTruthy defines which values are truthy in the Monkey language
func isTruthy(obj object.Object) bool {
	for _, falsyVal := range FALSY_VALUES {
		if falsyVal == obj {
			return false
		}
	}

	return true
}

/// Types

type ErrorFormat string
