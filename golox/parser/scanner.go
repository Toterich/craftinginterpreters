package parser

import (
	"strconv"
	"toterich/golox/error"
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
	tokens  []Token
}

func (s *Scanner) ScanTokens(source string) []Token {
	s.start = 0
	s.current = 0
	s.line = 1
	s.source = source
	s.tokens = make([]Token, 0)

	for !s.isAtEnd() {
		s.start = s.current
		c := source[s.current]
		s.current += 1

		switch c {
		case '(':
			s.addToken(LEFT_PAREN)
		case ')':
			s.addToken(RIGHT_PAREN)
		case '{':
			s.addToken(LEFT_BRACE)
		case '}':
			s.addToken(RIGHT_BRACE)
		case ',':
			s.addToken(COMMA)
		case '.':
			s.addToken(DOT)
		case ';':
			s.addToken(SEMICOLON)
		case '-':
			s.addToken(MINUS)
		case '+':
			s.addToken(PLUS)
		case '*':
			s.addToken(STAR)
		case '!':
			if s.match('=') {
				s.addToken(BANG_EQUAL)
			} else {
				s.addToken(BANG)
			}
		case '=':
			if s.match('=') {
				s.addToken(EQUAL_EQUAL)
			} else {
				s.addToken(EQUAL)
			}
		case '<':
			if s.match('=') {
				s.addToken(LESS_EQUAL)
			} else {
				s.addToken(LESS)
			}
		case '>':
			if s.match('=') {
				s.addToken(GREATER_EQUAL)
			} else {
				s.addToken(GREATER)
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
				s.addToken(SLASH)
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
				error.LogError(s.line, "Unexpected character.")
			}
		}
	}

	s.tokens = append(s.tokens, Token{type_: EOF, lexeme: "", line: s.line})
	return s.tokens
}

func (s Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s Scanner) generateToken(type_ TokenType) Token {
	return Token{type_: type_, lexeme: s.source[s.start:s.current], line: s.line}
}

func (s *Scanner) addToken(type_ TokenType) {
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
		error.LogError(s.line, "Unterminated string.")
		return
	}

	// Consume closing '"'
	s.current += 1

	t := s.generateToken(STRING)
	// Store normalized String with Token
	t.literal = StringLiteral(s.source[s.start+1 : s.current-1])
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

	t := s.generateToken(NUMBER)

	// Store actual numberic value with Token
	num, err := strconv.ParseFloat(t.lexeme, 64)
	// If this triggers, the number parsing above has a bug
	error.AssertNoError(err)

	t.literal = NumberLiteral(num)

	s.tokens = append(s.tokens, t)
}

func (s *Scanner) matchIdentifier() {
	for isAlphaNumeric(s.peek()) {
		s.current += 1
	}

	t := s.generateToken(IDENTIFIER)

	// Check if this identifier is a keyword
	tokenType, ok := KeywordStrings[t.lexeme]
	if ok {
		t.type_ = tokenType
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
