package parser

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

type LiteralType int

const (
	LT_NIL LiteralType = iota
	LT_STRING
	LT_NUMBER
)

type LiteralValue struct {
	tag         LiteralType
	stringValue string
	numberValue float64
}

func StringLiteral(str string) LiteralValue {
	return LiteralValue{tag: LT_STRING, stringValue: str}
}

func NumberLiteral(num float64) LiteralValue {
	return LiteralValue{tag: LT_NUMBER, numberValue: num}
}

type Token struct {
	type_   TokenType
	lexeme  string
	literal LiteralValue
	line    int
}

func (t Token) String() string {
	return fmt.Sprintf("%d: %s (%d)", t.type_, t.lexeme, t.line)
}
