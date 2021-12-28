package parser

import (
	"github.com/smalldevshima/go-monkey/ast"
	"github.com/smalldevshima/go-monkey/lexer"
	"github.com/smalldevshima/go-monkey/token"
)

/// Types

type Parser struct {
	lx *lexer.Lexer

	currentToken token.Token
	peekToken    token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{lx: l}

	// Read two tokens, so currentToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.lx.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	return nil
}
