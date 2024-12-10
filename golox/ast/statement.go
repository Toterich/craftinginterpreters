package ast

type Stmt interface {
	isStmt()
}

type ExprStmt struct {
	Expr Expr
}

func (s ExprStmt) isStmt() {}

type PrintStmt struct {
	Expr Expr
}

func (s PrintStmt) isStmt() {}

type VarDeclStmt struct {
	Identifier Token
	Value      Expr
}

func (s VarDeclStmt) isStmt() {}

type BlockStmt struct {
	Body []Stmt
}

func (s BlockStmt) isStmt() {}

type IfStmt struct {
	Condition Expr
	Then      Stmt
	Else      Stmt
}

func (s IfStmt) isStmt() {}

type WhileStmt struct {
	Condition Expr
	Then      Stmt
}

func (s WhileStmt) isStmt() {}

type BreakStmt struct {
}

func (s BreakStmt) isStmt() {}

type FunDeclStmt struct {
	Name   Token
	Params []Token
	Body   []Stmt
}

func (s FunDeclStmt) isStmt() {}

type StmtStore struct {
	Expr    []ExprStmt
	Print   []PrintStmt
	VarDecl []VarDeclStmt
	Block   []BlockStmt
	If      []IfStmt
	While   []WhileStmt
	Break   []BreakStmt
	FunDecl []FunDeclStmt
}

func (ss *StmtStore) NewExpr(expr Expr) *ExprStmt {
	idx := len(ss.Expr)
	ss.Expr = append(ss.Expr, ExprStmt{Expr: expr})
	return &ss.Expr[idx]
}

func (ss *StmtStore) NewPrint(expr Expr) *PrintStmt {
	idx := len(ss.Print)
	ss.Print = append(ss.Print, PrintStmt{Expr: expr})
	return &ss.Print[idx]
}

func (ss *StmtStore) NewVarDecl(identifier Token, value Expr) *VarDeclStmt {
	idx := len(ss.VarDecl)
	ss.VarDecl = append(ss.VarDecl, VarDeclStmt{Identifier: identifier, Value: value})
	return &ss.VarDecl[idx]
}

func (ss *StmtStore) NewBlock(children []Stmt) *BlockStmt {
	idx := len(ss.Block)
	ss.Block = append(ss.Block, BlockStmt{Body: children})
	return &ss.Block[idx]
}

func (ss *StmtStore) NewIf(condition Expr, then Stmt, else_ Stmt) *IfStmt {
	idx := len(ss.If)
	ss.If = append(ss.If, IfStmt{Condition: condition, Then: then, Else: else_})
	return &ss.If[idx]
}

func (ss *StmtStore) NewWhile(condition Expr, then Stmt) *WhileStmt {
	idx := len(ss.While)
	ss.While = append(ss.While, WhileStmt{Condition: condition, Then: then})
	return &ss.While[idx]
}

func (ss *StmtStore) NewBreak() *BreakStmt {
	idx := len(ss.Break)
	ss.Break = append(ss.Break, BreakStmt{})
	return &ss.Break[idx]
}

func (ss *StmtStore) NewFunDecl(name Token, params []Token, children []Stmt) *FunDeclStmt {
	idx := len(ss.FunDecl)
	ss.FunDecl = append(ss.FunDecl, FunDeclStmt{Name: name, Params: params, Body: children})
	return &ss.FunDecl[idx]
}
