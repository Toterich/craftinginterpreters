package ast

type Ast struct {
	Body        []Stmt
	Statements  StmtStore
	Expressions ExprStore
}
