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
		{"bang/literal/zero-int", "!0", false},
		{"bang/literal/neg-zero-int", "!-0", false},
		{"bang/literal/positive-int", "!5", false},
		{"bang/literal/negative-int", "!-5", false},

		{"bang/twice/literal/true", "!!true", true},
		{"bang/twice/literal/false", "!!false", false},
		{"bang/twice/literal/zero-int", "!!0", true},
		{"bang/twice/literal/neg-zero-int", "!!-0", true},
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

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{"no-else/boolean", "if (true) {10==10}", true},
		{"no-else/integer", "if (true) {10}", 10},
		{"no-else/null", "if (false) {10}", nil},

		{"if-else/boolean", "if (false) {10} else {10!=5}", true},
		{"if-else/integer", "if (false) {10==10} else {5}", 5},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			evaluated := testEval(test.input)
			if integer, ok := test.expected.(int); ok {
				checkIntegerObject(t, evaluated, int64(integer))
			} else if boolean, ok := test.expected.(bool); ok {
				checkBooleanObject(t, evaluated, boolean)
			} else {
				checkNullObject(t, evaluated)
			}
		})
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{"only-return/integer", "return 10;", 10},
		{"block/first-pos/integer", "return 9; 5;", 9},
		{"block/middle-pos", "7; return 1; 34;", 1},
		{"block/last-pos", "23; 43; return -24;", -24},

		{"block/if-else", "if (true) { return 100; } else { return 200; }; 300;", 100},
		{"block/function-call", "fn() { return 213; } ()", 213},

		{
			"block/nested/if-else",
			`if (true) {
				if (true) {
					return 456;
				}
				return 123;
			}`,
			456,
		},
		{
			"block/nested/function-call",
			`fn() {
				fn () {
					return 789;
				} ();
				return 987
			} ()`,
			987,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			evaluated := testEval(test.input)
			checkIntegerObject(t, evaluated, test.expected)
		})
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		message string
	}{
		{
			"operator/type/mismatch/int-bool",
			"5 + true;",
			"type mismatch: @int@ + @bool@",
		},
		{
			"operator/type/mismatch/bool-int",
			"true - 4;",
			"type mismatch: @bool@ - @int@",
		},

		{
			"identifier/unknown",
			"hello",
			"unknown identifier: hello",
		},

		{
			"operator/type/unknown/negate-bool",
			"-true;",
			"unknown operator: -@bool@",
		},
		{
			"operator/type/unknown/sum-bool",
			"true + false;",
			"unknown operator: @bool@ + @bool@",
		},

		{
			"block/exit-early",
			"-true; false; 1234;",
			"unknown operator: -@bool@",
		},
		{
			"block/if-else/exit-early",
			`if (true) {
				-true;
				10;
			}`,
			"unknown operator: -@bool@",
		},
		{
			"block/function-call/exit-early",
			`fn () {
				-false;
				return 30;
			} ()`,
			"unknown operator: -@bool@",
		},
		{
			"block/nested/if-else/exit-early",
			`if (true) {
				if (true) {
					true * true;
				}
				5432;
			}`,
			"unknown operator: @bool@ * @bool@",
		},
		{
			"block/nested/function-call/exit-early",
			`fn () {
				fn () {
					true / true;
				} ()
				return 5432;
			} ()`,
			"unknown operator: @bool@ / @bool@",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			evaluated := testEval(test.input)
			checkErrorObject(t, evaluated, test.message)
		})
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			"bind-literal/integer",
			"let value = 10; value;",
			10,
		},
		{
			"bind-literal/boolean",
			"let some = true; some;",
			true,
		},
		{
			"bind-literal/function",
			"let func = fn () { 10; }; func();",
			10,
		},

		{
			"bind-variable/integer",
			"let first = 10; let second = first; second;",
			10,
		},
		{
			"bind-variable/integer",
			"let first = false; let second = first; second;",
			false,
		},

		{
			"bind-expression/integer",
			"let left = 23; let right = 31; let result = left + right * 10; result;",
			333,
		},
		{
			"bind-expression/boolean",
			"let left = true; let right = false; let result = !left == right != false; result;",
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			evaluated := testEval(test.input)
			switch exp := test.expected.(type) {
			case int:
				checkIntegerObject(t, evaluated, int64(exp))
			case bool:
				checkBooleanObject(t, evaluated, exp)
			default:
				t.Fatalf("unknown expected type %T=%v", exp, exp)
			}
		})
	}
}

func TestFunctionObject(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		params []string
		body   string
	}{
		{
			"no-param/single-line",
			"fn() { return 2; }",
			[]string{},
			"return 2;",
		},
		{
			"one-param/single-line",
			"fn(x) { x + 1; }",
			[]string{"x"},
			"(x + 1);",
		},
		{
			"multi-param/single-line",
			"fn(x, y, z) { x - y - z }",
			[]string{"x", "y", "z"},
			"((x - y) - z);",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			evaluated := testEval(test.input)
			fn, ok := evaluated.(*object.Function)
			if !ok {
				t.Fatalf("evaluated is not *object.Function. got=%T: %+v", evaluated, evaluated)
			}

			if len(fn.Parameters) != len(test.params) {
				t.Fatalf("fn.Parameters does not contain %d params. got=%d", len(test.params), len(fn.Parameters))
			}
			for index, param := range fn.Parameters {
				if param.String() != test.params[index] {
					t.Errorf("fn.Paramerts[%d].String is wrong.\nexpected:\n\t%s\ngot:\n\t%s", index, test.params[index], param)
				}
			}

			if fn.Body.String() != test.body {
				t.Fatalf("fn.Body.String is wrong.\nexpected:\n>>>\n%s\n<<<\ngot:\n>>>\n%s\n<<<", test.body, fn.Body)
			}
		})
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			"no-param/return-literal",
			"let five = fn() { 5 }; five();",
			5,
		},
		{
			"no-param/return-captured",
			"let x = 10; let getX = fn() { x }; getX();",
			10,
		},
		{
			"no-param/shadowed-identifier",
			"let x = 3; fn() { let x = 11; } (); x;",
			3,
		},

		{
			"one-param/return-identity",
			"let ident = fn(x) {x}; ident(1234);",
			1234,
		},
		{
			"one-param/return-shadowing-param",
			"let shadow = 10; let f = fn(shadow) { shadow }; f(42);",
			42,
		},

		{
			"multi-param/return-curried",
			`
			let curriedSum = fn(a) {
				return fn(b) { return second(a,b) };
			}
			let second = fn(a, b) {
				return fn(c) { return third(a,b,c)};
			}
			let third = fn (a,b,c) {
				return a+b+c;
			}
			curriedSum(1)(2)(3)
			`,
			6,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			evaluated := testEval(test.input)
			switch exp := test.expected.(type) {
			case int:
				checkIntegerObject(t, evaluated, int64(exp))
			case bool:
				checkBooleanObject(t, evaluated, exp)
			default:
				t.Fatalf("unknown expected type %T=%v", exp, exp)
			}
		})
	}
}

/// helpers

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
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

func checkNullObject(t *testing.T, obj object.Object) {
	t.Helper()
	null, ok := obj.(*object.Null)
	if !ok {
		t.Fatalf("obj is not *object.Null. got=%T: (%+v)", obj, obj)
	}

	if null.Inspect() != object.F_NULL {
		t.Errorf("null.Inspect is wrong. expected=%v, got=%v", object.F_NULL, null.Inspect())
	}
}

func checkErrorObject(t *testing.T, obj object.Object, message string) {
	t.Helper()
	err, ok := obj.(*object.Error)
	if !ok {
		t.Fatalf("obj is not *object.Error. got=%T: (%+v)", obj, obj)
	}

	if err.Message != message {
		t.Fatalf("err.Message is wrong.\nexpected:\n\t%s\ngot:\n\t%s", message, err.Message)
	}
	if err.Inspect() != fmt.Sprintf(object.F_ERROR, message) {
		t.Fatalf("err.Inspect is wrong.\nexpected:\n\t%s\ngot:\n\t%s", fmt.Sprintf(object.F_ERROR, message), err.Inspect())
	}
}
