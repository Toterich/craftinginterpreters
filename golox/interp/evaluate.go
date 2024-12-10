package interp

import (
	"fmt"
	"toterich/golox/ast"
	"toterich/golox/util"
	"toterich/golox/util/assert"
)

func (i Interpreter) Evaluate(expr ast.Expr) (ast.LoxValue, error) {
	switch expr := expr.(type) {
	case *ast.LiteralExpr:
		return expr.Token.Literal, nil
	case *ast.IdentifierExpr:
		val, ok := i.env.getVar(expr.Token.Lexeme)
		if !ok {
			return val, util.NewRuntimeError(expr.Token, "undeclared identifier.")
		}
		return val, nil
	case *ast.UnaryExpr:
		return i.evalUnary(expr)
	case *ast.BinaryExpr:
		return i.evalBinary(expr)
	case *ast.GroupingExpr:
		return i.evalGrouping(expr)
	case *ast.AssignExpr:
		return i.evalAssignment(expr)
	case *ast.OrExpr:
		return i.evalOr(expr)
	case *ast.AndExpr:
		return i.evalAnd(expr)
	case *ast.CallExpr:
		return i.evalCall(expr)
	default:
		panic(assert.MissingCase(expr))
	}
}

func (i Interpreter) evalUnary(expr *ast.UnaryExpr) (ast.LoxValue, error) {
	right, err := i.Evaluate(expr.Operand)
	if err != nil {
		return right, err
	}

	switch expr.Operator.Type {
	case ast.MINUS:
		err := checkType(expr.Operator, ast.LT_NUMBER, right.Type)
		if err != nil {
			return right, err
		}
		return ast.NewNumberValue(-right.AsNumber()), nil
	case ast.BANG:
		return ast.NewBoolValue(right.IsTruthy()), nil
	}

	panic(assert.MissingCase(expr.Operator.Type))
}

func (i Interpreter) evalBinary(expr *ast.BinaryExpr) (ast.LoxValue, error) {
	left, err := i.Evaluate(expr.Left)
	if err != nil {
		return left, err
	}
	right, err := i.Evaluate(expr.Right)
	if err != nil {
		return right, err
	}

	switch expr.Operator.Type {
	case ast.MINUS:
		err = checkTypes(expr.Operator, []ast.LoxType{ast.LT_NUMBER, ast.LT_NUMBER}, []ast.LoxType{left.Type, right.Type})
		if err != nil {
			return ast.NewNilValue(), err
		}
		return ast.NewNumberValue(left.AsNumber() - right.AsNumber()), nil
	case ast.SLASH:
		err = checkTypes(expr.Operator, []ast.LoxType{ast.LT_NUMBER, ast.LT_NUMBER}, []ast.LoxType{left.Type, right.Type})
		if err != nil {
			return ast.NewNilValue(), err
		}
		if right.AsNumber() == 0 {
			return ast.NewNilValue(), util.NewRuntimeError(expr.Operator, "division by zero")
		}
		return ast.NewNumberValue(left.AsNumber() / right.AsNumber()), nil
	case ast.STAR:
		err = checkTypes(expr.Operator, []ast.LoxType{ast.LT_NUMBER, ast.LT_NUMBER}, []ast.LoxType{left.Type, right.Type})
		if err != nil {
			return ast.NewNilValue(), err
		}
		return ast.NewNumberValue(left.AsNumber() * right.AsNumber()), nil
	case ast.PLUS:
		if left.Type == ast.LT_NUMBER && right.Type == ast.LT_NUMBER {
			return ast.NewNumberValue(left.AsNumber() + right.AsNumber()), nil
		} else if left.Type == ast.LT_STRING && right.Type == ast.LT_STRING {
			return ast.NewStringValue(left.AsString() + right.AsString()), nil
		} else {
			return ast.NewNilValue(),
				util.NewRuntimeError(expr.Operator,
					fmt.Sprintf("Expected either [Number Number] or [String String] as operator's arguments, got [%s %s]", left.Type, right.Type))
		}
	case ast.GREATER:
		err = checkTypes(expr.Operator, []ast.LoxType{ast.LT_NUMBER, ast.LT_NUMBER}, []ast.LoxType{left.Type, right.Type})
		if err != nil {
			return ast.NewNilValue(), err
		}
		return ast.NewBoolValue(left.AsNumber() > right.AsNumber()), nil
	case ast.GREATER_EQUAL:
		err = checkTypes(expr.Operator, []ast.LoxType{ast.LT_NUMBER, ast.LT_NUMBER}, []ast.LoxType{left.Type, right.Type})
		if err != nil {
			return ast.NewNilValue(), err
		}
		return ast.NewBoolValue(left.AsNumber() >= right.AsNumber()), nil
	case ast.LESS:
		err = checkTypes(expr.Operator, []ast.LoxType{ast.LT_NUMBER, ast.LT_NUMBER}, []ast.LoxType{left.Type, right.Type})
		if err != nil {
			return ast.NewNilValue(), err
		}
		return ast.NewBoolValue(left.AsNumber() < right.AsNumber()), nil
	case ast.LESS_EQUAL:
		err = checkTypes(expr.Operator, []ast.LoxType{ast.LT_NUMBER, ast.LT_NUMBER}, []ast.LoxType{left.Type, right.Type})
		if err != nil {
			return ast.NewNilValue(), err
		}
		return ast.NewBoolValue(left.AsNumber() <= right.AsNumber()), nil
	case ast.BANG_EQUAL:
		return ast.NewBoolValue(!left.IsEqual(right)), nil
	case ast.EQUAL_EQUAL:
		return ast.NewBoolValue(left.IsEqual(right)), nil
	}

	panic(assert.MissingCase(expr.Operator.Type))
}

func (i Interpreter) evalGrouping(expr *ast.GroupingExpr) (ast.LoxValue, error) {
	return i.Evaluate(expr.Grouped)
}

func (i Interpreter) evalAssignment(expr *ast.AssignExpr) (ast.LoxValue, error) {
	_, ok := i.env.getVar(expr.Target.Lexeme)
	if !ok {
		return ast.NewNilValue(), util.NewRuntimeError(expr.Target, "left hand side of assignment has not been declared")
	}

	val, err := i.Evaluate(expr.Value)
	if err != nil {
		return val, err
	}

	// This is already checked by getVar above
	assert.Assert(i.env.assignVal(expr.Target.Lexeme, val), "identifier to be assigned to has not been declared")

	return val, nil
}

func (i Interpreter) evalOr(expr *ast.OrExpr) (ast.LoxValue, error) {
	leftVal, err := i.Evaluate(expr.Left)
	if err != nil {
		return leftVal, err
	}

	// Short circuit
	if leftVal.IsTruthy() {
		return ast.NewBoolValue(true), nil
	}

	rightVal, err := i.Evaluate(expr.Right)
	return ast.NewBoolValue(rightVal.IsTruthy()), err
}

func (i Interpreter) evalAnd(expr *ast.AndExpr) (ast.LoxValue, error) {
	leftVal, err := i.Evaluate(expr.Left)
	if err != nil {
		return leftVal, err
	}

	// Short circuit
	if !leftVal.IsTruthy() {
		return ast.NewBoolValue(false), nil
	}

	rightVal, err := i.Evaluate(expr.Right)
	return ast.NewBoolValue(rightVal.IsTruthy()), err
}

func (i Interpreter) evalCall(expr *ast.CallExpr) (ast.LoxValue, error) {
	callee, err := i.Evaluate(expr.Callee)
	if err != nil {
		return callee, err
	}

	if callee.Type != ast.LT_FUNCTION {
		return callee, util.NewRuntimeError(expr.Location, "callee is not callable.")
	}
	fun := callee.AsFunction()

	var args []ast.LoxValue
	for _, arg := range expr.Arguments {
		arg, err := i.Evaluate(arg)
		if err != nil {
			return arg, err
		}
		args = append(args, arg)
	}

	if fun.Arity() != len(args) {
		return callee, util.NewRuntimeError(expr.Location,
			fmt.Sprintf("callee expects %d arguments, got %d", fun.Arity(), len(args)))
	}

	return i.call(fun, args)
}

func checkType(token ast.Token, expected ast.LoxType, actual ast.LoxType) error {
	return checkTypes(token, []ast.LoxType{expected}, []ast.LoxType{actual})
}

func checkTypes(token ast.Token, expected []ast.LoxType, actual []ast.LoxType) error {
	assert.Assert(len(expected) == len(actual), "expected and actual need to be of equal length")

	equal := true
	for i := 0; i < len(expected); i += 1 {
		if expected[i] != actual[i] {
			equal = false
			break
		}
	}

	if equal {
		return nil
	}

	var msg string
	if len(expected) == 1 {
		msg = fmt.Sprintf("Expected %s as argument, got %s", expected[0], actual[0])
	} else {
		msg = fmt.Sprintf("Expected %s as arguments, got %s", expected, actual)
	}

	return util.NewRuntimeError(token, msg)
}
