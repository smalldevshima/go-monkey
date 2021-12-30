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
		{"literal/zero", "0", 0},
		{"literal/positive", "123456", 123456},

		{"literal/negation/zero", "-0", 0},
		{"literal/negation/non-zero", "-123456", -123456},
		{"literal/negation/twice/zero", "--0", 0},
		{"literal/negation/twice/non-zero", "--123456", 123456},

		{"literal/sum", "1 + 2", 3},
		{"literal/difference", "5 - 4", 1},
		{"literal/product", "8 * 8", 64},
		{"literal/division", "30 / 3", 10},

		{"literal/grouped/simple", "(1 + 2) * 5", 15},
		{"literal/grouped/complex", "3 + (3 - (1 + 2) * 5 - (-10 / (90 - 88)))", -4},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
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

		{"bang/literal/true", "!true", false},
		{"bang/literal/false", "!false", true},
		{"bang/literal/zero-int", "!0", true},
		{"bang/literal/neg-zero-int", "!-0", true},
		{"bang/literal/positive-int", "!5", false},
		{"bang/literal/negative-int", "!-5", false},

		{"bang/twice/literal/true", "!!true", true},
		{"bang/twice/literal/false", "!!false", false},
		{"bang/twice/literal/zero-int", "!!0", false},
		{"bang/twice/literal/neg-zero-int", "!!-0", false},
		{"bang/twice/literal/positive-int", "!!5", true},
		{"bang/twice/literal/negative-int", "!!-5", true},
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
	t.Helper()
	integer, ok := obj.(*object.Integer)
	if !ok {
		t.Fatalf("obj is not *object.Integer. got=%T: (%+v)", obj, obj)
	}

	if integer.Value != value {
		t.Errorf("integer.Value is wrong. expected=%d, got=%d", value, integer.Value)
	}
	if integer.Inspect() != fmt.Sprintf(object.F_INTEGER, value) {
		t.Errorf("integer.Inspect is wrong. expected=%q, got=%q", fmt.Sprintf(object.F_INTEGER, value), integer.Inspect())
	}
}

func checkBooleanObject(t *testing.T, obj object.Object, value bool) {
	t.Helper()
	boolean, ok := obj.(*object.Boolean)
	if !ok {
		t.Fatalf("obj is not *object.Boolean. got=%T: (%+v)", obj, obj)
	}

	if boolean.Value != value {
		t.Errorf("boolean.Value is wrong. expected=%v, got=%v", value, boolean.Value)
	}
	if boolean.Inspect() != fmt.Sprintf(object.F_BOOLEAN, value) {
		t.Errorf("boolean.Inspect is wrong. expected=%v, got=%v", fmt.Sprintf(object.F_BOOLEAN, value), boolean.Inspect())
	}
}
