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
		`

	p := New(lexer.New(input))

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for index, test := range tests {
		stmt := program.Statements[index]
		ok := t.Run("statement"+fmt.Sprint(index+1), func(tt *testing.T) {
			checkLetStatement(tt, stmt, test.expectedIdentifier)
		})
		if !ok {
			t.Fail()
		}
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
		return 5;
		return 10;
		return add(5, 10);
		return x + z;
		`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 4 {
		t.Fatalf("program.Statements does not contain 4 statements. got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("smt not *ast.ReturnStatement. got=%T", stmt)
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', got=%q", returnStmt.TokenLiteral())
		}
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

func TestIntergerLiteralExpression(t *testing.T) {
	input := "5;"

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

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value is not 5. got=%d", literal.Value)
	}
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral is not '5'. got=%q", literal.TokenLiteral())
	}
}

func TestPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
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
			checkIntegerLiteral(tt, exp.Right, test.integerValue)
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
		t.Errorf("%2d: %s", i, msg)
	}
	t.FailNow()
}

func checkLetStatement(t *testing.T, s ast.Statement, name string) {
	t.Helper()
	if s.TokenLiteral() != "let" {
		t.Fatalf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Fatalf("s not *ast.LetStatement. got=%T", s)
	}

	if letStmt.Name.Value != name {
		t.Fatalf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Fatalf("letStmt.Name.TokenLiteral() not '%s'. got=%s", name, letStmt.Name.TokenLiteral())
	}
}

func checkIntegerLiteral(t *testing.T, exp ast.Expression, value int64) {
	t.Helper()
	intLit, ok := exp.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp is not *ast.IntegerLiteral. got=%T", exp)
	}
	if intLit.Token.Type != token.INTEGER {
		t.Errorf("intLit.Token.Type is not %q. got=%q", token.INTEGER, intLit.Token.Type)
	}
	if intLit.Value != value {
		t.Errorf("intLit.Value is not %d. got=%d", value, intLit.Value)
	}
}
