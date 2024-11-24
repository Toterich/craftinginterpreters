package parser

import (
	"toterich/golox/util"
)

// A parserError, which in addition to an error string contains the Token where the error occured
type parserError interface {
	Token() Token
	error
}

type parserErrorImpl struct {
	token Token
	msg   string
}

func (pe parserErrorImpl) Token() Token {
	return pe.token
}

func (pe parserErrorImpl) Error() string {
	return pe.msg
}

func newParserError(token Token, msg string) parserError {
	return parserErrorImpl{token: token, msg: msg}
}

// A recursive-descent parser for transforming a stream of Tokens into an AST
type Parser struct {
	tokens  []Token
	current int
}

func (p *Parser) Parse(input []Token) Expr {
	p.tokens = input
	p.current = 0

	expr, err := p.parseExpression()
	if err != nil {
		util.LogError(err.Token().line, err.Error())
	}

	return expr
}

/*
Expression grammar rules:

expression     -> comma_op
comma_op       -> equality ("," equality)* ;
equality       -> comparison ( ( "!=" | "==" ) comparison )* ;
comparison     -> term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term           -> factor ( ( "-" | "+" ) factor )* ;
factor         -> unary ( ( "/" | "*" ) unary )* ;
unary          -> ( "!" | "-" ) unary
               | primary ;
primary        -> NUMBER | STRING | "true" | "false" | "nil"
               | "(" expression ")" ;

Every rule is implemented as a single parser function below
*/

// expression     -> comma_op
func (p *Parser) parseExpression() (Expr, parserError) {
	return p.parseCommaOp()
}

// comma_op       -> equality ("," equality)* ;
func (p *Parser) parseCommaOp() (Expr, parserError) {
	expr, err := p.parseEquality()
	if err != nil {
		return Expr{}, err
	}

	for p.match(COMMA) {
		operator := p.previous()

		right, err := p.parseEquality()
		if err != nil {
			return expr, err
		}
		expr = NewBinaryExpr(expr, operator, right)
	}

	return expr, nil
}

// equality       → comparison ( ( "!=" | "==" ) comparison )* ;
func (p *Parser) parseEquality() (Expr, parserError) {
	expr, err := p.parseComparison()
	if err != nil {
		return Expr{}, err
	}

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		operator := p.previous()

		right, err := p.parseComparison()
		if err != nil {
			return expr, err
		}
		expr = NewBinaryExpr(expr, operator, right)
	}

	return expr, nil
}

// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
func (p *Parser) parseComparison() (Expr, parserError) {
	expr, err := p.parseTerm()
	if err != nil {
		return Expr{}, err
	}

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		operator := p.previous()

		right, err := p.parseTerm()
		if err != nil {
			return expr, err
		}
		expr = NewBinaryExpr(expr, operator, right)
	}

	return expr, nil
}

// term           → factor ( ( "-" | "+" ) factor )* ;
func (p *Parser) parseTerm() (Expr, parserError) {
	expr, err := p.parseFactor()
	if err != nil {
		return Expr{}, err
	}

	for p.match(MINUS, PLUS) {
		operator := p.previous()

		right, err := p.parseFactor()
		if err != nil {
			return expr, err
		}
		expr = NewBinaryExpr(expr, operator, right)
	}

	return expr, nil
}

// factor         → unary ( ( "/" | "*" ) unary )* ;
func (p *Parser) parseFactor() (Expr, parserError) {
	expr, err := p.parseUnary()
	if err != nil {
		return Expr{}, err
	}

	for p.match(SLASH, STAR) {
		operator := p.previous()

		right, err := p.parseUnary()
		if err != nil {
			return expr, err
		}

		expr = NewBinaryExpr(expr, operator, right)
	}

	return expr, nil
}

// unary          → ( "!" | "-" ) parseUnary | primary ;
func (p *Parser) parseUnary() (Expr, parserError) {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		child, err := p.parseUnary()
		if err != nil {
			return Expr{}, err
		}
		return NewUnaryExpr(child, operator), nil
	}

	return p.parsePrimary()
}

// primary        → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" ;
func (p *Parser) parsePrimary() (Expr, parserError) {
	if !p.match(NUMBER, STRING, TRUE, FALSE, NIL, LEFT_PAREN) {
		return Expr{}, newParserError(p.peek(), "Expected Expression")
	}

	token := p.previous()

	if token.type_ == LEFT_PAREN {
		expr, err := p.parseExpression()
		if err != nil {
			return Expr{}, err
		}
		_, err = p.consume(RIGHT_PAREN, "Expected ')' after expression.")
		return NewGroupingExpr(expr), err
	}

	return NewLiteralExpr(token), nil
}

// Checks if the current Token is one of the given types and if so, consumes it and returns true
func (p *Parser) match(tokens ...TokenType) bool {
	for _, token := range tokens {
		if p.check(token) {
			p.current += 1
			return true
		}
	}

	return false
}

// If the current Token is of the given type, consumes and returns it. Otherwise, returns the current Token without
// consuming it and an error containing the given message.
func (p *Parser) consume(token TokenType, errMsg string) (Token, parserError) {
	if p.match(token) {
		return p.previous(), nil
	}

	return p.peek(), newParserError(p.peek(), errMsg)
}

// Returns true if the current Token is of the given TokenType
func (p Parser) check(token TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().type_ == token
}

// Returns the next Token without consuming it
func (p Parser) peek() Token {
	return p.tokens[p.current]
}

// Returns the previously consumed Token
func (p Parser) previous() Token {
	util.Assert(p.current > 0, "This may only be called when p.current > 0")
	return p.tokens[p.current-1]
}

// Returns true if all Tokens have been consumed
func (p Parser) isAtEnd() bool {
	return p.peek().type_ == EOF
}

// Skips all tokens until the beginning of the next statement. This is required e.g. after a parsing error,
// where we discard the rest of the current statement.
func (p *Parser) skipToNextStatement() {
	p.current += 1

	for !p.isAtEnd() {
		// Statement ends after semicolon
		if p.previous().type_ == SEMICOLON {
			return
		}

		// Statement begins with one of these keywords
		switch p.peek().type_ {
		case CLASS:
			fallthrough
		case FUN:
			fallthrough
		case VAR:
			fallthrough
		case FOR:
			fallthrough
		case IF:
			fallthrough
		case WHILE:
			fallthrough
		case PRINT:
			fallthrough
		case RETURN:
			return
		}

		p.current += 1
	}
}
