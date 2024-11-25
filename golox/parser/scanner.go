package parser

import (
	"strconv"
	"toterich/golox/ast"
	"toterich/golox/util"
)

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z' || c == '_')
}

func isAlphaNumeric(c byte) bool {
	return isDigit(c) || isAlpha(c)
}

type Scanner struct {
	start   int
	current int
	line    int
	source  string
	tokens  []ast.Token
}

func (s *Scanner) ScanTokens(source string) []ast.Token {
	s.start = 0
	s.current = 0
	s.line = 1
	s.source = source
	s.tokens = make([]ast.Token, 0)

	for !s.isAtEnd() {
		s.start = s.current
		c := source[s.current]
		s.current += 1

		switch c {
		case '(':
			s.addToken(ast.LEFT_PAREN)
		case ')':
			s.addToken(ast.RIGHT_PAREN)
		case '{':
			s.addToken(ast.LEFT_BRACE)
		case '}':
			s.addToken(ast.RIGHT_BRACE)
		case ',':
			s.addToken(ast.COMMA)
		case '.':
			s.addToken(ast.DOT)
		case ';':
			s.addToken(ast.SEMICOLON)
		case '-':
			s.addToken(ast.MINUS)
		case '+':
			s.addToken(ast.PLUS)
		case '*':
			s.addToken(ast.STAR)
		case '!':
			if s.match('=') {
				s.addToken(ast.BANG_EQUAL)
			} else {
				s.addToken(ast.BANG)
			}
		case '=':
			if s.match('=') {
				s.addToken(ast.EQUAL_EQUAL)
			} else {
				s.addToken(ast.EQUAL)
			}
		case '<':
			if s.match('=') {
				s.addToken(ast.LESS_EQUAL)
			} else {
				s.addToken(ast.LESS)
			}
		case '>':
			if s.match('=') {
				s.addToken(ast.GREATER_EQUAL)
			} else {
				s.addToken(ast.GREATER)
			}
		case '/':
			if s.match('/') {
				// Comments
				for (s.peek() != '\n') && (!s.isAtEnd()) {
					s.current += 1
				}
			} else if s.match('*') {
				// Block Comments
				s.matchBlockComment()
			} else {
				s.addToken(ast.SLASH)
			}
		case '\n':
			s.line += 1
		case '"':
			s.matchString()

		// Ignore whitespace
		case ' ':
		case '\r':
		case '\t':

		default:
			if isDigit(c) {
				s.matchNumber()
			} else if isAlpha(c) {
				s.matchIdentifier()
			} else {
				util.LogError(s.line, "Unexpected character.")
			}
		}
	}

	s.tokens = append(s.tokens, ast.Token{Type: ast.EOF, Lexeme: "", Line: s.line})
	return s.tokens
}

func (s Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s Scanner) generateToken(type_ ast.TokenType) ast.Token {
	return ast.Token{Type: type_, Lexeme: s.source[s.start:s.current], Line: s.line}
}

func (s *Scanner) addToken(type_ ast.TokenType) {
	s.tokens = append(s.tokens, s.generateToken(type_))
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}

	s.current += 1
	return true
}

func (s Scanner) peek() byte {
	if s.isAtEnd() {
		return '\x00'
	}
	return s.source[s.current]
}

func (s Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return '\x00'
	}
	return s.source[s.current+1]
}

func (s *Scanner) matchString() {
	for (s.peek() != '"') && (!s.isAtEnd()) {
		if s.peek() == '\n' {
			s.line += 1
		}
		s.current += 1
	}

	if s.isAtEnd() {
		util.LogError(s.line, "Unterminated string.")
		return
	}

	// Consume closing '"'
	s.current += 1

	t := s.generateToken(ast.STRING)
	// Store normalized String with ast.Token
	t.Literal = ast.NewStringValue(s.source[s.start+1 : s.current-1])
	s.tokens = append(s.tokens, t)
}

func (s *Scanner) matchNumber() {
	for isDigit(s.peek()) {
		s.current += 1
	}

	if s.peek() == '.' && isDigit(s.peekNext()) {
		// Consume '.'
		s.current += 1

		// Consume all fractional digits
		for isDigit(s.peek()) {
			s.current += 1
		}
	}

	t := s.generateToken(ast.NUMBER)

	// Store actual numberic value with ast.Token
	num, err := strconv.ParseFloat(t.Lexeme, 64)
	// If this triggers, the number parsing above has a bug
	util.AssertNoError(err)

	t.Literal = ast.NewNumberValue(num)

	s.tokens = append(s.tokens, t)
}

func (s *Scanner) matchIdentifier() {
	for isAlphaNumeric(s.peek()) {
		s.current += 1
	}

	t := s.generateToken(ast.IDENTIFIER)

	// Check if this identifier is a keyword
	tokenType, ok := ast.KeywordStrings[t.Lexeme]
	if ok {
		t.Type = tokenType
		if tokenType == ast.TRUE {
			t.Literal = ast.NewBoolValue(true)
		} else if tokenType == ast.FALSE {
			t.Literal = ast.NewBoolValue(false)
		}
	}

	s.tokens = append(s.tokens, t)
}

func (s *Scanner) matchBlockComment() {
	nestingLevel := 1

	for nestingLevel > 0 {
		if s.peek() == '*' && s.peekNext() == '/' {
			nestingLevel -= 1
			s.current += 2
		} else if s.peek() == '/' && s.peekNext() == '*' {
			nestingLevel += 1
			s.current += 2
		} else if s.peek() == '\n' {
			s.line += 1
			s.current += 1
		} else {
			s.current += 1
		}
	}
}
