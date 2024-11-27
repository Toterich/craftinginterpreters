package interp

import (
	"fmt"
	"toterich/golox/ast"
)

type Interpreter struct {
	vars map[string]ast.LoxValue
}

func NewInterpreter() Interpreter {
	return Interpreter{vars: map[string]ast.LoxValue{}}
}

func (i *Interpreter) Execute(stmt ast.Stmt) error {

	value, err := i.Evaluate(stmt.Expr)

	// Execute side effects
	switch stmt.Type {
	case ast.ST_EXPR:
		// Evaluation happened already above
	case ast.ST_PRINT:
		if err == nil {
			fmt.Println(value)
		}
	case ast.ST_VARDECL:
		i.vars[stmt.Ident] = value

	default:
		panic("Incomplete Switch")
	}

	return err
}
