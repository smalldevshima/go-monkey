package token

/// Constants and Variables

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// Identifiers and literals
	IDENTIFIER = "IDENTIFIER" // add, foo, bar, "hello", true
	INTEGER    = "INTEGER"    // 123

	// Operators
	ASSIGN = "="
	PLUS   = "+"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
)

var (
	// keywords is a map of literals to their corresponding TokenType.
	keywords = map[string]TokenType{
		"fn":  FUNCTION,
		"let": LET,
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
