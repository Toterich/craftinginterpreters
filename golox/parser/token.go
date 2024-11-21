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

type Token struct {
	type_   TokenType
	lexeme  string
	literal LiteralValue
	line    int
}

func (t Token) ToString() string {
	return fmt.Sprintf("%d %s", t.type_, t.lexeme)
}
