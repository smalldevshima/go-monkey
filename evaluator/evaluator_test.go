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
		// * literal value only
		{"literal/true", "true", true},
		{"literal/false", "false", false},

		// * bang prefix operator
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

		// * equality infix operator
		{"eq/literal/booleans/tt", "true == true", true},
		{"eq/literal/booleans/tf", "true == false", false},
		{"eq/literal/booleans/ft", "false == true", false},
		{"eq/literal/booleans/ff", "false == false", true},

		{"eq/literal/integers/same", "0 == 0", true},
		{"eq/literal/integers/different", "0 == 10", false},

		// {"eq/literal/bool-int/t-zero", "true == 0", false},
		// {"eq/literal/bool-int/f-zero", "false == 0", true},
		// {"eq/literal/bool-int/t-positive", "true == 5", true},
		// {"eq/literal/bool-int/f-positive", "false == 5", false},
		// {"eq/literal/bool-int/t-negative", "true == -5", true},
		// {"eq/literal/bool-int/f-negative", "false == -5", false},

		// * inequality infix operator
		{"neq/literal/booleans/tt", "true != true", false},
		{"neq/literal/booleans/tf", "true != false", true},
		{"neq/literal/booleans/ft", "false != true", true},
		{"neq/literal/booleans/ff", "false != false", false},

		{"neq/literal/integers/same", "0 != 0", false},
		{"neq/literal/integers/different", "0 != 10", true},

		// {"neq/literal/bool-int/t-zero", "true != 0", true},
		// {"neq/literal/bool-int/f-zero", "false != 0", false},
		// {"neq/literal/bool-int/t-positive", "true != 5", false},
		// {"neq/literal/bool-int/f-positive", "false != 5", true},
		// {"neq/literal/bool-int/t-negative", "true != -5", false},
		// {"neq/literal/bool-int/f-negative", "false != -5", true},

		// * less-then infix operator
		// {"less/literal/booleans/tt", "true < true", false},
		// {"less/literal/booleans/tf", "true < false", false},
		// {"less/literal/booleans/ft", "false < true", true},
		// {"less/literal/booleans/ff", "false < false", false},

		{"less/literal/integers/same", "0 < 0", false},
		{"less/literal/integers/lesser", "0 < 10", true},
		{"less/literal/integers/greater", "10 < 0", false},

		// * greater-then infix operator
		// {"greater/literal/booleans/tt", "true > true", false},
		// {"greater/literal/booleans/tf", "true > false", true},
		// {"greater/literal/booleans/ft", "false > true", false},
		// {"greater/literal/booleans/ff", "false > false", false},

		{"greater/literal/integers/same", "0 > 0", false},
		{"greater/literal/integers/lesser", "0 > 10", false},
		{"greater/literal/integers/greater", "10 > 0", true},
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
