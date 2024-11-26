package util

import (
	"errors"
	"fmt"
	"log"
	"toterich/golox/ast"
)

// A Lexical Error
type LexError struct {
	Line int
	Char byte
	Msg  string
}

func NewLexError(line int, char byte, msg string) LexError {
	return LexError{Line: line, Char: char, Msg: msg}
}

func (e LexError) Error() string {
	return fmt.Sprintf("Lexical Error at line %d: %s", e.Line, e.Msg)
}

// A SyntaxError, which in addition to an error string contains the Token where the error occured
type SyntaxError struct {
	Token ast.Token
	Msg   string
}

func NewSyntaxError(token ast.Token, msg string) SyntaxError {
	return SyntaxError{Token: token, Msg: msg}
}

func (e SyntaxError) Error() string {
	return fmt.Sprintf("Syntax Error at line %d: %s", e.Token.Line, e.Msg)
}

// A RuntimeError indicating an issue with executing Lox Code
type RuntimeError struct {
	Token ast.Token
	Msg   string
}

func NewRuntimeError(token ast.Token, msg string) RuntimeError {
	return RuntimeError{Token: token, Msg: msg}
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf("Runtime Error at line %d: %s", e.Token.Line, e.Msg)
}

func LogErrors(errs ...error) {
	for _, err := range errs {
		{
			var e LexError
			if errors.As(err, &e) {
				log.Printf("[line %d] Lexical Error at Char '%c': %s", e.Line, e.Char, e.Msg)
				continue
			}
		}

		{
			var e SyntaxError
			if errors.As(err, &e) {
				log.Printf("[line %d] Syntax Error at Token '%s': %s", e.Token.Line, e.Token.Lexeme, e.Msg)
				continue
			}
		}

		{
			var e RuntimeError
			if errors.As(err, &e) {
				log.Printf("[line %d] Runtime Error at '%s': %s", e.Token.Line, e.Token.Lexeme, e.Msg)
				continue
			}
		}

		log.Println(err)
	}
}
