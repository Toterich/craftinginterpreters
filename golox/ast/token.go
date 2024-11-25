package ast

import "fmt"

type TokenType int

const (
	// Single-char
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	// Single or Double-char
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// Literals
	IDENTIFIER
	STRING
	NUMBER

	// Keywords
	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	EOF
)

var KeywordStrings = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"fun":    FUN,
	"for":    FOR,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

type LoxType int

const (
	LT_NIL LoxType = iota
	LT_STRING
	LT_NUMBER
	LT_BOOL
)

type LiteralValue struct {
	Type  LoxType
	Value any
}

func NewStringLiteral(str string) LiteralValue {
	return LiteralValue{Type: LT_STRING, Value: str}
}

func NewNumberLiteral(num float64) LiteralValue {
	return LiteralValue{Type: LT_NUMBER, Value: num}
}

func NewBoolLiteral(val bool) LiteralValue {
	return LiteralValue{Type: LT_BOOL, Value: val}
}

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal LiteralValue
	Line    int
}

func (t Token) String() string {
	return fmt.Sprintf("%d: %s (%d)", t.Type, t.Lexeme, t.Line)
}
