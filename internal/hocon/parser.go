/*
 * parser.go
 * This file is part of the gekka-config project.
 *
 * Copyright (c) 2026 Sopranoworks, Osamu Takahashi
 * SPDX-License-Identifier: MIT
 */
package hocon

import (
	"fmt"
	"strconv"
	"strings"
)

// Parser transforms a token stream into an AST.
type Parser struct {
	scanner       *Scanner
	curTok        Token
	peekTok       Token
	lastTokenLine int
}

// NewParser creates a new Parser for the given Scanner.
func NewParser(s *Scanner) *Parser {
	p := &Parser{scanner: s}
	// Read two tokens, so curTok and peekTok are both set
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.lastTokenLine = p.curTok.Line
	p.curTok = p.peekTok
	p.peekTok = p.scanner.NextToken()
	// Skip comments
	for p.peekTok.Type == TOKEN_COMMENT {
		p.peekTok = p.scanner.NextToken()
	}
}

func (p *Parser) curTokenIs(t TokenType) bool {
	return p.curTok.Type == t
}

func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekTok.Type == t
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	return false
}

// Parse consumes tokens and returns the root Object of the HOCON tree.
func (p *Parser) Parse() (*Object, error) {
	obj := NewObject()
	for !p.curTokenIs(TOKEN_EOF) {
		err := p.parseObjectField(obj)
		if err != nil {
			return nil, err
		}
		p.skipSeparators()
	}
	return obj, nil
}

func (p *Parser) skipSeparators() {
	for p.curTokenIs(TOKEN_COMMA) || p.curTokenIs(TOKEN_NEWLINE) {
		p.nextToken()
	}
}

func (p *Parser) parseObjectField(obj *Object) error {
	p.skipSeparators()
	if p.curTokenIs(TOKEN_EOF) || p.curTokenIs(TOKEN_RIGHT_BRACE) {
		return nil
	}

	// Handle include
	if p.curTokenIs(TOKEN_STRING) && p.curTok.Literal == "include" {
		return p.parseInclude()
	}

	// Parse path: k1.k2.k3
	path, err := p.parsePath()
	if err != nil {
		return err
	}

	// Figure out the separator (:, =, or {)
	p.skipWhitespaceTokens()

	hasSeparator := false
	if p.curTokenIs(TOKEN_COLON) || p.curTokenIs(TOKEN_EQUALS) {
		hasSeparator = true
		p.nextToken()
	}

	var val Value
	if !hasSeparator && p.curTokenIs(TOKEN_LEFT_BRACE) {
		// a { ... } syntax
		val, err = p.parseObject()
	} else {
		val, err = p.parseValue()
	}

	if err != nil {
		return err
	}

	p.setValueAtPath(obj, path, val)
	return nil
}

func (p *Parser) parsePath() ([]string, error) {
	var path []string
	for {
		if !p.curTokenIs(TOKEN_STRING) {
			return nil, fmt.Errorf("expected string for path, got %s at line %d", p.curTok.Type, p.curTok.Line)
		}

		// In HOCON, unquoted strings can contains dots.
		// For simplicity in this phase, we'll split by dot if it's an unquoted string,
		// but wait, the lexer might have already tokens or might not.
		// Actually, HOCON path expressions are a bit complex.
		// Let's assume for now that if it's a string, we check if it contains dots.

		parts := strings.Split(p.curTok.Literal, ".")
		path = append(path, parts...)

		p.nextToken()

		// HOCON also allows "a" . "b"
		// If the next token is a dot (illegal in my current lexer? no, lexer handles it as ILLEGAL or part of string?)
		// Wait, my lexer says isUnquotedStringChar has dots?
		// Let me check lexer.go

		if p.curTokenIs(TOKEN_STRING) && strings.HasPrefix(p.curTok.Literal, ".") {
			// This is tricky. Let's simplify:
			// If curTok is STRING, we take it.
			// We don't have a TOKEN_DOT yet.
			// I might need to update the lexer if I want to handle a.b properly if they are separate tokens.
			// But usually they are one unquoted string if not quoted.
		} else {
			break
		}
	}
	return path, nil
}

func (p *Parser) skipWhitespaceTokens() {
	for p.curTokenIs(TOKEN_NEWLINE) {
		p.nextToken()
	}
}

func (p *Parser) parseValue() (Value, error) {
	val, err := p.parseSingleValue()
	if err != nil {
		return nil, err
	}

	// Check for value concatenation
	var parts []Value
	for p.isConcatenable() {
		nextVal, err := p.parseSingleValue()
		if err != nil {
			return nil, err
		}
		if parts == nil {
			parts = append(parts, val)
		}
		parts = append(parts, nextVal)
	}

	if parts != nil {
		return &Concatenation{Parts: parts}, nil
	}

	return val, nil
}

func (p *Parser) isConcatenable() bool {
	t := p.curTok.Type
	if t == TOKEN_EOF || t == TOKEN_COMMA || t == TOKEN_EQUALS || t == TOKEN_COLON || t == TOKEN_RIGHT_BRACE || t == TOKEN_RIGHT_BRACKET || t == TOKEN_NEWLINE {
		return false
	}
	// Concatenation only happens on the same line as the previous part
	return p.curTok.Line == p.lastTokenLine
}

func (p *Parser) parseSingleValue() (Value, error) {
	p.skipWhitespaceTokens()

	switch p.curTok.Type {
	case TOKEN_STRING:
		lit := p.curTok.Literal
		p.nextToken()
		// Try to parse as bool or null, else string
		if lit == "true" {
			return &Literal{Value: true}, nil
		}
		if lit == "false" {
			return &Literal{Value: false}, nil
		}
		if lit == "null" {
			return &Literal{Value: nil}, nil
		}
		// Try number
		if i, err := strconv.Atoi(lit); err == nil {
			return &Literal{Value: i}, nil
		}
		if f, err := strconv.ParseFloat(lit, 64); err == nil {
			return &Literal{Value: f}, nil
		}
		return &Literal{Value: lit}, nil
	case TOKEN_SUBSTITUTION:
		return p.parseSubstitution()
	case TOKEN_LEFT_BRACE:
		return p.parseObject()
	case TOKEN_LEFT_BRACKET:
		return p.parseList()
	default:
		return nil, fmt.Errorf("unexpected token %s at line %d", p.curTok.Type, p.curTok.Line)
	}
}

func (p *Parser) parseObject() (*Object, error) {
	obj := NewObject()
	if p.curTokenIs(TOKEN_LEFT_BRACE) {
		p.nextToken()
	}

	for !p.curTokenIs(TOKEN_RIGHT_BRACE) && !p.curTokenIs(TOKEN_EOF) {
		err := p.parseObjectField(obj)
		if err != nil {
			return nil, err
		}
		p.skipSeparators()
	}

	if p.curTokenIs(TOKEN_RIGHT_BRACE) {
		p.nextToken()
	}
	return obj, nil
}

func (p *Parser) parseList() (*List, error) {
	list := &List{}
	p.nextToken() // skip [

	for !p.curTokenIs(TOKEN_RIGHT_BRACKET) && !p.curTokenIs(TOKEN_EOF) {
		p.skipSeparators()
		if p.curTokenIs(TOKEN_RIGHT_BRACKET) {
			break
		}
		val, err := p.parseValue()
		if err != nil {
			return nil, err
		}
		list.Elements = append(list.Elements, val)
		p.skipSeparators()
	}

	if p.curTokenIs(TOKEN_RIGHT_BRACKET) {
		p.nextToken()
	}
	return list, nil
}

func (p *Parser) parseSubstitution() (*Substitution, error) {
	lit := p.curTok.Literal // "${path}" or "${?path}"
	p.nextToken()

	optional := false
	path := lit[2 : len(lit)-1] // strip ${ and }
	if strings.HasPrefix(path, "?") {
		optional = true
		path = path[1:]
	}

	return &Substitution{Path: path, Optional: optional}, nil
}

func (p *Parser) parseInclude() error {
	p.nextToken() // skip "include"
	p.skipWhitespaceTokens()
	if !p.curTokenIs(TOKEN_STRING) {
		return fmt.Errorf("expected string after include at line %d", p.curTok.Line)
	}
	// For now, we just skip it as per specs ("focus on syntax parsing, leaving actual file loading for later")
	p.nextToken()
	return nil
}

func (p *Parser) setValueAtPath(root *Object, path []string, val Value) {
	cur := root
	for i := 0; i < len(path)-1; i++ {
		key := path[i]
		if next, ok := cur.Fields[key]; ok {
			if nextObj, ok := next.(*Object); ok {
				cur = nextObj
			} else {
				// Overwrite non-object with object
				newObj := NewObject()
				cur.Fields[key] = newObj
				cur = newObj
			}
		} else {
			newObj := NewObject()
			cur.Fields[key] = newObj
			cur = newObj
		}
	}

	lastKey := path[len(path)-1]
	p.mergeValue(cur, lastKey, val)
}

func (p *Parser) mergeValue(obj *Object, key string, val Value) {
	existing, ok := obj.Fields[key]
	if !ok {
		obj.Fields[key] = val
		return
	}

	existingObj, isExistingObj := existing.(*Object)
	newObj, isNewObj := val.(*Object)

	if isExistingObj && isNewObj {
		// Recursive merge
		for k, v := range newObj.Fields {
			p.mergeValue(existingObj, k, v)
		}
		return
	}

	// Default: overwrite
	obj.Fields[key] = val
}
