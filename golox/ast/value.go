package ast

import (
	"strconv"
	"toterich/golox/util/assert"
)

type LoxType int

// All supported types in Lox
const (
	LT_NIL LoxType = iota
	LT_STRING
	LT_NUMBER // 64bit float
	LT_BOOL
	LT_FUNCTION
)

func (t LoxType) String() string {
	switch t {
	case LT_NIL:
		return "nil"
	case LT_STRING:
		return "String"
	case LT_NUMBER:
		return "Number"
	case LT_BOOL:
		return "Bool"
	case LT_FUNCTION:
		return "Function"
	default:
		panic(assert.MissingCase(t))
	}
}

type LoxFunction struct {
	Declaration Stmt
}

func (lf LoxFunction) Arity() int {
	return len(lf.Declaration.Tokens) - 1
}

// A Value in Lox, represented by a type and a pointer to the actual value.
// Use Type Assertions (see below) to extract the value
type LoxValue struct {
	Type  LoxType
	Value any
}

func NewNilValue() LoxValue {
	return LoxValue{Type: LT_NIL}
}

func NewStringValue(str string) LoxValue {
	return LoxValue{Type: LT_STRING, Value: str}
}

func NewNumberValue(num float64) LoxValue {
	return LoxValue{Type: LT_NUMBER, Value: num}
}

func NewBoolValue(val bool) LoxValue {
	return LoxValue{Type: LT_BOOL, Value: val}
}

func NewFunction(fun LoxFunction) LoxValue {
	return LoxValue{Type: LT_FUNCTION, Value: fun}
}

func (v LoxValue) IsTruthy() bool {
	switch v.Type {
	case LT_NIL:
		return false
	case LT_BOOL:
		return v.AsBool()
	}

	return true
}

func (v LoxValue) IsEqual(other LoxValue) bool {
	if v.Type == LT_NIL && other.Type == LT_NIL {
		return true
	}

	return v == other
}

func (v LoxValue) AsString() string {
	return v.Value.(string)
}

func (v LoxValue) AsNumber() float64 {
	return v.Value.(float64)
}

func (v LoxValue) AsBool() bool {
	return v.Value.(bool)
}

func (v LoxValue) AsFunction() LoxFunction {
	return v.Value.(LoxFunction)
}

// String representation of the LoxValue, don't confuse with AsString()!
func (v LoxValue) String() string {
	switch v.Type {
	case LT_NIL:
		return "nil"
	case LT_BOOL:
		return strconv.FormatBool(v.AsBool())
	case LT_NUMBER:
		return strconv.FormatFloat(v.AsNumber(), 'g', -1, 64)
	case LT_STRING:
		return v.AsString()
	case LT_FUNCTION:
		return v.AsFunction().Declaration.Tokens[0].Lexeme
	default:
		panic(assert.MissingCase(v.Type))
	}
}
