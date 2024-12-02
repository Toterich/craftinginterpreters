package ast

type ExprType int

const (
	EXPR_INVALID ExprType = iota
	EXPR_LITERAL
	EXPR_UNARY
	EXPR_BINARY
	EXPR_GROUPING
	EXPR_IDENTIFIER
	EXPR_ASSIGN
	EXPR_OR
	EXPR_AND
)

type Expr struct {
	Type     ExprType
	Token    Token
	Children []Expr
}

func NewLiteralExpr(token Token) Expr {
	return Expr{Type: EXPR_LITERAL, Token: token}
}

func NewUnaryExpr(child Expr, operand Token) Expr {
	return Expr{Type: EXPR_UNARY, Token: operand, Children: []Expr{child}}
}

func NewBinaryExpr(left Expr, operand Token, right Expr) Expr {
	return Expr{Type: EXPR_BINARY, Token: operand, Children: []Expr{left, right}}
}

// Created by wrapping expression in ()
func NewGroupingExpr(expr Expr) Expr {
	return Expr{Type: EXPR_GROUPING, Children: []Expr{expr}}
}

func NewIdentifierExpression(token Token) Expr {
	return Expr{Type: EXPR_IDENTIFIER, Token: token}
}

// Lhs is a Token because it needs to point to a storage location
func NewAssignExpr(left Token, right Expr) Expr {
	return Expr{Type: EXPR_ASSIGN, Token: left, Children: []Expr{right}}
}

func NewOrExpr(left Expr, right Expr) Expr {
	return Expr{Type: EXPR_OR, Children: []Expr{left, right}}
}

func NewAndExpr(left Expr, right Expr) Expr {
	return Expr{Type: EXPR_AND, Children: []Expr{left, right}}
}
