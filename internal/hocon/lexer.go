/*
 * lexer.go
 * This file is part of the gekka-config project.
 *
 * Copyright (c) 2026 Sopranoworks, Osamu Takahashi
 * SPDX-License-Identifier: MIT
 */
package hocon

import (
	"strings"
	"unicode"
)

// Scanner is responsible for tokenizing HOCON input.
type Scanner struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	line         int
	column       int
}

// NewScanner creates a new Scanner for the given input string.
func NewScanner(input string) *Scanner {
	s := &Scanner{input: input, line: 1, column: 0}
	s.readChar()
	return s
}

func (s *Scanner) readChar() {
	prevWasNewline := (s.ch == '\n')
	if s.readPosition >= len(s.input) {
		s.ch = 0
	} else {
		s.ch = s.input[s.readPosition]
	}
	s.position = s.readPosition
	s.readPosition++

	if prevWasNewline {
		s.line++
		s.column = 1
	} else {
		s.column++
	}
}

func (s *Scanner) peekChar() byte {
	if s.readPosition >= len(s.input) {
		return 0
	}
	return s.input[s.readPosition]
}

func (s *Scanner) skipWhitespace() {
	for s.ch == ' ' || s.ch == '\t' || s.ch == '\r' {
		s.readChar()
	}
}

// NextToken returns the next token from the input.
func (s *Scanner) NextToken() Token {
	var tok Token

	s.skipWhitespace()

	line := s.line
	col := s.column

	switch s.ch {
	case '{':
		tok = Token{Type: TOKEN_LEFT_BRACE, Literal: string(s.ch), Line: line, Column: col}
	case '}':
		tok = Token{Type: TOKEN_RIGHT_BRACE, Literal: string(s.ch), Line: line, Column: col}
	case '[':
		tok = Token{Type: TOKEN_LEFT_BRACKET, Literal: string(s.ch), Line: line, Column: col}
	case ']':
		tok = Token{Type: TOKEN_RIGHT_BRACKET, Literal: string(s.ch), Line: line, Column: col}
	case ':':
		tok = Token{Type: TOKEN_COLON, Literal: string(s.ch), Line: line, Column: col}
	case '=':
		tok = Token{Type: TOKEN_EQUALS, Literal: string(s.ch), Line: line, Column: col}
	case ',':
		tok = Token{Type: TOKEN_COMMA, Literal: string(s.ch), Line: line, Column: col}
	case '\n':
		tok = Token{Type: TOKEN_NEWLINE, Literal: "\n", Line: line, Column: col}
	case '#', '/':
		if s.ch == '/' && s.peekChar() != '/' {
			tok = s.readUnquotedStringToken(line, col)
			return tok
		}
		tok = Token{Type: TOKEN_COMMENT, Literal: s.readComment(), Line: line, Column: col}
		return tok
	case '"':
		if s.peekChar() == '"' && s.peekReadFar(2) == '"' {
			tok = Token{Type: TOKEN_STRING, Literal: s.readTripleQuotedString(), Line: line, Column: col}
		} else {
			tok = Token{Type: TOKEN_STRING, Literal: s.readQuotedString(), Line: line, Column: col}
		}
		return tok
	case '$':
		if s.peekChar() == '{' {
			tok = Token{Type: TOKEN_SUBSTITUTION, Literal: s.readSubstitution(), Line: line, Column: col}
			return tok
		}
		tok = s.readUnquotedStringToken(line, col)
		return tok
	case 0:
		tok.Literal = ""
		tok.Type = TOKEN_EOF
		tok.Line = line
		tok.Column = col
	default:
		if isUnquotedStringChar(s.ch) {
			tok = s.readUnquotedStringToken(line, col)
			return tok
		}
		tok = Token{Type: TOKEN_ILLEGAL, Literal: string(s.ch), Line: line, Column: col}
	}

	s.readChar()
	return tok
}

func (s *Scanner) readComment() string {
	startPos := s.position
	for s.ch != '\n' && s.ch != 0 {
		s.readChar()
	}
	return s.input[startPos:s.position]
}

func (s *Scanner) readQuotedString() string {
	s.readChar() // skip "
	startPos := s.position
	for s.ch != '"' && s.ch != 0 {
		if s.ch == '\\' {
			s.readChar() // skip \
		}
		s.readChar()
	}
	content := s.input[startPos:s.position]
	if s.ch == '"' {
		s.readChar() // skip closing "
	}
	return content
}

func (s *Scanner) readTripleQuotedString() string {
	s.readChar() // skip "
	s.readChar() // skip "
	s.readChar() // skip "
	startPos := s.position
	for {
		if s.ch == 0 {
			break
		}
		if s.ch == '"' && s.peekChar() == '"' && s.peekReadFar(2) == '"' {
			break
		}
		s.readChar()
	}
	content := s.input[startPos:s.position]
	if s.ch == '"' {
		s.readChar() // skip "
		s.readChar() // skip "
		s.readChar() // skip "
	}
	return content
}

func (s *Scanner) readUnquotedStringToken(line, col int) Token {
	startPos := s.position
	for isUnquotedStringChar(s.ch) {
		s.readChar()
	}
	return Token{
		Type:    TOKEN_STRING,
		Literal: strings.TrimSpace(s.input[startPos:s.position]),
		Line:    line,
		Column:  col,
	}
}

func (s *Scanner) readSubstitution() string {
	startPos := s.position
	s.readChar() // skip $
	s.readChar() // skip {
	for s.ch != '}' && s.ch != 0 {
		s.readChar()
	}
	if s.ch == '}' {
		s.readChar() // skip closing }
	}
	return s.input[startPos:s.position]
}

func (s *Scanner) peekReadFar(n int) byte {
	pos := s.readPosition + n - 1
	if pos >= len(s.input) {
		return 0
	}
	return s.input[pos]
}

func isUnquotedStringChar(ch byte) bool {
	if ch == 0 || unicode.IsSpace(rune(ch)) {
		return false
	}
	switch ch {
	case '{', '}', '[', ']', ':', '=', ',', '+', '#', '`', '^', '?', '!', '@', '*', '&', '/', '\\':
		return false
	case '$':
		return false // handled separately for substitution
	case '"':
		return false // handled separately for quoted
	}
	return true
}
