/*
 * parser_test.go
 * This file is part of the gekka-config project.
 *
 * Copyright (c) 2026 Sopranoworks, Osamu Takahashi
 * SPDX-License-Identifier: MIT
 */
package hocon

import (
	"testing"
)

func TestParser_Basic(t *testing.T) {
	input := `
		a : 1
		b = "hello"
		c {
			d : true
		}
		e = [1, 2, 3]
	`
	scanner := NewScanner(input)
	parser := NewParser(scanner)
	obj, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	tests := []struct {
		path  string
		want  string
		vtype ValueType
	}{
		{"a", "1", LiteralType},
		{"b", "hello", LiteralType},
		{"c", "{d: true}", ObjectType},
		{"e", "[1, 2, 3]", ListType},
	}

	for _, tt := range tests {
		val, ok := obj.Fields[tt.path]
		if !ok {
			t.Errorf("Field %s not found", tt.path)
			continue
		}
		if val.Type() != tt.vtype {
			t.Errorf("Field %s has wrong type: got %v, want %v", tt.path, val.Type(), tt.vtype)
		}
		if val.String() != tt.want {
			t.Errorf("Field %s has wrong value: got %s, want %s", tt.path, val.String(), tt.want)
		}
	}
}

func TestParser_NestedPath(t *testing.T) {
	input := `
		a.b.c = 42
		x.y = { z: 10 }
	`
	scanner := NewScanner(input)
	parser := NewParser(scanner)
	obj, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Verify a.b.c
	a, ok := obj.Fields["a"].(*Object)
	if !ok {
		t.Fatalf("a is not an object: %v", obj.Fields["a"])
	}
	b, ok := a.Fields["b"].(*Object)
	if !ok {
		t.Fatalf("b is not an object")
	}
	c, ok := b.Fields["c"].(*Literal)
	if !ok || c.Value != 42 {
		t.Errorf("a.b.c is not 42: %v", b.Fields["c"])
	}

	// Verify x.y.z
	x, _ := obj.Fields["x"].(*Object)
	y, _ := x.Fields["y"].(*Object)
	z, _ := y.Fields["z"].(*Literal)
	if z.Value != 10 {
		t.Errorf("x.y.z is not 10: %v", y.Fields["z"])
	}
}

func TestParser_ObjectMerging(t *testing.T) {
	input := `
		a { x: 1 }
		a { y: 2 }
		b.c = 10
		b.d = 20
	`
	scanner := NewScanner(input)
	parser := NewParser(scanner)
	obj, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Verify a
	a, ok := obj.Fields["a"].(*Object)
	if !ok {
		t.Fatal("a is not an object")
	}
	if len(a.Fields) != 2 {
		t.Errorf("a should have 2 fields, got %d: %v", len(a.Fields), a.Fields)
	}

	// Verify b
	b, ok := obj.Fields["b"].(*Object)
	if !ok {
		t.Fatal("b is not an object")
	}
	if len(b.Fields) != 2 {
		t.Errorf("b should have 2 fields, got %d: %v", len(b.Fields), b.Fields)
	}
}

func TestParser_Substitutions(t *testing.T) {
	input := `
		a = ${path.to.val}
		b = ${?optional.val}
	`
	scanner := NewScanner(input)
	parser := NewParser(scanner)
	obj, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	a := obj.Fields["a"].(*Substitution)
	if a.Path != "path.to.val" || a.Optional != false {
		t.Errorf("Wrong substitution for a: %v", a)
	}

	b := obj.Fields["b"].(*Substitution)
	if b.Path != "optional.val" || b.Optional != true {
		t.Errorf("Wrong substitution for b: %v", b)
	}
}

func TestParser_Include(t *testing.T) {
	input := `
		include "foo.conf"
		a = 1
	`
	scanner := NewScanner(input)
	parser := NewParser(scanner)
	obj, err := parser.Parse()
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if obj.Fields["a"].(*Literal).Value != 1 {
		t.Error("a should be 1")
	}
}
