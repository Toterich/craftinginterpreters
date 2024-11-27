package ast

type StmtType int

const (
	ST_INVALID StmtType = iota
	ST_EXPR
	ST_PRINT
)

type Stmt struct {
	Type StmtType
	Expr Expr
}

func NewInvalidStmt() Stmt {
	return Stmt{Type: ST_INVALID}
}

func NewExprStmt(expr Expr) Stmt {
	return Stmt{Type: ST_EXPR, Expr: expr}
}

func NewPrintStmt(expr Expr) Stmt {
	return Stmt{Type: ST_PRINT, Expr: expr}
}
