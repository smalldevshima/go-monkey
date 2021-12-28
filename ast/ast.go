package ast

import "github.com/smalldevshima/go-monkey/token"

/// Types

// Node is the base interface of the AST.
// Every node in the AST has to implement this interface.
type Node interface {
	// TokenLiteral produces the string literal that the node is associated with.
	// It is only used for debugging and testing.
	TokenLiteral() string
}

type Statement interface {
	Node
	// statementNode is a marker function that is only added to help the go tools.
	statementNode()
}

type Expression interface {
	Node
	// expressionNode is a marker function that is only added to help the go tools.
	expressionNode()
}

// Program is the type of the root node of the AST.
// Every Monkey program consists of a series of statements.
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

type LetStatement struct {
	// the token.LET token
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return i.Token.Literal }

type Identifier struct {
	// the token.IDENTIFIER token
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
