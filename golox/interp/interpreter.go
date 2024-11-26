package interp

import (
	"fmt"
	"toterich/golox/ast"
)

func Evaluate(expr ast.Expr) (ast.LoxValue, error) {
	switch expr.Type {
	case ast.EXPR_UNARY:
		return evalUnary(expr)
	case ast.EXPR_BINARY:
		return evalBinary(expr)
	case ast.EXPR_GROUPING:
		return evalGrouping(expr)
	}

	panic(fmt.Sprintf("Unhandled expression type %d in Evaluate", expr.Type))
}

func evalUnary(expr ast.Expr) (ast.LoxValue, error) {
	right, err := Evaluate(expr.Children[0])
	if err != nil {
		return ast.NewNilValue(), err
	}

	switch expr.Token.Type {
	case ast.MINUS:
		if right.Type != ast.LT_NUMBER {
			return ast.NewNilValue(), fmt.Errorf("expected number after unary operator -, got %s", right.Type)
		}
		return ast.NewNumberValue(-right.Number()), nil
	case ast.BANG:
		return ast.NewBoolValue(isTruthy(right)), nil
	}

	panic(fmt.Sprintf("Unhandled operator %d in Unary Expression", expr.Token.Type))
}

func evalBinary(expr ast.Expr) (ast.LoxValue, error) {
	left, err := Evaluate(expr.Children[0])
	if err != nil {
		return ast.NewNilValue(), err
	}
	right, err := Evaluate(expr.Children[1])
	if err != nil {
		return ast.NewNilValue(), err
	}

	switch expr.Token.Type {
	case ast.MINUS:
		return ast.NewNumberValue(left.Number() - right.Number()), nil
	case ast.SLASH:
		return ast.NewNumberValue(left.Number() / right.Number()), nil
	case ast.STAR:
		return ast.NewNumberValue(left.Number() * right.Number()), nil
	case ast.PLUS:
		if left.Type == ast.LT_NUMBER && right.Type == ast.LT_NUMBER {
			return ast.NewNumberValue(left.Number() + right.Number()), nil
		} else if left.Type == ast.LT_STRING && right.Type == ast.LT_STRING {
			return ast.NewStringValue(left.String() + right.String()), nil
		} else {
			// TODO: Report Error, mismatching types for operand +
		}
	case ast.GREATER:
		return ast.NewBoolValue(left.Number() > right.Number()), nil
	case ast.GREATER_EQUAL:
		return ast.NewBoolValue(left.Number() >= right.Number()), nil
	case ast.LESS:
		return ast.NewBoolValue(left.Number() < right.Number()), nil
	case ast.LESS_EQUAL:
		return ast.NewBoolValue(left.Number() <= right.Number()), nil
	case ast.BANG_EQUAL:
		return ast.NewBoolValue(!isEqual(left, right)), nil
	case ast.EQUAL_EQUAL:
		return ast.NewBoolValue(isEqual(left, right)), nil
	}

	panic(fmt.Sprintf("Unhandled operator %d in Binary Expression", expr.Token.Type))
}

func evalGrouping(expr ast.Expr) (ast.LoxValue, error) {
	return Evaluate(expr.Children[0])
}

func isTruthy(value ast.LoxValue) bool {
	switch value.Type {
	case ast.LT_NIL:
		return false
	case ast.LT_BOOL:
		return value.Bool()
	}

	return true
}

func isEqual(left ast.LoxValue, right ast.LoxValue) bool {
	if left.Type == ast.LT_NIL && right.Type == ast.LT_NIL {
		return true
	}

	return left == right
}
