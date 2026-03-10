/*
 * resolver_test.go
 * This file is part of the gekka-config project.
 *
 * Copyright (c) 2026 Sopranoworks, Osamu Takahashi
 * SPDX-License-Identifier: MIT
 */
package hocon

import (
	"os"
	"testing"
)

func TestResolver_Basic(t *testing.T) {
	input := `
		a : 1
		b : ${a}
		c : {
			d : ${a}
		}
	`
	scanner := NewScanner(input)
	parser := NewParser(scanner)
	obj, _ := parser.Parse()

	resolved, err := Resolve(obj)
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}

	// Verify b
	if bVal, ok := resolved.Fields["b"].(*Literal); ok {
		if bVal.Value != 1 {
			t.Errorf("b = %v, want 1", bVal.Value)
		}
	} else {
		t.Errorf("b is not a literal: %T", resolved.Fields["b"])
	}

	// Verify c.d
	if cVal, ok := resolved.Fields["c"].(*Object); ok {
		if dVal, ok := cVal.Fields["d"].(*Literal); ok {
			if dVal.Value != 1 {
				t.Errorf("c.d = %v, want 1", dVal.Value)
			}
		} else {
			t.Errorf("c.d is not a literal: %T", cVal.Fields["d"])
		}
	} else {
		t.Errorf("c is not an object: %T", resolved.Fields["c"])
	}
}

func TestResolver_Optional(t *testing.T) {
	input := `
		a : ${?missing}
		b : 42
		c : ${?b}
	`
	scanner := NewScanner(input)
	parser := NewParser(scanner)
	obj, _ := parser.Parse()

	resolved, err := Resolve(obj)
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}

	// a should be missing
	if _, ok := resolved.Fields["a"]; ok {
		t.Error("a should be missing from resolved object")
	}

	// c should be 42
	if cVal, ok := resolved.Fields["c"].(*Literal); ok {
		if cVal.Value != 42 {
			t.Errorf("c = %v, want 42", cVal.Value)
		}
	} else {
		t.Errorf("c is not a literal: %T", resolved.Fields["c"])
	}
}

func TestResolver_EnvVar(t *testing.T) {
	os.Setenv("HOCON_TEST_VAR", "env-value")
	defer os.Unsetenv("HOCON_TEST_VAR")

	input := `a : ${hocon_test_var}`
	scanner := NewScanner(input)
	parser := NewParser(scanner)
	obj, _ := parser.Parse()

	resolved, err := Resolve(obj)
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}

	if aVal, ok := resolved.Fields["a"].(*Literal); ok {
		if aVal.Value != "env-value" {
			t.Errorf("a = %v, want env-value", aVal.Value)
		}
	} else {
		t.Errorf("a is not a literal: %T", resolved.Fields["a"])
	}
}

func TestResolver_Circular(t *testing.T) {
	input := `
		a : ${b}
		b : ${a}
	`
	scanner := NewScanner(input)
	parser := NewParser(scanner)
	obj, _ := parser.Parse()

	_, err := Resolve(obj)
	if err == nil {
		t.Fatal("expected error for circular dependency")
	}
}

func TestResolver_NestedSubstitutions(t *testing.T) {
	input := `
		a : 1
		b : ${a}
		c : ${b}
	`
	scanner := NewScanner(input)
	parser := NewParser(scanner)
	obj, _ := parser.Parse()

	resolved, err := Resolve(obj)
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}

	if cVal, ok := resolved.Fields["c"].(*Literal); ok {
		if cVal.Value != 1 {
			t.Errorf("c = %v, want 1", cVal.Value)
		}
	} else {
		t.Errorf("c is not a literal: %T", resolved.Fields["c"])
	}
}
