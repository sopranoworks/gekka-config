/*
 * config_test.go
 * This file is part of the gekka-config project.
 *
 * Copyright (c) 2026 Sopranoworks, Osamu Takahashi
 * SPDX-License-Identifier: MIT
 */
package hocon

import (
	"testing"
	"time"
)

func TestConfig_Retrieval(t *testing.T) {
	input := `
		a : "hello"
		b : 42
		c : true
		d : {
			e : "inner"
			f : 100ms
		}
	`
	scanner := NewScanner(input)
	parser := NewParser(scanner)
	obj, _ := parser.Parse()
	conf := NewConfig(obj)

	// String
	if val, _ := conf.GetString("a"); val != "hello" {
		t.Errorf("GetString(a) = %v, want hello", val)
	}

	// Int
	if val, _ := conf.GetInt("b"); val != 42 {
		t.Errorf("GetInt(b) = %v, want 42", val)
	}

	// Boolean
	if val, _ := conf.GetBoolean("c"); val != true {
		t.Errorf("GetBoolean(c) = %v, want true", val)
	}

	// Nested String
	if val, _ := conf.GetString("d.e"); val != "inner" {
		t.Errorf("GetString(d.e) = %v, want inner", val)
	}

	// Duration
	if val, _ := conf.GetDuration("d.f"); val != 100*time.Millisecond {
		t.Errorf("GetDuration(d.f) = %v, want 100ms", val)
	}
}

func TestConfig_WithFallback(t *testing.T) {
	input1 := `
		a : 1
		common : {
			x : "primary"
		}
	`
	input2 := `
		b : 2
		common : {
			y : "fallback"
			x : "should-be-ignored"
		}
	`

	p1 := NewParser(NewScanner(input1))
	o1, _ := p1.Parse()
	c1 := NewConfig(o1)

	p2 := NewParser(NewScanner(input2))
	o2, _ := p2.Parse()
	c2 := NewConfig(o2)

	merged := c1.WithFallback(c2)

	// primary values take precedence
	if val, _ := merged.GetInt("a"); val != 1 {
		t.Errorf("a = %v, want 1", val)
	}
	// fallback values are used for missing keys
	if val, _ := merged.GetInt("b"); val != 2 {
		t.Errorf("b = %v, want 2", val)
	}
	// objects are merged recursively
	if val, _ := merged.GetString("common.x"); val != "primary" {
		t.Errorf("common.x = %v, want primary", val)
	}
	if val, _ := merged.GetString("common.y"); val != "fallback" {
		t.Errorf("common.y = %v, want fallback", val)
	}
}

func TestConfig_ErrorCases(t *testing.T) {
	input := `a : 1`
	root, _ := NewParser(NewScanner(input)).Parse()
	conf := NewConfig(root)

	if _, err := conf.GetString("non-existent"); err == nil {
		t.Error("expected error for non-existent key")
	}

	if _, err := conf.GetConfig("a"); err == nil {
		t.Error("expected error when GetConfig on a literal")
	}
}
