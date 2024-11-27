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

func (stmt Stmt) PrettyPrint() string {
	switch stmt.Type {
	case ST_INVALID:
		return "INVALID;"
	case ST_EXPR:
		return stmt.Expr.PrettyPrint() + ";"
	case ST_PRINT:
		return "Print: " + stmt.Expr.PrettyPrint() + ";"
	}

	panic("Incomplete switch")
}
