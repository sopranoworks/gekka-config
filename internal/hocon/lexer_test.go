/*
 * lexer_test.go
 * This file is part of the gekka-config project.
 *
 * Copyright (c) 2026 Sopranoworks, Osamu Takahashi
 * SPDX-License-Identifier: MIT
 */
package hocon

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `
{
  "akka" : {
    "actor" : {
      provider = "cluster"
      timeout = 10s # this is a comment
      // double slash comment
    }
  }
  list = [1, 2, 3]
  multi = """
    hello
    world
  """
  sub = ${path.to.value}
}
`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{TOKEN_NEWLINE, "\n"},
		{TOKEN_LEFT_BRACE, "{"},
		{TOKEN_NEWLINE, "\n"},
		{TOKEN_STRING, "akka"},
		{TOKEN_COLON, ":"},
		{TOKEN_LEFT_BRACE, "{"},
		{TOKEN_NEWLINE, "\n"},
		{TOKEN_STRING, "actor"},
		{TOKEN_COLON, ":"},
		{TOKEN_LEFT_BRACE, "{"},
		{TOKEN_NEWLINE, "\n"},
		{TOKEN_STRING, "provider"},
		{TOKEN_EQUALS, "="},
		{TOKEN_STRING, "cluster"},
		{TOKEN_NEWLINE, "\n"},
		{TOKEN_STRING, "timeout"},
		{TOKEN_EQUALS, "="},
		{TOKEN_STRING, "10s"},
		{TOKEN_COMMENT, "# this is a comment"},
		{TOKEN_NEWLINE, "\n"},
		{TOKEN_COMMENT, "// double slash comment"},
		{TOKEN_NEWLINE, "\n"},
		{TOKEN_RIGHT_BRACE, "}"},
		{TOKEN_NEWLINE, "\n"},
		{TOKEN_RIGHT_BRACE, "}"},
		{TOKEN_NEWLINE, "\n"},
		{TOKEN_STRING, "list"},
		{TOKEN_EQUALS, "="},
		{TOKEN_LEFT_BRACKET, "["},
		{TOKEN_STRING, "1"},
		{TOKEN_COMMA, ","},
		{TOKEN_STRING, "2"},
		{TOKEN_COMMA, ","},
		{TOKEN_STRING, "3"},
		{TOKEN_RIGHT_BRACKET, "]"},
		{TOKEN_NEWLINE, "\n"},
		{TOKEN_STRING, "multi"},
		{TOKEN_EQUALS, "="},
		{TOKEN_STRING, "\n    hello\n    world\n  "},
		{TOKEN_NEWLINE, "\n"},
		{TOKEN_STRING, "sub"},
		{TOKEN_EQUALS, "="},
		{TOKEN_SUBSTITUTION, "${path.to.value}"},
		{TOKEN_NEWLINE, "\n"},
		{TOKEN_RIGHT_BRACE, "}"},
		{TOKEN_NEWLINE, "\n"},
		{TOKEN_EOF, ""},
	}

	s := NewScanner(input)

	for i, tt := range tests {
		tok := s.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestCoordinateTracking(t *testing.T) {
	input := `a
bc`
	s := NewScanner(input)

	tok1 := s.NextToken() // 'a'
	if tok1.Line != 1 || tok1.Column != 1 {
		t.Errorf("tok1 wrong position: line %d, col %d. expected 1,1", tok1.Line, tok1.Column)
	}

	tok2 := s.NextToken() // '\n'
	if tok2.Line != 1 || tok2.Column != 2 {
		t.Errorf("tok2 wrong position: line %d, col %d. expected 1,2", tok2.Line, tok2.Column)
	}

	tok3 := s.NextToken() // 'bc'
	if tok3.Line != 2 || tok3.Column != 1 {
		t.Errorf("tok3 wrong position: line %d, col %d. expected 2,1", tok3.Line, tok3.Column)
	}
}
