package lexer

import (
	"testing"

	"github.com/smalldevshima/go-monkeyi/token"
)

/// Constants / Variables

var (
	testLetAssignment = lexerTest{
		name: "let assignment",
		input: `
			let five = 5;
			let ten = 10;
			`,
		expectedTokens: []token.Token{
			{Type: token.LET, Literal: "let"},
			{Type: token.IDENTIFIER, Literal: "five"},
			{Type: token.ASSIGN, Literal: "="},
			{Type: token.INTEGER, Literal: "5"},
			{Type: token.SEMICOLON, Literal: ";"},
			{Type: token.LET, Literal: "let"},
			{Type: token.IDENTIFIER, Literal: "ten"},
			{Type: token.ASSIGN, Literal: "="},
			{Type: token.INTEGER, Literal: "10"},
			{Type: token.SEMICOLON, Literal: ";"},
			{Type: token.EOF, Literal: ""},
		},
	}
	testFunctionDefinition = lexerTest{
		name: "function definition",
		input: `
			let add = fn(x, y) {
				x + y;
			};
			`,
		expectedTokens: []token.Token{
			{Type: token.LET, Literal: "let"},
			{Type: token.IDENTIFIER, Literal: "add"},
			{Type: token.ASSIGN, Literal: "="},
			{Type: token.FUNCTION, Literal: "fn"},
			{Type: token.LPAREN, Literal: "("},
			{Type: token.IDENTIFIER, Literal: "x"},
			{Type: token.COMMA, Literal: ","},
			{Type: token.IDENTIFIER, Literal: "y"},
			{Type: token.RPAREN, Literal: ")"},
			{Type: token.LBRACE, Literal: "{"},
			{Type: token.IDENTIFIER, Literal: "x"},
			{Type: token.PLUS, Literal: "+"},
			{Type: token.IDENTIFIER, Literal: "y"},
			{Type: token.SEMICOLON, Literal: ";"},
			{Type: token.RBRACE, Literal: "}"},
			{Type: token.SEMICOLON, Literal: ";"},
			{Type: token.EOF, Literal: ""},
		},
	}
	testFunctionCall = lexerTest{
		name: "function call",
		input: `
		let result = add(five, ten);
		`,
		expectedTokens: []token.Token{
			{Type: token.LET, Literal: "let"},
			{Type: token.IDENTIFIER, Literal: "result"},
			{Type: token.ASSIGN, Literal: "="},
			{Type: token.IDENTIFIER, Literal: "add"},
			{Type: token.LPAREN, Literal: "("},
			{Type: token.IDENTIFIER, Literal: "five"},
			{Type: token.COMMA, Literal: ","},
			{Type: token.IDENTIFIER, Literal: "ten"},
			{Type: token.RPAREN, Literal: ")"},
			{Type: token.SEMICOLON, Literal: ";"},
			{Type: token.EOF, Literal: ""},
		},
	}
)

/// Tests

func TestNextToken(t *testing.T) {
	lexerTests := []lexerTest{
		testLetAssignment,
		testFunctionDefinition,
		testFunctionCall,
	}
	for index, lexTest := range lexerTests {
		good := t.Run(lexTest.name, func(tt *testing.T) {
			lex := New(lexTest.input)

			for index, expected := range lexTest.expectedTokens {
				have := lex.NextToken()

				if have.Type != expected.Type {
					tt.Fatalf("token-test[%d] - tokentype wrong. expected=%q, have=%q", index, expected.Type, have.Type)
				}

				if have.Literal != expected.Literal {
					tt.Fatalf("token-test[%d] - literal wrong. expected=%q, have=%q", index, expected.Literal, have.Literal)
				}
			}
		})
		if !good {
			t.Logf("lexer-test[%d] %q failed.", index, lexTest.name)
			t.Fail()
		}
	}
}

/// Types

type lexerTest struct {
	name           string
	input          string
	expectedTokens []token.Token
}
