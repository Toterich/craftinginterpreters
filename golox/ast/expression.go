package ast

// Tag interface for Expression types
type Expr interface {
	isExpr()
}

type LiteralExpr struct {
	Token Token
}

func (e LiteralExpr) isExpr() {}

type UnaryExpr struct {
	Operator Token
	Operand  Expr
}

func (e UnaryExpr) isExpr() {}

type BinaryExpr struct {
	Operator Token
	Left     Expr
	Right    Expr
}

func (e BinaryExpr) isExpr() {}

type GroupingExpr struct {
	Grouped Expr
}

func (e GroupingExpr) isExpr() {}

type IdentifierExpr struct {
	Token Token
}

func (e IdentifierExpr) isExpr() {}

type AssignExpr struct {
	Target Token
	Value  Expr
}

func (e AssignExpr) isExpr() {}

type OrExpr struct {
	Left  Expr
	Right Expr
}

func (e OrExpr) isExpr() {}

type AndExpr struct {
	Left  Expr
	Right Expr
}

func (e AndExpr) isExpr() {}

type CallExpr struct {
	Location  Token
	Callee    Expr
	Arguments []Expr
}

func (e CallExpr) isExpr() {}

type ExprStore struct {
	Literal    []LiteralExpr
	Unary      []UnaryExpr
	Binary     []BinaryExpr
	Grouping   []GroupingExpr
	Identifier []IdentifierExpr
	Assign     []AssignExpr
	Or         []OrExpr
	And        []AndExpr
	Call       []CallExpr
}

func (es *ExprStore) NewLiteralExpr(token Token) *LiteralExpr {
	idx := len(es.Literal)
	es.Literal = append(es.Literal, LiteralExpr{Token: token})
	return &es.Literal[idx]
}

func (es *ExprStore) NewUnaryExpr(operator Token, operand Expr) *UnaryExpr {
	idx := len(es.Unary)
	es.Unary = append(es.Unary, UnaryExpr{Operator: operator, Operand: operand})
	return &es.Unary[idx]
}

func (es *ExprStore) NewBinaryExpr(operator Token, left Expr, right Expr) *BinaryExpr {
	idx := len(es.Binary)
	es.Binary = append(es.Binary, BinaryExpr{Operator: operator, Left: left, Right: right})
	return &es.Binary[idx]
}

func (es *ExprStore) NewGroupingExpr(grouped Expr) *GroupingExpr {
	idx := len(es.Grouping)
	es.Grouping = append(es.Grouping, GroupingExpr{Grouped: grouped})
	return &es.Grouping[idx]
}

func (es *ExprStore) NewIdentifierExpr(token Token) *IdentifierExpr {
	idx := len(es.Identifier)
	es.Identifier = append(es.Identifier, IdentifierExpr{Token: token})
	return &es.Identifier[idx]
}

func (es *ExprStore) NewAssignExpr(target Token, value Expr) *AssignExpr {
	idx := len(es.Assign)
	es.Assign = append(es.Assign, AssignExpr{Target: target, Value: value})
	return &es.Assign[idx]
}

func (es *ExprStore) NewOrExpr(left Expr, right Expr) *OrExpr {
	idx := len(es.Or)
	es.Or = append(es.Or, OrExpr{Left: left, Right: right})
	return &es.Or[idx]
}

func (es *ExprStore) NewAndExpr(left Expr, right Expr) *AndExpr {
	idx := len(es.And)
	es.And = append(es.And, AndExpr{Left: left, Right: right})
	return &es.And[idx]
}

func (es *ExprStore) NewCallExpr(location Token, callee Expr, arguments []Expr) *CallExpr {
	idx := len(es.Call)
	es.Call = append(es.Call, CallExpr{Location: location, Callee: callee, Arguments: arguments})
	return &es.Call[idx]
}
