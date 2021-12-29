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
	_ Precedence = iota
	LOWEST
	EQUALS
	LTGT
	SUM
	PRODUCT
	PREFIX
	CALL
)

var (
	// prefixTokens is the list of all tokens that are parsed in prefix position
	prefixTokens = []token.TokenType{token.IDENTIFIER, token.INTEGER, token.BANG, token.DASH, token.TRUE, token.FALSE, token.LPAREN, token.COMMA, token.IF}
	// infixTokens is the list of all tokens that are parsed in infix position
	infixTokens = []token.TokenType{token.EQ, token.NEQ, token.LT, token.GT, token.PLUS, token.DASH, token.SLASH, token.ASTERISK}

	// precedences maps every infix operator to its corresponding precedence value
	precedences = map[token.TokenType]Precedence{
		token.EQ:       EQUALS,
		token.NEQ:      EQUALS,
		token.LT:       LTGT,
		token.GT:       LTGT,
		token.PLUS:     SUM,
		token.DASH:     SUM,
		token.SLASH:    PRODUCT,
		token.ASTERISK: PRODUCT,
	}
)

/// Types

type Precedence uint

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
	for _, tok := range prefixTokens {
		p.registerPrefix(tok, p.parsePrefixExpression)
	}
	for _, tok := range infixTokens {
		p.registerInfix(tok, p.parseInfixExpression)
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

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if !p.expectPeek(token.SEMICOLON) {
		return nil
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.currentToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	if !p.expectPeek(token.SEMICOLON) {
		return nil
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

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{
		Token:      p.currentToken,
		Statements: []ast.Statement{},
	}

	p.nextToken()

	for !p.currentTokenIs(token.RBRACE) && !p.currentTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

func (p *Parser) parseExpression(precedence Precedence) ast.Expression {
	prefix, ok := p.prefixParseFns[p.currentToken.Type]
	if !ok {
		p.noPrefixParseFnError(p.currentToken.Type)
		return nil
	}

	leftExp := prefix()

	// * continue extending the expression with infix operators, until you find:
	// * - a semicolon ";",
	// * - the end of the file "EOF",
	// * - a token that has a lower-or-equal precedence (does not bind stronger than the current token), or
	// * - a token that is not an infix token
	for !p.peekTokenIs(token.SEMICOLON) && !p.peekTokenIs(token.EOF) && precedence < p.peekPrecedence() {
		infix, ok := p.infixParseFns[p.peekToken.Type]
		if !ok {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

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

func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{Token: p.currentToken, Value: p.currentTokenIs(token.TRUE)}
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	switch p.currentToken.Type {
	case token.INTEGER:
		if s := p.parseIntegerLiteral(); s != nil {
			return s
		}
	case token.IDENTIFIER:
		if s := p.parseIdentifier(); s != nil {
			return s
		}
	case token.TRUE, token.FALSE:
		if s := p.parseBooleanLiteral(); s != nil {
			return s
		}
	case token.LPAREN:
		p.nextToken()

		exp := p.parseExpression(LOWEST)
		if exp == nil || !p.expectPeek(token.RPAREN) {
			return nil
		}
		return exp
	case token.IF:
		if exp := p.parseIfExpression(); exp != nil {
			return exp
		}
	case token.COMMA:
		// todo
	case token.BANG, token.DASH:
		exp := &ast.PrefixExpression{
			Token:    p.currentToken,
			Operator: p.currentToken.Literal,
		}

		p.nextToken()

		exp.Right = p.parseExpression(PREFIX)
		return exp
	}
	msg := fmt.Sprintf("unhandled token %q with value %q when trying to parse prefix expression", p.currentToken.Type, p.currentToken.Literal)
	p.errors = append(p.errors, msg)
	return nil
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	exp := &ast.InfixExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
		Left:     left,
	}

	pre := p.currentPrecedence()
	p.nextToken()
	exp.Right = p.parseExpression(pre)
	return exp
}

func (p *Parser) parseIfExpression() ast.Expression {
	exp := &ast.IfExpression{Token: p.currentToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	exp.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	exp.Then = p.parseBlockStatement()

	if !p.peekTokenIs(token.ELSE) {
		return exp
	}

	p.nextToken()
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	exp.Otherwise = p.parseBlockStatement()

	return exp
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
	msg := fmt.Sprintf("token %q cannot appear in prefix position", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) currentPrecedence() Precedence {
	if p, ok := precedences[p.currentToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) peekPrecedence() Precedence {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}
