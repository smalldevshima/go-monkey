package parser

import (
	"fmt"
	"strconv"

	"github.com/smalldevshima/go-monkey/ast"
	"github.com/smalldevshima/go-monkey/lexer"
	"github.com/smalldevshima/go-monkey/token"
)

/// Constant / Variables

// Expression evaluation precedence constants where a higher value means more precedence.
const (
	_ uint = iota
	LOWEST
	EQUALS
	LTGT
	SUM
	PRODUCT
	PREFIX
	CALL
)

/// Types

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// The Parser consumes the output of a given lexer.Lexer and produces an ast.Program as its output.
// A Parser's zero value is not usable and new ones need to be created using parser.New.
type Parser struct {
	lx *lexer.Lexer

	currentToken token.Token
	peekToken    token.Token

	errors []string

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		lx:             l,
		errors:         []string{},
		prefixParseFns: make(map[token.TokenType]prefixParseFn),
		infixParseFns:  make(map[token.TokenType]infixParseFn),
	}

	// register expression parsing fns
	p.registerPrefix(token.IDENTIFIER, p.parseIdentifier)
	p.registerPrefix(token.INTEGER, p.parseIntegerLiteral)

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
	// * always check if parsed statement is nil, else the wrapped interface type will mask the nil value
	switch p.currentToken.Type {
	case token.LET:
		if s := p.parseLetStatement(); s != nil {
			return s
		}
	case token.RETURN:
		if s := p.parseReturnStatement(); s != nil {
			return s
		}
	default:
		if s := p.parseExpressionStatement(); s != nil {
			return s
		}
	}
	msg := fmt.Sprintf("unexpected token %q with value %q when trying to parse statement", p.currentToken.Type, p.currentToken.Literal)
	p.errors = append(p.errors, msg)
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
			p.peekError(token.SEMICOLON)
			return nil
		}
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currentToken}

	p.nextToken()

	// todo: currently expressions are skipped until a semicolon is found
	for !p.currentTokenIs(token.SEMICOLON) {
		if p.currentTokenIs(token.EOF) {
			p.peekError(token.SEMICOLON)
			return nil
		}
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currentToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence uint) ast.Expression {
	prefix, ok := p.prefixParseFns[p.currentToken.Type]
	if !ok {
		p.noPrefixParseFnError(p.currentToken.Type)
		return nil
	}

	leftExp := prefix()
	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currentToken, Value: p.currentToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as int64", p.currentToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit := &ast.IntegerLiteral{Token: p.currentToken, Value: value}
	return lit
}

// nextToken advances the tokens read from the internal Lexer.
func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lx.NextToken()
}

// expectPeek compares the next token against the provided.
// If they are the same, it advances the tokens and returns true.
// Otherwise it leaves the tokens as is, adds an error to the internal list and returns false.
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

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}
