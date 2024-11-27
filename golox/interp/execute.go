package interp

import (
	"fmt"
	"toterich/golox/ast"
)

func Execute(stmt ast.Stmt) error {
	var err error

	switch stmt.Type {
	case ast.ST_EXPR:
		_, err = Evaluate(stmt.Expr)
	case ast.ST_PRINT:
		var value ast.LoxValue
		value, err = Evaluate(stmt.Expr)
		fmt.Println(value)
	default:
		panic("Incomplete Switch")
	}

	return err
}
