package token

/// Constants and Variables

// Possible token types for lexer/ parser/ ast
const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	IDENTIFIER TokenType = "IDENTIFIER"
	INTEGER    TokenType = "INTEGER"
	STRING     TokenType = "STRING"

	ASSIGN   TokenType = "="
	PLUS     TokenType = "+"
	DASH     TokenType = "-"
	ASTERISK TokenType = "*"
	SLASH    TokenType = "/"
	BANG     TokenType = "!"
	LT       TokenType = "<"
	GT       TokenType = ">"
	EQ       TokenType = "=="
	NEQ      TokenType = "!="

	COMMA     TokenType = ","
	SEMICOLON TokenType = ";"

	LPAREN   TokenType = "("
	RPAREN   TokenType = ")"
	LBRACE   TokenType = "{"
	RBRACE   TokenType = "}"
	LBRACKET TokenType = "["
	RBRACKET TokenType = "]"

	FUNCTION TokenType = "FUNCTION"
	RETURN   TokenType = "RETURN"
	LET      TokenType = "LET"

	TRUE  TokenType = "TRUE"
	FALSE TokenType = "FALSE"

	IF   TokenType = "IF"
	ELSE TokenType = "ELSE"
)

var (
	// keywords is a map of literals to their corresponding TokenType.
	keywords = map[string]TokenType{
		"fn":     FUNCTION,
		"return": RETURN,
		"let":    LET,
		"true":   TRUE,
		"false":  FALSE,
		"if":     IF,
		"else":   ELSE,
	}
)

/// Functions

// LookupIdent checks if the given identifier is a keyword and if so, returns its TokenType.
// If the given identifier is not a keyword, it returns the TokenType for user defined identifiers.
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENTIFIER
}

/// Types

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}
