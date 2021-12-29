package parser

import (
	"fmt"

	"github.com/smalldevshima/go-monkey/ast"
	"github.com/smalldevshima/go-monkey/lexer"
	"github.com/smalldevshima/go-monkey/token"
)

/// Constant / Variables

/// Types

// The Parser consumes the output of a given lexer.Lexer and produces an ast.Program as its output.
// A Parser's zero value is not usable and new ones need to be created using parser.New.
type Parser struct {
	lx *lexer.Lexer

	currentToken token.Token
	peekToken    token.Token

	errors []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		lx:     l,
		errors: []string{},
	}

	// Read two tokens, so currentToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

// ParseProgram consumes the internal Lexer's token list and produces a Program from them.
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for !p.currentTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

// parseStatement checks the current token type and calls the corresponding parse method.
func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.LET:
		// * check if s is nil, else the wrapped interface type will mask the nil value
		if s := p.parseLetStatement(); s != nil {
			return s
		}
	}
	return nil
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.currentToken}

	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// todo: currently expressions are skipped until a semicolon is found
	for !p.currentTokenIs(token.SEMICOLON) {
		if p.currentTokenIs(token.EOF) {
			return nil
		}
		p.nextToken()
	}

	return stmt
}

// nextToken advances the tokens read from the internal Lexer.
func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lx.NextToken()
}

// expectPeek compares the next token against the provided.
// If they are the same, it advances the tokens and returns true.
// Otherwise it leaves the tokens as is and returns false.
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

// peekError creates a new unexpected-token error message and appends it to the error list.
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("unexpected token of type %q with literal %q, expected token of type %q", p.peekToken.Type, p.peekToken.Literal, t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) currentTokenIs(t token.TokenType) bool {
	return p.currentToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}
