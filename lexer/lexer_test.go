package lexer

import (
	"testing"

	"github.com/smalldevshima/go-monkeyi/token"
)

/// Constants / Variables

var (
	testLetInitialization = lexerTest{
		name: "let variable initialization",
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
		},
	}
	testOperators = lexerTest{
		name:  "operators",
		input: `+ - * / ! < >`,
		expectedTokens: []token.Token{
			{Type: token.PLUS, Literal: "+"},
			{Type: token.DASH, Literal: "-"},
			{Type: token.ASTERISK, Literal: "*"},
			{Type: token.SLASH, Literal: "/"},
			{Type: token.BANG, Literal: "!"},
			{Type: token.LESS, Literal: "<"},
			{Type: token.GREATER, Literal: ">"},
		},
	}
	testKeywords = lexerTest{
		name:  "keywords",
		input: `fn return true false let if else`,
		expectedTokens: []token.Token{
			{Type: token.FUNCTION, Literal: "fn"},
			{Type: token.RETURN, Literal: "return"},
			{Type: token.TRUE, Literal: "true"},
			{Type: token.FALSE, Literal: "false"},
			{Type: token.LET, Literal: "let"},
			{Type: token.IF, Literal: "if"},
			{Type: token.ELSE, Literal: "else"},
		},
	}
)

/// Tests

func TestNextToken(t *testing.T) {
	lexerTests := []lexerTest{
		testLetInitialization,
		testFunctionDefinition,
		testFunctionCall,
		testOperators,
		testKeywords,
	}
	for index, lexTest := range lexerTests {
		good := t.Run(lexTest.name, func(tt *testing.T) {
			lex := New(lexTest.input)

			for index, expected := range lexTest.expectedTokens {
				if expected.Type == token.EOF {
					// EOF is handled after loop
					break
				}
				have := lex.NextToken()

				if have.Type != expected.Type {
					tt.Fatalf("token-test[%d] - tokentype wrong. expected=%q, have=%q", index, expected.Type, have.Type)
				}

				if have.Literal != expected.Literal {
					tt.Fatalf("token-test[%d] - literal wrong. expected=%q, have=%q", index, expected.Literal, have.Literal)
				}
			}
			end := lex.NextToken()
			if end.Type != token.EOF {
				tt.Fatalf("token-test[%d] - unconsumed token in input. expected=%q, have=%q with literal=%q", index, token.EOF, end.Type, end.Literal)
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
