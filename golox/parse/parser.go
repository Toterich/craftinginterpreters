package parse

import (
	"toterich/golox/ast"
	"toterich/golox/util"
	"toterich/golox/util/assert"
)

// A recursive-descent parser for transforming a stream of Tokens into an AST
type Parser struct {
	tokens    []ast.Token
	errs      []error
	current   int
	loopLevel int
}

// For grammar rules, see lox_spec/grammar.txt
// Every rule is implemented in a parse... function below

// program        -> statement* EOF;
func (p *Parser) Parse(input []ast.Token) ([]ast.Stmt, []error) {
	p.tokens = input
	p.current = 0
	p.errs = nil
	p.loopLevel = 0

	var statements []ast.Stmt

	for !p.isAtEnd() {
		stmt, errs := p.parseDeclaration()
		// Discard statements with a parse error
		if errs != nil {
			p.errs = append(p.errs, errs...)
			p.skipToNextStatement()
		} else {
			statements = append(statements, stmt)
		}
	}

	return statements, p.errs
}

// declaration    -> varDeclStmt | statement ;
func (p *Parser) parseDeclaration() (ast.Stmt, []error) {
	var stmt ast.Stmt
	var err error
	if p.match(ast.VAR) {
		stmt, err = p.parseVarDecl()
		if err != nil {
			return stmt, []error{err}
		} else {
			return stmt, nil
		}
	}

	return p.parseStatement()
}

// statement
// -> exprStmt | ifStmt | printStmt | whileStmt | forStmt | breakStmt | blockStmt ;
func (p *Parser) parseStatement() (ast.Stmt, []error) {
	if p.match(ast.LEFT_BRACE) {
		return p.parseBlockStmt()
	}
	if p.match(ast.IF) {
		return p.parseIfStmt()
	}
	if p.match(ast.WHILE) {
		return p.parseWhileStmt()
	}
	if p.match(ast.FOR) {
		return p.parseForStmt()
	}

	// The following statements can only produce a single error each, which
	// is packed inside a single-element array

	var stmt ast.Stmt
	var err error
	if p.match(ast.PRINT) {
		stmt, err = p.parsePrintStmt()
	} else if p.match(ast.BREAK) {
		if p.loopLevel < 1 {
			err = util.NewSyntaxError(p.previous(), "break statement outside of loop.")
		} else {
			stmt = ast.NewBreakStmt()
			_, err = p.consume(ast.SEMICOLON, "expected ';' after break.")
		}
	} else {
		stmt, err = p.parseExprStmt()
	}

	if err != nil {
		return stmt, []error{err}
	} else {
		return stmt, nil
	}
}

// blockStmt      -> "{" statement* "}";
// parseBlockStmt() can return multiple errors because each nested statement can
// produce an error
func (p *Parser) parseBlockStmt() (ast.Stmt, []error) {
	children := make([]ast.Stmt, 0)
	var errs []error = nil

	// Empty blocks are allowed
	if p.match(ast.RIGHT_BRACE) {
		return ast.NewBlockStmt(children), errs
	}

	for !p.isAtEnd() {
		stmt, err := p.parseDeclaration()
		if err != nil {
			errs = append(errs, err...)
			// We continue parsing until the end of the block so the parser doesn't trip up.
			p.skipToNextStatement()
			// Check if we skipped over the end of the block
			if p.previous().Type == ast.RIGHT_BRACE {
				return ast.NewBlockStmt(children), errs
			}
			continue
		}
		children = append(children, stmt)
		if p.match(ast.RIGHT_BRACE) {
			return ast.NewBlockStmt(children), errs
		}
	}

	errs = append(errs, util.NewSyntaxError(p.peek(), "missing closing '}'."))
	return ast.NewBlockStmt(children), errs
}

// ifStmt         -> "if" "(" expression ")" statement ("else" statement)? ;
func (p *Parser) parseIfStmt() (ast.Stmt, []error) {
	if !p.match(ast.LEFT_PAREN) {
		return ast.NewInvalidStmt(),
			[]error{util.NewSyntaxError(p.peek(), "expected condition after 'if'.")}
	}

	condition, err := p.parseExpression()
	if err != nil {
		return ast.NewInvalidStmt(), []error{util.NewSyntaxError(p.peek(), "'if' condition must be a valid expression.")}
	}

	if !p.match(ast.RIGHT_PAREN) {
		return ast.NewInvalidStmt(),
			[]error{util.NewSyntaxError(p.peek(), "missing closing ')' after 'if' condition.")}
	}

	ifStmt, errs := p.parseStatement()
	if errs != nil {
		return ifStmt, errs
	}

	elseStmt := ast.NewInvalidStmt()

	if p.match(ast.ELSE) {
		elseStmt, errs = p.parseStatement()
		if errs != nil {
			return elseStmt, errs
		}
	}

	return ast.NewIfStmt(condition, ifStmt, elseStmt), nil
}

// whileStmt      -> "while" "(" expression ")" statement ;
func (p *Parser) parseWhileStmt() (ast.Stmt, []error) {
	p.incLoopLevel()
	defer p.decLoopLevel()

	if !p.match(ast.LEFT_PAREN) {
		return ast.NewInvalidStmt(),
			[]error{util.NewSyntaxError(p.peek(), "expected condition after 'while'.")}
	}

	condition, err := p.parseExpression()
	if err != nil {
		return ast.NewInvalidStmt(), []error{util.NewSyntaxError(p.peek(), "'while' condition must be a valid expression.")}
	}

	if !p.match(ast.RIGHT_PAREN) {
		return ast.NewInvalidStmt(),
			[]error{util.NewSyntaxError(p.peek(), "missing closing ')' after 'while' condition.")}
	}

	loopStmt, errs := p.parseStatement()
	if errs != nil {
		return loopStmt, errs
	}

	return ast.NewWhileStmt(condition, loopStmt), nil
}

// forStmt        -> "for" "(" (varDeclStmt | exprStmt | ";" ) expression? ";" expression? ")" statement ;
func (p *Parser) parseForStmt() (ast.Stmt, []error) {
	p.incLoopLevel()
	defer p.decLoopLevel()

	_, err := p.consume(ast.LEFT_PAREN, "expected condition after 'for'.")
	if err != nil {
		return ast.NewInvalidStmt(), []error{err}
	}

	// Parse the initializer, which can either be a variable declaration, an expression, or nothing
	var initializer ast.Stmt
	if p.match(ast.VAR) {
		initializer, err = p.parseVarDecl()
	} else if p.match(ast.SEMICOLON) {
		initializer, err = ast.NewInvalidStmt(), nil
	} else {
		initializer, err = p.parseExprStmt()
	}
	if err != nil {
		return initializer, []error{err}
	}

	// Parse the end condition, which is either an expression or nothing
	var condition ast.Expr
	if p.check(ast.SEMICOLON) {
		// If the condition is omitted, we do the same thing as C and replace it by a non-zero constant,
		// in this case a true literal. This means the loop will run indefinitely.
		condition, err =
			ast.NewLiteralExpr(
				ast.Token{Type: ast.TRUE, Literal: ast.NewBoolValue(true), Line: p.peek().Line}),
			nil
	} else {
		condition, err = p.parseExpression()
		if err != nil {
			return ast.NewExprStmt(condition), []error{err}
		}
	}
	_, err = p.consume(ast.SEMICOLON, "expected ';' after 'for' condition")
	if err != nil {
		return ast.NewExprStmt(condition), []error{err}
	}

	// Parse the increment, which is either an expression or nothing
	var increment ast.Expr
	if p.check(ast.RIGHT_PAREN) {
		increment, err = ast.Expr{}, nil
	} else {
		increment, err = p.parseExpression()
		if err != nil {
			return ast.NewExprStmt(increment), []error{err}
		}
	}
	_, err = p.consume(ast.RIGHT_PAREN, "expected closing ')' after 'for' increment")
	if err != nil {
		return ast.NewExprStmt(increment), []error{err}
	}

	body, errs := p.parseStatement()
	if errs != nil {
		return body, errs
	}

	// Desugar the for loop to a while statement by adding the initializer and increment as
	// statements
	while := ast.NewWhileStmt(condition, body)

	if increment.Type != ast.EXPR_INVALID {
		// If the loop body is already a block statement, append the increment to the end, otherwise create
		// a block statement of the original body and the increment
		if while.Children[0].Type == ast.ST_BLOCK {
			while.Children[0].Children = append(while.Children[0].Children, ast.NewExprStmt(increment))
		} else {
			while.Children[0] = ast.NewBlockStmt([]ast.Stmt{while.Children[0], ast.NewExprStmt(increment)})
		}
	}

	if initializer.Type != ast.ST_INVALID {
		// Wrap the whole while statement in a block and prepend the initializer
		while = ast.NewBlockStmt([]ast.Stmt{initializer, while})
	}

	return while, nil
}

// varDeclStmt    -> "var" IDENTIFIER ("=" expression)? ";";
func (p *Parser) parseVarDecl() (ast.Stmt, error) {
	if !p.match(ast.IDENTIFIER) {
		return ast.NewInvalidStmt(),
			util.NewSyntaxError(p.peek(), "expected Identifier after 'var'.")
	}

	stmt := ast.NewVarDeclStmt(p.previous().Lexeme)

	if p.match(ast.EQUAL) {
		expr, err := p.parseExpression()
		if err != nil {
			return ast.NewInvalidStmt(), err
		}
		stmt.Expr = expr
	}

	_, err := p.consume(ast.SEMICOLON, "expected ; after variable declaration.")
	return stmt, err
}

// printStmt      -> "print" expression ";"
func (p *Parser) parsePrintStmt() (ast.Stmt, error) {
	expr, err := p.parseExpression()
	if err != nil {
		return ast.NewInvalidStmt(), err
	}
	_, err = p.consume(ast.SEMICOLON, "expected ; after print statement.")
	return ast.NewPrintStmt(expr), err
}

// exprStmt       -> expression ";"
func (p *Parser) parseExprStmt() (ast.Stmt, error) {
	expr, err := p.parseExpression()
	if err != nil {
		return ast.NewInvalidStmt(), err
	}
	_, err = p.consume(ast.SEMICOLON, "expected ; after expression.")
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

// assignment     -> IDENTIFIER "=" assignment | logic_or;
func (p *Parser) parseAssignment() (ast.Expr, error) {
	// We parse the lhs of the assignment first as a general expression and only check if it is a
	// valid assignment target further below. This allows parsing complex l-values, e.g.
	// makeInst().foo.bar = val

	expr, err := p.parseLogicOr()
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

// logic_or       -> logic_and ("or" logic_and)* ;
func (p *Parser) parseLogicOr() (ast.Expr, error) {
	expr, err := p.parseLogicAnd()
	if err != nil {
		return expr, err
	}

	for p.match(ast.OR) {
		right, err := p.parseLogicAnd()
		if err != nil {
			return right, err
		}

		expr = ast.NewOrExpr(expr, right)
	}

	return expr, nil
}

// logic_and      -> equality ("and" equality)* ;
func (p *Parser) parseLogicAnd() (ast.Expr, error) {
	expr, err := p.parseEquality()
	if err != nil {
		return expr, err
	}

	for p.match(ast.AND) {
		right, err := p.parseEquality()
		if err != nil {
			return right, err
		}

		expr = ast.NewAndExpr(expr, right)
	}

	return expr, nil
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

// unary          -> ( "!" | "-" ) unary | call ;
func (p *Parser) parseUnary() (ast.Expr, error) {
	if p.match(ast.BANG, ast.MINUS) {
		operator := p.previous()
		child, err := p.parseUnary()
		if err != nil {
			return child, err
		}
		return ast.NewUnaryExpr(child, operator), nil
	}

	return p.parseCall()
}

// call           -> primary ( "(" arguments? ")" )* ;
// arguments      -> expression ( ", " expression )* ;
func (p *Parser) parseCall() (ast.Expr, error) {
	expr, err := p.parsePrimary()
	if err != nil {
		return expr, err
	}

	for p.match(ast.LEFT_PAREN) {
		args := make([]ast.Expr, 0)

		// empty argument list
		if p.match(ast.RIGHT_PAREN) {
			expr = ast.NewCallExpr(expr, p.previous(), args)
			continue
		}

		// first argument
		arg, err := p.parseExpression()
		if err != nil {
			return expr, err
		}
		args = append(args, arg)

		// additional arguments
		for p.match(ast.COMMA) {
			if len(args) >= 255 {
				return expr, util.NewSyntaxError(p.peek(), "can't have more than 255 arguments.")
			}
			arg, err := p.parseExpression()
			if err != nil {
				return expr, err
			}
			args = append(args, arg)
		}

		close, err := p.consume(ast.RIGHT_PAREN, "expected ')' after argument list.")
		if err != nil {
			return expr, err
		}

		expr = ast.NewCallExpr(expr, close, args)
	}

	return expr, nil
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

	return ast.Expr{}, util.NewSyntaxError(p.peek(), "expected expression.")
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
	assert.Assert(p.current > 0, "This may only be called when p.current > 0")
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
		// Statement ends after semicolon or }
		switch p.previous().Type {
		case ast.SEMICOLON:
			fallthrough
		case ast.RIGHT_BRACE:
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
			fallthrough
		case ast.LEFT_BRACE:
			return
		}

		p.current += 1
	}
}

func (p *Parser) incLoopLevel() {
	p.loopLevel += 1
}

func (p *Parser) decLoopLevel() {
	p.loopLevel -= 1
}
