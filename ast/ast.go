package ast

import (
	"fmt"
	"strings"

	"github.com/smalldevshima/go-monkey/token"
)

/// Constants / Variables

const (
	emptyExpressionValue = "<NIL>"
)

/// Types

// Node is the base interface of the AST.
// Every node in the AST has to implement this interface.
type Node interface {
	// TokenLiteral produces the string literal that the node is associated with.
	// It is only used for debugging and testing.
	TokenLiteral() string
	String() string
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

func (p *Program) String() string {
	var out strings.Builder

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type ExpressionStatement struct {
	// the first token of the expression
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	value := emptyExpressionValue
	if es.Expression != nil {
		value = es.Expression.String()
	}
	return fmt.Sprintf("%s;", value)
}

type LetStatement struct {
	// the token.LET token
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	value := emptyExpressionValue
	if ls.Value != nil {
		value = ls.Value.String()
	}
	return fmt.Sprintf("%s %s = %s;", ls.TokenLiteral(), ls.Name, value)
}

type ReturnStatement struct {
	// the token.RETURN token
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	value := emptyExpressionValue
	if rs.ReturnValue != nil {
		value = rs.ReturnValue.String()
	}
	return fmt.Sprintf("%s %s;", rs.TokenLiteral(), value)
}

type BlockStatement struct {
	// the token.LBRACE token
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out strings.Builder

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type Identifier struct {
	// the token.IDENTIFIER token
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

type IntegerLiteral struct {
	// the token.INTEGER token
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.TokenLiteral() }

type BooleanLiteral struct {
	// the boolean token, e.g. token.TRUE or token.FALSE
	Token token.Token
	Value bool
}

func (bl *BooleanLiteral) expressionNode()      {}
func (bl *BooleanLiteral) TokenLiteral() string { return bl.Token.Literal }
func (bl *BooleanLiteral) String() string       { return bl.TokenLiteral() }

type PrefixExpression struct {
	// the prefix token, e.g. token.BANG or token.DASH
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	return fmt.Sprintf("(%s%s)", pe.Operator, pe.Right)
}

type InfixExpression struct {
	// the token of the operator
	Token       token.Token
	Left, Right Expression
	Operator    string
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	return fmt.Sprintf("(%s %s %s)", ie.Left, ie.Operator, ie.Right)
}

type IfExpression struct {
	// the token.IF token
	Token     token.Token
	Condition Expression
	// The if-branch of the if-else expression
	Then *BlockStatement
	// The optional else-branch of the if-else expression, possibly nil
	Otherwise *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	condition := emptyExpressionValue
	then := ""
	otherwise := ""
	if ie.Condition != nil {
		condition = ie.Condition.String()
	}
	if ie.Then != nil {
		then = ie.Then.String()
	}
	if ie.Otherwise != nil {
		otherwise = " " + ie.Otherwise.String()
	}
	return fmt.Sprintf("if %s %s%s", condition, then, otherwise)
}
