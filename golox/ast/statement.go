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
	ST_FUNDECL
)

type Stmt struct {
	Type     StmtType
	Expr     Expr
	Tokens   []Token
	Children []Stmt
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

// In a Var Decl, the token is the identifier of the new variable
func NewVarDeclStmt(token Token) Stmt {
	return Stmt{Type: ST_VARDECL, Tokens: []Token{token}, Expr: nil}
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

// In a Function Stmt, the first Token is the function identifier and the subsequent ones are
// the function parameters
func NewFunDeclStmt(name Token, params []Token, body []Stmt) Stmt {
	return Stmt{Type: ST_FUNDECL, Tokens: append([]Token{name}, params...), Children: body}
}
