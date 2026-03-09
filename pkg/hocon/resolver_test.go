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
	conf := NewConfig(obj)

	resolved, err := conf.Resolve()
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}

	if val, _ := resolved.GetInt("b"); val != 1 {
		t.Errorf("b = %v, want 1", val)
	}

	if val, _ := resolved.GetInt("c.d"); val != 1 {
		t.Errorf("c.d = %v, want 1", val)
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
	conf := NewConfig(obj)

	resolved, err := conf.Resolve()
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}

	// a should be missing
	if _, err := resolved.GetValue("a"); err == nil {
		t.Error("a should be missing from resolved config")
	}

	// c should be 42
	if val, _ := resolved.GetInt("c"); val != 42 {
		t.Errorf("c = %v, want 42", val)
	}
}

func TestResolver_EnvVar(t *testing.T) {
	os.Setenv("HOCON_TEST_VAR", "env-value")
	defer os.Unsetenv("HOCON_TEST_VAR")

	input := `a : ${hocon.test.var}`
	scanner := NewScanner(input)
	parser := NewParser(scanner)
	obj, _ := parser.Parse()
	conf := NewConfig(obj)

	resolved, err := conf.Resolve()
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}

	if val, _ := resolved.GetString("a"); val != "env-value" {
		t.Errorf("a = %v, want env-value", val)
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
	conf := NewConfig(obj)

	_, err := conf.Resolve()
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
	conf := NewConfig(obj)

	resolved, err := conf.Resolve()
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}

	if val, _ := resolved.GetInt("c"); val != 1 {
		t.Errorf("c = %v, want 1", val)
	}
}
