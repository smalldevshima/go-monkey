package evaluator

import (
	"fmt"
	"testing"

	"github.com/smalldevshima/go-monkey/lexer"
	"github.com/smalldevshima/go-monkey/object"
	"github.com/smalldevshima/go-monkey/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{"literal/0", "0", 0},
		{"literal/123456", "123456", 123456},
	}

	for index, test := range tests {
		t.Run("integer/"+fmt.Sprint(index), func(t *testing.T) {
			evaluated := testEval(test.input)
			checkIntegerObject(t, evaluated, test.expected)
		})
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"literal/true", "true", true},
		{"literal/false", "false", false},
	}

	for index, test := range tests {
		t.Run("boolean/"+fmt.Sprint(index), func(t *testing.T) {
			evaluated := testEval(test.input)
			checkBooleanObject(t, evaluated, test.expected)
		})
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"literal/true", "!true", false},
		{"literal/false", "!false", true},
		{"literal/zero-int", "!0", true},
		{"literal/non-zero-int", "!5", false},
		{"literal/twice/true", "!!true", true},
		{"literal/twice/false", "!!false", false},
		{"literal/twice/zero-int", "!!0", false},
		{"literal/twice/non-zero-int", "!!5", true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			evaluated := testEval(test.input)
			checkBooleanObject(t, evaluated, test.expected)
		})
	}
}

/// helpers

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	return Eval(program)
}

func checkIntegerObject(t *testing.T, obj object.Object, value int64) {
	integer, ok := obj.(*object.Integer)
	if !ok {
		t.Fatalf("obj is not *object.Integer. got=%T: (%+v)", obj, obj)
	}

	if integer.Value != value {
		t.Errorf("integer.Value is wrong. expected=%q, got=%q", integer.Value, value)
	}
	if integer.Inspect() != fmt.Sprintf(object.F_INTEGER, value) {
		t.Errorf("integer.Inspect is wrong. expected=%q, got=%q", fmt.Sprintf(object.F_INTEGER, value), integer.Inspect())
	}
}

func checkBooleanObject(t *testing.T, obj object.Object, value bool) {
	boolean, ok := obj.(*object.Boolean)
	if !ok {
		t.Fatalf("obj is not *object.Boolean. got=%T: (%+v)", obj, obj)
	}

	if boolean.Value != value {
		t.Errorf("boolean.Value is wrong. expected=%v, got=%v", boolean.Value, value)
	}
	if boolean.Inspect() != fmt.Sprintf(object.F_BOOLEAN, value) {
		t.Errorf("boolean.Inspect is wrong. expected=%v, got=%v", fmt.Sprintf(object.F_BOOLEAN, value), boolean.Inspect())
	}
}
