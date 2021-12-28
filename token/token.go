package token

/// Constants and Variables

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// Identifiers and literals
	IDENT = "IDENT" // add, foo, bar, "hello", true
	INT   = "INT"   // 123

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
// If the given identifier is not a keyword, it returns the TokenType for identifiers.
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

/// Types

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}
