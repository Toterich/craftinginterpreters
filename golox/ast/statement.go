package ast

type StmtType int

const (
	ST_INVALID StmtType = iota
	ST_EXPR
	ST_PRINT
	ST_VARDECL
	ST_BLOCK
	ST_IF
	ST_WHILE
	ST_FOR
	ST_BREAK
)

type Stmt struct {
	Type     StmtType
	Expr     Expr
	Ident    string // For VarDecl Statement
	Children []Stmt // For Block Statement
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

func NewVarDeclStmt(ident string) Stmt {
	// TODO: This is a bit weird because Expressions always store a Token indicating the position where
	// the Expression occurs in the source, but for empty variable declarations we need to assign
	// an implicit nil Expression. Maybe Expressions need to be decoupled form source locations in general?
	return Stmt{Type: ST_VARDECL, Ident: ident, Expr: NewLiteralExpr(Token{Type: NIL, Line: -1})}
}

func NewBlockStmt(children []Stmt) Stmt {
	return Stmt{Type: ST_BLOCK, Children: children}
}

func NewIfStmt(condition Expr, ifBranch Stmt, elseBranch Stmt) Stmt {
	return Stmt{Type: ST_IF, Expr: condition, Children: []Stmt{ifBranch, elseBranch}}
}

func NewWhileStmt(condition Expr, loop Stmt) Stmt {
	return Stmt{Type: ST_WHILE, Expr: condition, Children: []Stmt{loop}}
}

func NewBreakStmt() Stmt {
	return Stmt{Type: ST_BREAK}
}
