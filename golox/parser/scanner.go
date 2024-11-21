package parser

import "toterich/golox/error"

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
	s.line = 0
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
			// Comments
			if s.match('/') {
				for (s.peek() != '\n') && (!s.isAtEnd()) {
					s.current += 1
				}
			} else {
				s.addToken(SLASH)
			}
		case '\n':
			s.line += 1
		case '"':
			s.stringLiteral()

		// Ignore whitespace
		case ' ':
		case '\r':
		case '\t':

		default:
			error.LogError(s.line, "Unexpected character.")
		}
	}

	s.tokens = append(s.tokens, Token{type_: EOF, lexeme: "", line: s.line})
	return s.tokens
}

func (s Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) addToken(type_ TokenType) {
	s.tokens = append(s.tokens, Token{type_: type_, lexeme: s.source[s.start:s.current], line: s.line})
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

func (s *Scanner) stringLiteral() {
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

	// Store String without wrapping ""
	s.tokens = append(s.tokens, Token{type_: STRING, lexeme: s.source[s.start+1 : s.current-1], line: s.line})
}
