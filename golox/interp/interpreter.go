package interp

import (
	"fmt"
	"toterich/golox/ast"
	"toterich/golox/util/assert"
)

type Interpreter struct {
	env     environment
	doBreak bool
}

func NewInterpreter() Interpreter {
	return Interpreter{env: newEnvironment()}
}

func (i *Interpreter) Execute(stmt ast.Stmt) error {
	var exprValue ast.LoxValue
	var err error

	// Execute side effects
	switch stmt.Type {

	case ast.ST_EXPR:
		_, err = i.Evaluate(stmt.Expr)

	case ast.ST_PRINT:
		exprValue, err = i.Evaluate(stmt.Expr)
		if err == nil {
			fmt.Println(exprValue)
		}

	case ast.ST_VARDECL:
		exprValue, err = i.Evaluate(stmt.Expr)
		if err == nil {
			i.env.declareVal(stmt.Tokens[0].Lexeme, exprValue)
		}

	case ast.ST_BLOCK:
		i.env.push()

		for _, child := range stmt.Children {
			err = i.Execute(child)
			if err != nil || i.doBreak {
				break
			}
		}

		i.env.pop()

	case ast.ST_IF:
		exprValue, err = i.Evaluate(stmt.Expr)
		if exprValue.IsTruthy() {
			err = i.Execute(stmt.Children[0])
		} else if stmt.Children[1].Type != ast.ST_INVALID {
			err = i.Execute(stmt.Children[1])
		}

	case ast.ST_WHILE:
		exprValue, err = i.Evaluate(stmt.Expr)
		for exprValue.IsTruthy() && !i.doBreak {
			err = i.Execute(stmt.Children[0])
			if err != nil {
				break
			}
			exprValue, err = i.Evaluate(stmt.Expr)
			if err != nil {
				break
			}
		}
		// Only break out of the innermost loop
		i.doBreak = false

	case ast.ST_BREAK:
		i.doBreak = true

	default:
		panic(assert.MissingCase(stmt.Type))
	}

	return err
}
