/*
 * token.go
 * This file is part of the gekka-config project.
 *
 * Copyright (c) 2026 Sopranoworks, Osamu Takahashi
 * SPDX-License-Identifier: MIT
 */
package hocon

// TokenType represents the type of a HOCON token.
type TokenType string

const (
	// TOKEN_ILLEGAL represents an unknown character.
	TOKEN_ILLEGAL TokenType = "ILLEGAL"
	// TOKEN_EOF represents the end of the input.
	TOKEN_EOF TokenType = "EOF"
	// TOKEN_LEFT_BRACE represents '{'.
	TOKEN_LEFT_BRACE TokenType = "{"
	// TOKEN_RIGHT_BRACE represents '}'.
	TOKEN_RIGHT_BRACE TokenType = "}"
	// TOKEN_LEFT_BRACKET represents '['.
	TOKEN_LEFT_BRACKET TokenType = "["
	// TOKEN_RIGHT_BRACKET represents ']'.
	TOKEN_RIGHT_BRACKET TokenType = "]"
	// TOKEN_COLON represents ':'.
	TOKEN_COLON TokenType = ":"
	// TOKEN_EQUALS represents '='.
	TOKEN_EQUALS TokenType = "="
	// TOKEN_COMMA represents ','.
	TOKEN_COMMA TokenType = ","
	// TOKEN_SUBSTITUTION represents a HOCON substitution like ${path}.
	TOKEN_SUBSTITUTION TokenType = "SUBSTITUTION"
	// TOKEN_STRING represents a quoted or unquoted string.
	TOKEN_STRING TokenType = "STRING"
	// TOKEN_COMMENT represents a comment starting with # or //.
	TOKEN_COMMENT TokenType = "COMMENT"
	// TOKEN_NEWLINE represents a newline character.
	TOKEN_NEWLINE TokenType = "NEWLINE"
)

// Token represents a single HOCON lexical token.
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}
