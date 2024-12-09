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
		i.env.push(false)

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

	case ast.ST_FUNDECL:
		fun := ast.LoxFunction{Declaration: stmt}
		i.env.declareVal(stmt.Tokens[0].Lexeme, ast.NewFunction(fun))

	default:
		panic(assert.MissingCase(stmt.Type))
	}

	return err
}

func (i *Interpreter) call(callee ast.LoxFunction, arguments []ast.LoxValue) (ast.LoxValue, error) {
	// For the duration of the call, create a new environment that only inherits from the global env
	// TODO: Functions don't necessarily have access to only global scope. For those declared inside
	// another scope, the call to them should inherit that scope instead
	i.env.push(true)
	defer func() { i.env.pop() }()

	// Declare passed function parameters in local env
	for idx, param := range callee.Declaration.Tokens[1:] {
		i.env.declareVal(param.Lexeme, arguments[idx])
	}

	// Execute statements one after another in the local env
	for _, statement := range callee.Declaration.Children {
		err := i.Execute(statement)
		if err != nil {
			return ast.NewNilValue(), err
		}
	}

	return ast.NewNilValue(), nil
}
