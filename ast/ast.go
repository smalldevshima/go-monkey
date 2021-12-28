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
	statementNode()
}

type Expression interface {
	Node
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
	Token token.Token // the token.LET token

}
