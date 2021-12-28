package lexer

import "github.com/smalldevshima/go-monkeyi/token"

/// Types

type Lexer struct {
	input string
	// current position in input (point to current char)
	position int
	// current reading position in input (after current char)
	readPosition int
	// current char under examination
	char byte
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	newToken := func(tokenType token.TokenType, char byte) token.Token {
		return token.Token{Type: tokenType, Literal: string(char)}
	}

	switch l.char {
	case '=':
		tok = newToken(token.ASSIGN, l.char)
	case ';':
		tok = newToken(token.SEMICOLON, l.char)
	case '(':
		tok = newToken(token.LPAREN, l.char)
	case ')':
		tok = newToken(token.RPAREN, l.char)
	case ',':
		tok = newToken(token.COMMA, l.char)
	case '+':
		tok = newToken(token.PLUS, l.char)
	case '{':
		tok = newToken(token.LBRACE, l.char)
	case '}':
		tok = newToken(token.RBRACE, l.char)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	}
	l.readChar()
	return tok
}

// readChar sets the char field to the next character at read position of the input.
// The position is updated to the read position and the read position is advanced by 1.
// If the read position exceeds the size of the input, then the char field is set to 0.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}
