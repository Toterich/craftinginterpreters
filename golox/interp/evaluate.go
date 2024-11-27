package interp

import (
	"fmt"
	"toterich/golox/ast"
	"toterich/golox/util"
)

func (i Interpreter) Evaluate(expr ast.Expr) (ast.LoxValue, error) {
	switch expr.Type {
	case ast.EXPR_LITERAL:
		return expr.Token.Literal, nil
	case ast.EXPR_IDENTIFIER:
		val, ok := i.vars[expr.Token.Lexeme]
		if !ok {
			return val, util.NewRuntimeError(expr.Token, "undeclared identifier.")
		}
		return val, nil
	case ast.EXPR_UNARY:
		return i.evalUnary(expr)
	case ast.EXPR_BINARY:
		return i.evalBinary(expr)
	case ast.EXPR_GROUPING:
		return i.evalGrouping(expr)
	case ast.EXPR_ASSIGN:
		return i.evalAssignment(expr)
	}

	panic(fmt.Sprintf("Unhandled expression type %d in Evaluate", expr.Type))
}

func (i Interpreter) evalUnary(expr ast.Expr) (ast.LoxValue, error) {
	right, err := i.Evaluate(expr.Children[0])
	if err != nil {
		return right, err
	}

	switch expr.Token.Type {
	case ast.MINUS:
		err := checkType(expr.Token, ast.LT_NUMBER, right.Type)
		if err != nil {
			return right, err
		}
		return ast.NewNumberValue(-right.AsNumber()), nil
	case ast.BANG:
		return ast.NewBoolValue(isTruthy(right)), nil
	}

	panic(fmt.Sprintf("Unhandled operator %d in Unary Expression", expr.Token.Type))
}

func (i Interpreter) evalBinary(expr ast.Expr) (ast.LoxValue, error) {
	left, err := i.Evaluate(expr.Children[0])
	if err != nil {
		return left, err
	}
	right, err := i.Evaluate(expr.Children[1])
	if err != nil {
		return right, err
	}

	switch expr.Token.Type {
	case ast.MINUS:
		err = checkTypes(expr.Token, []ast.LoxType{ast.LT_NUMBER, ast.LT_NUMBER}, []ast.LoxType{left.Type, right.Type})
		if err != nil {
			return ast.NewNilValue(), err
		}
		return ast.NewNumberValue(left.AsNumber() - right.AsNumber()), nil
	case ast.SLASH:
		err = checkTypes(expr.Token, []ast.LoxType{ast.LT_NUMBER, ast.LT_NUMBER}, []ast.LoxType{left.Type, right.Type})
		if err != nil {
			return ast.NewNilValue(), err
		}
		if right.AsNumber() == 0 {
			return ast.NewNilValue(), util.NewRuntimeError(expr.Token, "division by zero")
		}
		return ast.NewNumberValue(left.AsNumber() / right.AsNumber()), nil
	case ast.STAR:
		err = checkTypes(expr.Token, []ast.LoxType{ast.LT_NUMBER, ast.LT_NUMBER}, []ast.LoxType{left.Type, right.Type})
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
				util.NewRuntimeError(expr.Token,
					fmt.Sprintf("Expected either [Number Number] or [String String] as operator's arguments, got [%s %s]", left.Type, right.Type))
		}
	case ast.GREATER:
		err = checkTypes(expr.Token, []ast.LoxType{ast.LT_NUMBER, ast.LT_NUMBER}, []ast.LoxType{left.Type, right.Type})
		if err != nil {
			return ast.NewNilValue(), err
		}
		return ast.NewBoolValue(left.AsNumber() > right.AsNumber()), nil
	case ast.GREATER_EQUAL:
		err = checkTypes(expr.Token, []ast.LoxType{ast.LT_NUMBER, ast.LT_NUMBER}, []ast.LoxType{left.Type, right.Type})
		if err != nil {
			return ast.NewNilValue(), err
		}
		return ast.NewBoolValue(left.AsNumber() >= right.AsNumber()), nil
	case ast.LESS:
		err = checkTypes(expr.Token, []ast.LoxType{ast.LT_NUMBER, ast.LT_NUMBER}, []ast.LoxType{left.Type, right.Type})
		if err != nil {
			return ast.NewNilValue(), err
		}
		return ast.NewBoolValue(left.AsNumber() < right.AsNumber()), nil
	case ast.LESS_EQUAL:
		err = checkTypes(expr.Token, []ast.LoxType{ast.LT_NUMBER, ast.LT_NUMBER}, []ast.LoxType{left.Type, right.Type})
		if err != nil {
			return ast.NewNilValue(), err
		}
		return ast.NewBoolValue(left.AsNumber() <= right.AsNumber()), nil
	case ast.BANG_EQUAL:
		return ast.NewBoolValue(!isEqual(left, right)), nil
	case ast.EQUAL_EQUAL:
		return ast.NewBoolValue(isEqual(left, right)), nil
	}

	panic(fmt.Sprintf("Unhandled operator %d in Binary Expression", expr.Token.Type))
}

func (i Interpreter) evalGrouping(expr ast.Expr) (ast.LoxValue, error) {
	return i.Evaluate(expr.Children[0])
}

func (i Interpreter) evalAssignment(expr ast.Expr) (ast.LoxValue, error) {
	_, ok := i.vars[expr.Token.Lexeme]
	if !ok {
		return ast.NewNilValue(), util.NewRuntimeError(expr.Token, "left hand side of assignment has not been declared")
	}

	val, err := i.Evaluate(expr.Children[0])
	if err != nil {
		return val, err
	}

	i.vars[expr.Token.Lexeme] = val
	return val, nil
}

func isTruthy(value ast.LoxValue) bool {
	switch value.Type {
	case ast.LT_NIL:
		return false
	case ast.LT_BOOL:
		return value.AsBool()
	}

	return true
}

func isEqual(left ast.LoxValue, right ast.LoxValue) bool {
	if left.Type == ast.LT_NIL && right.Type == ast.LT_NIL {
		return true
	}

	return left == right
}

func checkType(token ast.Token, expected ast.LoxType, actual ast.LoxType) error {
	return checkTypes(token, []ast.LoxType{expected}, []ast.LoxType{actual})
}

func checkTypes(token ast.Token, expected []ast.LoxType, actual []ast.LoxType) error {
	util.Assert(len(expected) == len(actual), "expected and actual need to be of equal length")

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
