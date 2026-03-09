/*
 * error_test.go
 * This file is part of the gekka-config project.
 *
 * Copyright (c) 2026 Sopranoworks, Osamu Takahashi
 * SPDX-License-Identifier: MIT
 */
package hocon

import (
	"testing"
)

func TestError_CircularDependency(t *testing.T) {
	input := `a = ${b}, b = ${a}`
	conf, err := ParseString(input)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}
	_, err = conf.Resolve()
	if err == nil {
		t.Error("Expected error for circular dependency, got nil")
	}
}

func TestError_SyntaxError(t *testing.T) {
	input := `a { b = }`
	_, err := ParseString(input)
	if err == nil {
		t.Error("Expected syntax error, got nil")
	}
}
