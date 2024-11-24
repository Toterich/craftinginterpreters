package parser

type ExprType int

const (
	EXPR_LITERAL ExprType = iota
	EXPR_UNARY
	EXPR_BINARY
	EXPR_GROUPING
)

type Expr struct {
	tag      ExprType
	token    Token
	children []Expr
}

func NewLiteralExpr(token Token) Expr {
	return Expr{tag: EXPR_LITERAL, token: token}
}

func NewUnaryExpr(child Expr, operand Token) Expr {
	return Expr{tag: EXPR_UNARY, token: operand, children: []Expr{child}}
}

func NewBinaryExpr(left Expr, operand Token, right Expr) Expr {
	return Expr{tag: EXPR_BINARY, token: operand, children: []Expr{left, right}}
}

func NewGroupingExpr(expr Expr) Expr {
	return Expr{tag: EXPR_GROUPING, children: []Expr{expr}}
}
