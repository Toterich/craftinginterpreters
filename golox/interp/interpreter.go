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
	var err error

	switch stmt := stmt.(type) {

	case *ast.ExprStmt:
		_, err = i.Evaluate(stmt.Expr)

	case *ast.PrintStmt:
		value, err := i.Evaluate(stmt.Expr)
		if err == nil {
			fmt.Println(value)
		}

	case *ast.VarDeclStmt:
		value, err := i.Evaluate(stmt.Value)
		if err == nil {
			i.env.declareVal(stmt.Identifier.Lexeme, value)
		}

	case *ast.BlockStmt:
		i.env.push(false)

		for _, child := range stmt.Body {
			err = i.Execute(child)
			if err != nil || i.doBreak {
				break
			}
		}

		i.env.pop()

	case *ast.IfStmt:
		var doIf ast.LoxValue
		doIf, err = i.Evaluate(stmt.Condition)
		if err != nil {
			break
		}
		if doIf.IsTruthy() {
			err = i.Execute(stmt.Then)
		} else if stmt.Else != nil {
			err = i.Execute(stmt.Else)
		}

	case *ast.WhileStmt:
		var doWhile ast.LoxValue
		doWhile, err = i.Evaluate(stmt.Condition)
		for doWhile.IsTruthy() && !i.doBreak {
			err = i.Execute(stmt.Then)
			if err != nil {
				break
			}
			doWhile, err = i.Evaluate(stmt.Condition)
			if err != nil {
				break
			}
		}
		// Only break out of the innermost loop
		i.doBreak = false

	case *ast.BreakStmt:
		i.doBreak = true

	case *ast.FunDeclStmt:
		fun := ast.LoxFunction{Declaration: stmt}
		i.env.declareVal(stmt.Name.Lexeme, ast.NewFunction(fun))

	default:
		panic(assert.MissingCase(stmt))
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
	for idx, param := range callee.Declaration.Params {
		i.env.declareVal(param.Lexeme, arguments[idx])
	}

	// Execute statements one after another in the local env
	for _, statement := range callee.Declaration.Body {
		err := i.Execute(statement)
		if err != nil {
			return ast.NewNilValue(), err
		}
	}

	return ast.NewNilValue(), nil
}
