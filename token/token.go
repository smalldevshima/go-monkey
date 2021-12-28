package token

/// Constants and Variables

const (
	ILLEGAL TokenType = "ILLEGAL"
	EOF     TokenType = "EOF"

	// Identifiers and literals
	IDENTIFIER = "IDENTIFIER" // add, foo, bar, "hello", true
	INTEGER    = "INTEGER"    // 123

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	DASH     = "-"
	ASTERISK = "*"
	SLASH    = "/"
	BANG     = "!"
	LT       = "<"
	GT       = ">"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	RETURN   = "RETURN"
	LET      = "LET"

	TRUE  = "TRUE"
	FALSE = "FALSE"

	IF   = "IF"
	ELSE = "ELSE"
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
