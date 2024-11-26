package parse

import (
	"toterich/golox/ast"
	"toterich/golox/util"
)

// A recursive-descent parser for transforming a stream of Tokens into an AST
type Parser struct {
	tokens  []ast.Token
	errs    []error
	current int
}

// For grammar rules, see lox_spec/grammar.txt
// Every rule is implemented in a parse... function below

// program        -> statement* EOF;
func (p *Parser) Parse(input []ast.Token) ([]ast.Stmt, []error) {
	p.tokens = input
	p.current = 0
	p.errs = nil

	var statements []ast.Stmt

	for !p.isAtEnd() {
		stmt, err := p.parseStatement()
		// Discard statements with a parse error
		if err != nil {
			p.errs = append(p.errs, err)
			p.skipToNextStatement()
		} else {
			statements = append(statements, stmt)
		}
	}

	return statements, p.errs
}

// statement      -> exprStmt | printStmt | varDeclStmt;
func (p *Parser) parseStatement() (ast.Stmt, error) {
	if p.match(ast.VAR) {
		return p.parseVarDeclStmt()
	}
	if p.match(ast.PRINT) {
		return p.parsePrintStmt()
	}

	return p.parseExprStmt()
}

// varDeclStmt    -> "var" IDENTIFIER ("=" expression)? ";";
func (p *Parser) parseVarDeclStmt() (ast.Stmt, error) {
	if !p.match(ast.IDENTIFIER) {
		return ast.NewInvalidStmt(),
			util.NewSyntaxError(p.peek(), "expected Identifier after 'var'")
	}

	stmt := ast.NewVarDeclStmt(p.previous().Lexeme)

	if p.match(ast.EQUAL) {
		expr, err := p.parseExpression()
		if err != nil {
			return ast.NewInvalidStmt(), err
		}
		stmt.Expr = expr
	}

	_, err := p.consume(ast.SEMICOLON, "expected ; after variable declaration")
	return stmt, err
}

// printStmt      -> "print" expression ";"
func (p *Parser) parsePrintStmt() (ast.Stmt, error) {
	expr, err := p.parseExpression()
	if err != nil {
		return ast.NewInvalidStmt(), err
	}
	_, err = p.consume(ast.SEMICOLON, "expected ; after print expression")
	return ast.NewPrintStmt(expr), err
}

// exprStmt       -> expression ";"
func (p *Parser) parseExprStmt() (ast.Stmt, error) {
	expr, err := p.parseExpression()
	if err != nil {
		return ast.NewInvalidStmt(), err
	}
	_, err = p.consume(ast.SEMICOLON, "expected ; after expression")
	return ast.NewExprStmt(expr), err
}

// expression     -> comma_op
func (p *Parser) parseExpression() (ast.Expr, error) {
	return p.parseCommaOp()
}

// comma_op       -> assignment ("," assignment)* ;
func (p *Parser) parseCommaOp() (ast.Expr, error) {
	expr, err := p.parseAssignment()
	if err != nil {
		return expr, err
	}

	for p.match(ast.COMMA) {
		operator := p.previous()

		right, err := p.parseAssignment()
		if err != nil {
			return expr, err
		}
		expr = ast.NewBinaryExpr(expr, operator, right)
	}

	return expr, nil
}

// assignment     -> IDENTIFIER "=" assignment | equality;
func (p *Parser) parseAssignment() (ast.Expr, error) {
	// We parse the lhs of the assignment first as a general expression and only check if it is a
	// valid assignment target further below. This allows parsing complex l-values, e.g.
	// makeInst().foo.bar = val

	expr, err := p.parseEquality()
	if err != nil {
		return expr, err
	}

	if p.match(ast.EQUAL) {
		equals := p.previous()
		right, err := p.parseAssignment()
		if err != nil {
			return right, err
		}

		if expr.Type == ast.EXPR_IDENTIFIER {
			left := expr.Token
			return ast.NewAssignExpr(left, right), nil
		}

		return expr, util.NewRuntimeError(equals, "lhs of assignment is not an identifier.")
	}

	return expr, err
}

// equality       → comparison ( ( "!=" | "==" ) comparison )* ;
func (p *Parser) parseEquality() (ast.Expr, error) {
	expr, err := p.parseComparison()
	if err != nil {
		return expr, err
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
func (p *Parser) parseComparison() (ast.Expr, error) {
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
func (p *Parser) parseTerm() (ast.Expr, error) {
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
func (p *Parser) parseFactor() (ast.Expr, error) {
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
func (p *Parser) parseUnary() (ast.Expr, error) {
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

// primary        → NUMBER | STRING | IDENTIFIER | "true" | "false" | "nil" | "(" expression ")" ;
func (p *Parser) parsePrimary() (ast.Expr, error) {
	if p.match(ast.NUMBER, ast.STRING, ast.TRUE, ast.FALSE, ast.NIL) {
		return ast.NewLiteralExpr(p.previous()), nil
	}

	if p.match(ast.IDENTIFIER) {
		return ast.NewIdentifierExpression(p.previous()), nil
	}

	if p.match(ast.LEFT_PAREN) {
		expr, err := p.parseExpression()
		if err != nil {
			return ast.Expr{}, err
		}
		_, err = p.consume(ast.RIGHT_PAREN, "expected ')' after expression.")
		return ast.NewGroupingExpr(expr), err
	}

	return ast.Expr{}, util.NewSyntaxError(p.peek(), "expected Expression.")
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
func (p *Parser) consume(token ast.TokenType, errMsg string) (ast.Token, error) {
	if p.match(token) {
		return p.previous(), nil
	}

	return p.peek(), util.NewSyntaxError(p.peek(), errMsg)
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
	if p.isAtEnd() {
		return
	}

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
