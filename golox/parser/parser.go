package parser

type Parser struct {
	tokens  []Token
	current int
}

/*
Expression grammar rules:

expression     → equality ;
equality       → comparison ( ( "!=" | "==" ) comparison )* ;
comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term           → factor ( ( "-" | "+" ) factor )* ;
factor         → unary ( ( "/" | "*" ) unary )* ;
unary          → ( "!" | "-" ) unary
               | primary ;
primary        → NUMBER | STRING | "true" | "false" | "nil"
               | "(" expression ")" ;

Every rule is implemented as a single parser function below
*/

// expression     → equality ;
func (p *Parser) expression() Expr {
	return p.equality()
}

// equality       → comparison ( ( "!=" | "==" ) comparison )* ;
func (p *Parser) equality() Expr {
	expr := p.comparison()

	for {
		match, operator := p.match(BANG_EQUAL, EQUAL_EQUAL)
		if !match {
			break
		}

		right := p.comparison()
		expr = NewBinaryExpr(expr, operator, right)
	}

	return expr
}

// comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
func (p *Parser) comparison() Expr {
	expr := p.term()

	for {
		match, operator := p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL)
		if !match {
			break
		}

		right := p.term()
		expr = NewBinaryExpr(expr, operator, right)
	}

	return expr
}

// term           → factor ( ( "-" | "+" ) factor )* ;
func (p *Parser) term() Expr {
	expr := p.factor()

	for {
		match, operator := p.match(MINUS, PLUS)
		if !match {
			break
		}

		right := p.factor()
		expr = NewBinaryExpr(expr, operator, right)
	}

	return expr
}

// factor         → unary ( ( "/" | "*" ) unary )* ;
func (p *Parser) factor() Expr {
	expr := p.unary()

	for {
		match, operator := p.match(SLASH, STAR)
		if !match {
			break
		}

		right := p.unary()
		expr = NewBinaryExpr(expr, operator, right)
	}

	return expr
}

// unary          → ( "!" | "-" ) unary | primary ;
func (p *Parser) unary() Expr {
	match, operator := p.match(BANG, MINUS)
	if match {
		child := p.unary()
		return NewUnaryExpr(child, operator)
	}

	return p.primary()
}

// primary        → NUMBER | STRING | "true" | "false" | "nil" | "(" expression ")" ;
func (p *Parser) primary() Expr {
	match, token := p.match(NUMBER, STRING, TRUE, FALSE, NIL, LEFT_PAREN)
	if !match {
		panic("Can't parse")
	}

	if token.type_ != LEFT_PAREN {
		return NewLiteralExpr(token)
	}

}

// Checks if the current Token is one of the given types and if so, consumes it
// Returns True and the matched Token if it was matched, false and the zero-value of Token
// if it wasn't matched
func (p *Parser) match(tokens ...TokenType) (bool, Token) {
	for _, token := range tokens {
		if p.check(token) {
			p.current += 1
			return true, p.tokens[p.current-1]
		}
	}

	return false, Token{}
}

// Checks if the current Token is of the given type
func (p Parser) check(token TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().type_ == token
}

func (p Parser) peek() Token {
	return p.tokens[p.current]
}

func (p Parser) isAtEnd() bool {
	return p.peek().type_ == EOF
}
