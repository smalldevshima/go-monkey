package parser

import (
	"fmt"
	"testing"

	"github.com/smalldevshima/go-monkey/ast"
	"github.com/smalldevshima/go-monkey/lexer"
	"github.com/smalldevshima/go-monkey/token"
)

func TestLetStatements(t *testing.T) {
	input := `
		let x = 5;
		let y = 10;
		let foobar = 838383;
		let foo = bar;
		`

	p := New(lexer.New(input))

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 4 {
		t.Fatalf("program.Statements does not contain 4 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"x", 5},
		{"y", 10},
		{"foobar", 838383},
		{"foo", "bar"},
	}

	for index, test := range tests {
		stmt := program.Statements[index]
		ok := t.Run("statement"+fmt.Sprint(index+1), func(tt *testing.T) {
			checkLetStatement(tt, stmt, test.expectedIdentifier, test.expectedValue)
		})
		if !ok {
			t.Fail()
		}
	}
}

func TestReturnStatements(t *testing.T) {
	returnTests := []struct {
		name              string
		input             string
		returnValueString string
	}{
		{
			"integerLiteral",
			"return 5;",
			"5",
		},
		{
			"booleanLiteral",
			"return true;",
			"true",
		},
		{
			"infixExpression",
			"return x + z;",
			"(x + z)",
		},
		{
			"functionCall",
			"return add(5, 10);",
			"add(5, 10)",
		},
	}
	for _, test := range returnTests {
		t.Run("return/"+test.name, func(tt *testing.T) {
			l := lexer.New(test.input)
			p := New(l)

			program := p.ParseProgram()
			checkParserErrors(tt, p)
			if len(program.Statements) != 1 {
				tt.Fatalf("program.Statements does not contain 1 statement. got=%d: %s", len(program.Statements), program.Statements)
			}
			stmt, ok := program.Statements[0].(*ast.ReturnStatement)
			if !ok {
				tt.Fatalf("program.Statements[0] is not *ast.ReturnStatement. got=%T", program.Statements[0])
			}
			if stmt.TokenLiteral() != "return" {
				tt.Errorf("stmt.TokenLiteral is not 'return', got=%q", stmt.TokenLiteral())
			}
			if stmt.ReturnValue.String() != test.returnValueString {
				tt.Fatalf("stmt.ReturnValue.String is wrong.\nexpected:\n\t%s\ngot:\n\t%s", test.returnValueString, stmt.ReturnValue.String())
			}
		})
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := `foobar;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.Identifier. got=%T", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value is not 'foobar'. got=%q", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral is not 'foobar'. got=%q", ident.TokenLiteral())
	}
}

func TestLiteralExpression(t *testing.T) {
	literalTests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"5;", 5},
		{"true;", true},
		{"false;", false},
	}

	for _, test := range literalTests {
		t.Run("literal/"+fmt.Sprint(test.expectedValue), func(tt *testing.T) {
			l := lexer.New(test.input)
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(tt, p)

			if len(program.Statements) != 1 {
				tt.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
			}
			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				tt.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
			}
			checkLiteralExpression(tt, stmt.Expression, test.expectedValue)
		})
	}
}

func TestPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, test := range prefixTests {
		t.Run("prefix"+test.operator, func(tt *testing.T) {
			l := lexer.New(test.input)
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(tt, p)

			if len(program.Statements) != 1 {
				tt.Fatalf("program.Statements does not contain 1 statement. got=%d: %s", len(program.Statements), program.Statements)
			}

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				tt.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
			}

			exp, ok := stmt.Expression.(*ast.PrefixExpression)
			if !ok {
				t.Fatalf("stmt.Expression is not *ast.PrefixExpression. got=%T", stmt.Expression)
			}
			if exp.Operator != test.operator {
				t.Fatalf("exp.Operator is not %q. got=%q", test.operator, exp.Operator)
			}
			checkLiteralExpression(tt, exp.Right, test.value)
		})
	}
}

func TestInfixExpression(t *testing.T) {
	infixTests := []struct {
		input    string
		left     interface{}
		operator string
		right    interface{}
	}{
		{"1 + 2", 1, "+", 2},
		{"3 - 4", 3, "-", 4},
		{"5 * 6", 5, "*", 6},
		{"7 / 8", 7, "/", 8},
		{"9 > 10", 9, ">", 10},
		{"11 < 12", 11, "<", 12},
		{"13 == 14", 13, "==", 14},
		{"15 != 16", 15, "!=", 16},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, test := range infixTests {
		t.Run("infix"+test.operator, func(tt *testing.T) {
			l := lexer.New(test.input)
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(tt, p)
			if len(program.Statements) != 1 {
				tt.Fatalf("program.Statements does not contain 1 statement. got=%d: %s", len(program.Statements), program.Statements)
			}

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				tt.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
			}
			checkInfixExpression(tt, stmt.Expression, test.left, test.operator, test.right)
		})
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b);",
		},
		{
			"!-a",
			"(!(-a));",
		},
		{
			"a + b + c",
			"((a + b) + c);",
		},
		{
			"a + b - c",
			"((a + b) - c);",
		},
		{
			"a * b * c",
			"((a * b) * c);",
		},
		{
			"a * b / c",
			"((a * b) / c);",
		},
		{
			"a + b / c",
			"(a + (b / c));",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f);",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4);((-5) * 5);",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4));",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4));",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)));",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4);",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2);",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5));",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5));",
		},
		{
			"!(true == true)",
			"(!(true == true));",
		},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("\nexpected=\n\t%q\ngot=\n\t%q", tt.expected, actual)
		}
	}
}

func TestIfExpression(t *testing.T) {
	ifTests := []struct {
		input     string
		condLeft  interface{}
		condOp    string
		condRight interface{}
		thenName  string
		// leave empty for no else-branch
		otherwiseName string
	}{
		{
			`if (x < y) { a } else { b }`,
			"x", "<", "y",
			"a", "b",
		},
		{
			`if (previous == current) { current }`,
			"previous", "==", "current",
			"current", "",
		},
	}

	for index, test := range ifTests {
		t.Run("if-else/"+fmt.Sprint(index), func(tt *testing.T) {
			l := lexer.New(test.input)
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(tt, p)
			if len(program.Statements) != 1 {
				tt.Fatalf("program.Statements does not contain 1 statement. got=%d: %s", len(program.Statements), program.Statements)
			}

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				tt.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
			}

			exp, ok := stmt.Expression.(*ast.IfExpression)
			if !ok {
				tt.Fatalf("stmt.Expression is not *ast.IfExpression. got=%T", stmt.Expression)
			}

			checkInfixExpression(tt, exp.Condition, test.condLeft, test.condOp, test.condRight)

			if len(exp.Then.Statements) != 1 {
				tt.Fatalf("then-block does not contain 1 statement. got=%d: %s", len(exp.Then.Statements), exp.Then.Statements)
			}
			then, ok := exp.Then.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				tt.Fatalf("exp.Then.Statements[0] is not *ast.ExpressionStatement. got=%T", exp.Then.Statements[0])
			}
			checkLiteralExpression(tt, then.Expression, test.thenName)

			if test.otherwiseName == "" {
				return
			}

			if exp.Otherwise == nil {
				tt.Fatalf("else-block is nil")
			}
			if len(exp.Otherwise.Statements) != 1 {
				tt.Fatalf("else-block does not contain 1 statement. got=%d: %s", len(exp.Otherwise.Statements), exp.Otherwise.Statements)
			}
			otherwise, ok := exp.Otherwise.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				tt.Fatalf("exp.Otherwise.Statements[0] is not *ast.ExpressionStatement. got=%T", exp.Otherwise.Statements[0])
			}

			checkLiteralExpression(tt, otherwise.Expression, test.otherwiseName)
		})
	}
}

func TestFunctionLiteral(t *testing.T) {
	fnTests := []struct {
		input   string
		params  []string
		bodyLen int
	}{
		{
			`fn(x, y) { x + y; }`,
			[]string{"x", "y"},
			1,
		},
		{
			`fn() { x; y; z; }`,
			[]string{},
			3,
		},
		{
			`fn(a,b,c) {}`,
			[]string{"a", "b", "c"},
			0,
		},
		{
			`fn(x) { x; }`,
			[]string{"x"},
			1,
		},
	}
	for index, test := range fnTests {
		t.Run("functionLiteral"+fmt.Sprint(index), func(tt *testing.T) {
			l := lexer.New(test.input)
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(tt, p)
			if len(program.Statements) != 1 {
				tt.Fatalf("program.Statements does not contain 1 statement. got=%d: %s", len(program.Statements), program.Statements)
			}

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				tt.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
			}
			fn, ok := stmt.Expression.(*ast.FunctionLiteral)
			if !ok {
				tt.Fatalf("stmt.Expression is not *ast.FunctionLiteral. got=%T", stmt.Expression)
			}

			if len(fn.Parameters) != len(test.params) {
				tt.Fatalf("fn.Parameters does not contain %d identifiers. got=%d", len(test.params), len(fn.Parameters))
			}
			for index, param := range fn.Parameters {
				checkIdentifier(tt, param, test.params[index])
			}

			if len(fn.Body.Statements) != test.bodyLen {
				tt.Fatalf("fn.Body.Statements does not contain %d statements. got=%d", test.bodyLen, len(fn.Body.Statements))
			}
		})
	}
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}
	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)
		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want %d, got=%d\n", len(tt.expectedParams), len(function.Parameters))
		}
		for i, ident := range tt.expectedParams {
			checkLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestFunctionCallExpression(t *testing.T) {
	callTests := []struct {
		name      string
		input     string
		function  string
		arguments []string
	}{
		{
			"identifier/no-arg",
			"next()",
			"next",
			[]string{},
		},
		{
			"identifier/one-arg",
			"incr(3 + 3 * 2)",
			"incr",
			[]string{"(3 + (3 * 2))"},
		},
		{
			"identifier/mutli-arg",
			"add(5, 30 / 5, 20)",
			"add",
			[]string{"5", "(30 / 5)", "20"},
		},
		{
			"literal/no-arg",
			"fn(){ return 2 + 2 }()",
			"fn () { return (2 + 2); }",
			[]string{},
		},
		{
			"literal/one-arg",
			"fn(x){ x * 2 }(10)",
			"fn (x) { (x * 2); }",
			[]string{"10"},
		},
		{
			"literal/multi-arg",
			"fn(a,b,c){ a<b == a<c } (5*2,5,20)",
			"fn (a, b, c) { ((a < b) == (a < c)); }",
			[]string{"(5 * 2)", "5", "20"},
		},
		{
			"if-else/no-arg",
			"if(neg){incr}else{decr}()",
			"if (neg) { incr; } else { decr; }",
			[]string{},
		},
	}
	for _, test := range callTests {
		t.Run("call/"+test.name, func(tt *testing.T) {
			l := lexer.New(test.input)
			p := New(l)
			program := p.ParseProgram()
			checkParserErrors(tt, p)
			if len(program.Statements) != 1 {
				tt.Fatalf("program.Statements does not have 1 statement. got=%d: %q", len(program.Statements), program.Statements)
			}

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				tt.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
			}

			exp, ok := stmt.Expression.(*ast.CallExpression)
			if !ok {
				tt.Fatalf("stmt.Expression is not *ast.CallExpression. got=%T", stmt.Expression)
			}

			if exp.Function.String() != test.function {
				tt.Errorf("exp.Function.String is wrong.\nexpected:\n\t%s\ngot:\n\t%s", test.function, exp.Function)
			}

			if len(exp.Arguments) != len(test.arguments) {
				tt.Fatalf("exp.Arguments does not contain %d arguments. got=%d: %q", len(test.arguments), len(exp.Arguments), exp.Arguments)
			}
			for index, arg := range exp.Arguments {
				if arg.String() != test.arguments[index] {
					tt.Errorf("exp.Arguments[%d] is wrong.\nexpected:\n\t%s\ngot:\n\t%s", index, test.arguments[index], arg)
				}
			}
		})
	}
}

/// helpers

func checkParserErrors(t *testing.T, p *Parser) {
	t.Helper()
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors:", len(errors))
	for i, msg := range errors {
		if i >= 10 {
			t.Errorf("omitting more errors ...")
			break
		}
		t.Errorf("%3d: %s", i+1, msg)
	}
	t.FailNow()
}

func checkLetStatement(t *testing.T, s ast.Statement, name string, value interface{}) {
	t.Helper()
	if s.TokenLiteral() != "let" {
		t.Fatalf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Fatalf("s not *ast.LetStatement. got=%T", s)
	}
	checkIdentifier(t, letStmt.Name, name)
	checkLiteralExpression(t, letStmt.Value, value)
}

func checkIntegerLiteral(t *testing.T, exp ast.Expression, value int64) {
	t.Helper()
	intLit, ok := exp.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("exp is not *ast.IntegerLiteral. got=%T", exp)
		return
	}
	if intLit.Token.Type != token.INTEGER {
		t.Errorf("intLit.Token.Type is not %q. got=%q", token.INTEGER, intLit.Token.Type)
	}
	if intLit.Value != value {
		t.Errorf("intLit.Value is not %d. got=%d", value, intLit.Value)
	}
}

func checkBooleanLiteral(t *testing.T, exp ast.Expression, value bool) {
	t.Helper()
	boolLit, ok := exp.(*ast.BooleanLiteral)
	if !ok {
		t.Errorf("exp is not *ast.BooleanLiteral. got=%T", exp)
		return
	}
	if boolLit.Token.Type != token.TRUE && boolLit.Token.Type != token.FALSE {
		t.Errorf("boolLit.Token.Type is neither %q nor %q. got=%q", token.TRUE, token.FALSE, boolLit.Token.Type)
	}
	if boolLit.Value != value {
		t.Errorf("boolLit.Value is not %v. got=%v", value, boolLit.Value)
	}
}

func checkIdentifier(t *testing.T, exp ast.Expression, value string) {
	t.Helper()
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
	}
	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
	}
	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value, ident.TokenLiteral())
	}
}

func checkLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) {
	t.Helper()
	switch v := expected.(type) {
	case int:
		checkIntegerLiteral(t, exp, int64(v))
	case int64:
		checkIntegerLiteral(t, exp, v)
	case string:
		checkIdentifier(t, exp, v)
	case bool:
		checkBooleanLiteral(t, exp, v)
	default:
		t.Errorf("type of expected not handled. got=%T", expected)
	}
}

func checkInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) {
	t.Helper()
	infixExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not *ast.InfixExpression. got=%T(%s)", exp, exp)
		return
	}
	checkLiteralExpression(t, infixExp.Left, left)
	if infixExp.Operator != operator {
		t.Errorf("infixExp.Operator is not %q. got=%q", operator, infixExp.Operator)
	}
	checkLiteralExpression(t, infixExp.Right, right)
}
