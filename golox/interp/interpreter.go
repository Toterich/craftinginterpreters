package interp

import (
	"fmt"
	"toterich/golox/ast"
)

type Interpreter struct {
	env environment
}

func NewInterpreter() Interpreter {
	return Interpreter{env: newEnvironment()}
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
		i.env.setVar(stmt.Ident, value)
	case ast.ST_BLOCK:
		i.env.push()

		for _, child := range stmt.Children {
			i.Execute(child)
		}

		i.env.pop()

	default:
		panic("Incomplete Switch")
	}

	return err
}
