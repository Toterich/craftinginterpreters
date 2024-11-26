package ast

import "strconv"

type LoxType int

// All supported types in Lox
const (
	LT_NIL LoxType = iota
	LT_STRING
	LT_NUMBER // 64bit float
	LT_BOOL
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
	default:
		panic("Incomplete Switch")
	}
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

func (v LoxValue) AsString() string {
	return v.Value.(string)
}

func (v LoxValue) AsNumber() float64 {
	return v.Value.(float64)
}

func (v LoxValue) AsBool() bool {
	return v.Value.(bool)
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
	default:
		panic("Incomplete Switch")
	}
}
