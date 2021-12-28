package parser

import (
	"fmt"
	"testing"

	"github.com/smalldevshima/go-monkey/ast"
	"github.com/smalldevshima/go-monkey/lexer"
)

func TestLetStatements(t *testing.T) {
	input := `
		let x = 5;
		let y = 10;
		let foobar - 838383;
		`

	p := New(lexer.New(input))

	program := p.ParseProgram()
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

func checkLetStatement(t *testing.T, s ast.Statement, name string) {
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