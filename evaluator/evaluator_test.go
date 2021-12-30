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