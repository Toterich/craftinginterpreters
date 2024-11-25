package parser

import (
	"toterich/golox/ast"
	"toterich/golox/util"
)

// A parserError, which in addition to an error string contains the Token where the error occured
type parserError interface {
	Token() ast.Token
	error
}

type parserErrorImpl struct {
	token ast.Token
	msg   string
}

func (pe parserErrorImpl) Token() ast.Token {
	return pe.token
}

func (pe parserErrorImpl) Error() string {
	return pe.msg
}

func newParserError(token ast.Token, msg string) parserError {
	return parserErrorImpl{token: token, msg: msg}
}

// A recursive-descent parser for transforming a stream of Tokens into an AST
type Parser struct {
	tokens  []ast.Token
	current int
}

func (p *Parser) Parse(input []ast.Token) ast.Expr {
	p.tokens = input
	p.current = 0

	expr, err := p.parseExpression()
	if err != nil {
		util.LogError(err.Token().Line, err.Error())
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
func (p *Parser) parseExpression() (ast.Expr, parserError) {
	return p.parseCommaOp()
}

// comma_op       -> equality ("," equality)* ;
func (p *Parser) parseCommaOp() (ast.Expr, parserError) {
	expr, err := p.parseEquality()
	if err != nil {
		return ast.Expr{}, err
	}

	for p.match(ast.COMMA) {
		operator := p.previous()

		right, err := p.parseEquality()
		if err != nil {
			return expr, err
		}
		expr = ast.NewBinaryExpr(expr, operator, right)
	}

	return expr, nil
}

// equality       → comparison ( ( "!=" | "==" ) comparison )* ;
func (p *Parser) parseEquality() (ast.Expr, parserError) {
	expr, err := p.parseComparison()
	if err != nil {
		return ast.Expr{}, err
	}

	for p.match(ast.BANG_EQUAL, ast.EQUAL_EQUAL) {
		operator := p.previous()

		right, err := p.parseComparison()
		if err != nil {
			return expr, err
		}
		expr = ast.NewBinaryExpr(expr, operator, right)
	}

	return expr, nil
}

// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
func (p *Parser) parseComparison() (ast.Expr, parserError) {
	expr, err := p.parseTerm()
	if err != nil {
		return ast.Expr{}, err
	}

	for p.match(ast.GREATER, ast.GREATER_EQUAL, ast.LESS, ast.LESS_EQUAL) {
		operator := p.previous()

		right, err := p.parseTerm()
		if err != nil {
			return expr, err
		}
		expr = ast.NewBinaryExpr(expr, operator, right)
	}

	return expr, nil
}

// term           → factor ( ( "-" | "+" ) factor )* ;
func (p *Parser) parseTerm() (ast.Expr, parserError) {
	expr, err := p.parseFactor()
	if err != nil {
		return ast.Expr{}, err
	}

	for p.match(ast.MINUS, ast.PLUS) {
		operator := p.previous()

		right, err := p.parseFactor()
		if err != nil {
			return expr, err
		}
		expr = ast.NewBinaryExpr(expr, operator, right)
	}

	return expr, nil
}

// factor         → unary ( ( "/" | "*" ) unary )* ;
func (p *Parser) parseFactor() (ast.Expr, parserError) {
	expr, err := p.parseUnary()
	if err != nil {
		return ast.Expr{}, err
	}

	for p.match(ast.SLASH, ast.STAR) {
		operator := p.previous()

		right, err := p.parseUnary()
		if err != nil {
			return expr, err
		}

		expr = ast.NewBinaryExpr(expr, operator, right)
	}

	return expr, nil
}

// unary          → ( "!" | "-" ) parseUnary | primary ;
func (p *Parser) parseUnary() (ast.Expr, parserError) {
	if p.match(ast.BANG, ast.MINUS) {
		operator := p.previous()
		child, err := p.parseUnary()
		if err != nil {
			return ast.Expr{}, err
		}
		return ast.NewUnaryExpr(child, operator), nil
	}

	return p.parsePrimary()
}

// primary        → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" ;
func (p *Parser) parsePrimary() (ast.Expr, parserError) {
	if !p.match(ast.NUMBER, ast.STRING, ast.TRUE, ast.FALSE, ast.NIL, ast.LEFT_PAREN) {
		return ast.Expr{}, newParserError(p.peek(), "Expected Expression")
	}

	token := p.previous()

	if token.Type == ast.LEFT_PAREN {
		expr, err := p.parseExpression()
		if err != nil {
			return ast.Expr{}, err
		}
		_, err = p.consume(ast.RIGHT_PAREN, "Expected ')' after expression.")
		return ast.NewGroupingExpr(expr), err
	}

	return ast.NewLiteralExpr(token), nil
}

// Checks if the current Token is one of the given types and if so, consumes it and returns true
func (p *Parser) match(tokens ...ast.TokenType) bool {
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
func (p *Parser) consume(token ast.TokenType, errMsg string) (ast.Token, parserError) {
	if p.match(token) {
		return p.previous(), nil
	}

	return p.peek(), newParserError(p.peek(), errMsg)
}

// Returns true if the current ast.Token is of the given ast.TokenType
func (p Parser) check(token ast.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == token
}

// Returns the next ast.Token without consuming it
func (p Parser) peek() ast.Token {
	return p.tokens[p.current]
}

// Returns the previously consumed ast.Token
func (p Parser) previous() ast.Token {
	util.Assert(p.current > 0, "This may only be called when p.current > 0")
	return p.tokens[p.current-1]
}

// Returns true if all ast.Tokens have been consumed
func (p Parser) isAtEnd() bool {
	return p.peek().Type == ast.EOF
}

// Skips all tokens until the beginning of the next statement. This is required e.g. after a parsing error,
// where we discard the rest of the current statement.
func (p *Parser) skipToNextStatement() {
	p.current += 1

	for !p.isAtEnd() {
		// Statement ends after semicolon
		if p.previous().Type == ast.SEMICOLON {
			return
		}

		// Statement begins with one of these keywords
		switch p.peek().Type {
		case ast.CLASS:
			fallthrough
		case ast.FUN:
			fallthrough
		case ast.VAR:
			fallthrough
		case ast.FOR:
			fallthrough
		case ast.IF:
			fallthrough
		case ast.WHILE:
			fallthrough
		case ast.PRINT:
			fallthrough
		case ast.RETURN:
			return
		}

		p.current += 1
	}
}
