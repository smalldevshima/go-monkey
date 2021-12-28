package lexer

import "github.com/smalldevshima/go-monkeyi/token"

/// Functions

// isDigit returns true for all ASCII decimal number characters.
func isDigit(char byte) bool {
	return '0' <= char && char <= '9'
}

// isLetter returns true for all ASCII characters that are valid to be used for keywords and identifiers.
func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

// isWhitespace returns true for all ASCII characters that are considered whitespace for the lexer.
func isWhitespace(char byte) bool {
	return char == ' ' || char == '\t' || char == '\n' || char == '\r'
}

// newToken returns a Token with the given type and the given Literal ASCII charcode.
func newToken(tokenType token.TokenType, char byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(char)}
}

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

	l.skipWhitespace()

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
	default:
		if isLetter(l.char) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			// return immediately to not advance read position further
			return tok
		} else if isDigit(l.char) {
			tok.Literal = l.readInteger()
			tok.Type = token.INTEGER
			// return immediately to not advance read position further
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.char)
		}
	}
	l.readChar()
	return tok
}

// readChar sets the char field to the next character at read position of the input.
// The position is updated to the read position and the read position is advanced by 1.
//
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

// readIdentifier consumes and returns a whole word up to the next character where isLetter=false.
func (l *Lexer) readIdentifier() string {
	start := l.position
	for isLetter(l.char) {
		l.readChar()
	}
	return l.input[start:l.position]
}

// readInteger consumes and returns a whole number up to the next character where isDigit=false.
func (l *Lexer) readInteger() string {
	start := l.position
	for isDigit(l.char) {
		l.readChar()
	}
	return l.input[start:l.position]
}

// skipWhitespace consumes the input until the next character where isWhitespace=false.
func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.char) {
		l.readChar()
	}
}
