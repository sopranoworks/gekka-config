/*
 * ast.go
 * This file is part of the gekka-config project.
 *
 * Copyright (c) 2026 Sopranoworks, Osamu Takahashi
 * SPDX-License-Identifier: MIT
 */
package hocon

import (
	"fmt"
	"strings"
)

// ValueType represents the type of a HOCON value.
type ValueType int

const (
	// ObjectType represents a HOCON object.
	ObjectType ValueType = iota
	// ListType represents a HOCON list.
	ListType
	// LiteralType represents a HOCON primitive value.
	LiteralType
	// SubstitutionType represents a HOCON substitution.
	SubstitutionType
	// ConcatenationType represents multiple values joined together.
	ConcatenationType
)

// Value is the interface for all AST nodes.
type Value interface {
	Type() ValueType
	String() string
}

// Object represents a HOCON object (map of keys to values).
type Object struct {
	// Fields maps keys to HOCON values.
	Fields map[string]Value
}

// Type returns ObjectType.
func (o *Object) Type() ValueType { return ObjectType }

// String returns a string representation of the object.
func (o *Object) String() string {
	var sb strings.Builder
	sb.WriteString("{")
	first := true
	for k, v := range o.Fields {
		if !first {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%s: %s", k, v.String()))
		first = false
	}
	sb.WriteString("}")
	return sb.String()
}

// NewObject creates a new Object.
func NewObject() *Object {
	return &Object{Fields: make(map[string]Value)}
}

// List represents a HOCON list.
type List struct {
	// Elements is the slice of values in the list.
	Elements []Value
}

// Type returns ListType.
func (l *List) Type() ValueType { return ListType }

// String returns a string representation of the list.
func (l *List) String() string {
	var sb strings.Builder
	sb.WriteString("[")
	for i, v := range l.Elements {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(v.String())
	}
	sb.WriteString("]")
	return sb.String()
}

// Literal represents primitive values: strings, numbers, booleans, or null.
type Literal struct {
	// Value holds the actual Go value (string, int, float64, bool, or nil).
	Value interface{}
}

// Type returns LiteralType.
func (l *Literal) Type() ValueType { return LiteralType }

// String returns a string representation of the literal.
func (l *Literal) String() string {
	if l.Value == nil {
		return "null"
	}
	return fmt.Sprintf("%v", l.Value)
}

// Substitution represents a ${path} or ${?path} expression.
type Substitution struct {
	// Path is the dot-notation path to resolve.
	Path string
	// Optional indicates if this is a ${?path} substitution.
	Optional bool
}

// Type returns SubstitutionType.
func (s *Substitution) Type() ValueType { return SubstitutionType }

// String returns a string representation of the substitution.
func (s *Substitution) String() string {
	prefix := "$"
	if s.Optional {
		prefix = "$?"
	}
	return fmt.Sprintf("%s{%s}", prefix, s.Path)
}

// Concatenation represents multiple values joined together (e.g., "${host}:${port}").
type Concatenation struct {
	// Parts is the list of values to be concatenated.
	Parts []Value
}

// Type returns ConcatenationType.
func (c *Concatenation) Type() ValueType { return ConcatenationType }

// String returns a string representation of the concatenation.
func (c *Concatenation) String() string {
	var sb strings.Builder
	for _, p := range c.Parts {
		sb.WriteString(p.String())
	}
	return sb.String()
}
