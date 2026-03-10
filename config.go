/*
 * config.go
 * This file is part of the gekka-config project.
 *
 * Copyright (c) 2026 Sopranoworks, Osamu Takahashi
 * SPDX-License-Identifier: MIT
 */
// Package hocon provides a pure Go implementation of the HOCON (Human-Optimized Config Object Notation) format.
// It is designed for zero dependencies and high compatibility with Lightbend's Java implementation.
package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/sopranoworks/gekka-config/internal/hocon"
)

// Config wraps the root Object and provides a high-level API for configuration traversal and retrieval.
type Config struct {
	root *hocon.Object
}

// newConfig creates a new Config instance wrapping the provided root Object.
func newConfig(root *hocon.Object) Config {
	return Config{root: root}
}

// getValue retrieves the raw Value at the given dot-notation path (e.g., "a.b.c").
func (c Config) getValue(path string) (hocon.Value, error) {
	if path == "" {
		return c.root, nil
	}

	parts := strings.Split(path, ".")
	var current hocon.Value = c.root

	for _, part := range parts {
		obj, ok := current.(*hocon.Object)
		if !ok {
			return nil, fmt.Errorf("path '%s' not found: '%s' is not an object", path, part)
		}

		val, ok := obj.Fields[part]
		if !ok {
			return nil, fmt.Errorf("path '%s' not found: key '%s' missing", path, part)
		}
		current = val
	}

	return current, nil
}

// GetConfig returns a new Config instance rooted at the given path.
func (c Config) GetConfig(path string) (Config, error) {
	val, err := c.getValue(path)
	if err != nil {
		return Config{}, err
	}

	obj, ok := val.(*hocon.Object)
	if !ok {
		return Config{}, fmt.Errorf("value at '%s' is not an object", path)
	}

	return newConfig(obj), nil
}

// GetString retrieves the value at path as a string.
func (c Config) GetString(path string) (string, error) {
	val, err := c.getValue(path)
	if err != nil {
		return "", err
	}

	lit, ok := val.(*hocon.Literal)
	if !ok {
		return "", fmt.Errorf("value at '%s' is not a literal", path)
	}

	return fmt.Sprint(lit.Value), nil
}

// GetInt retrieves the value at path as an integer.
func (c Config) GetInt(path string) (int, error) {
	val, err := c.getValue(path)
	if err != nil {
		return 0, err
	}

	lit, ok := val.(*hocon.Literal)
	if !ok {
		return 0, fmt.Errorf("value at '%s' is not a literal", path)
	}

	switch v := lit.Value.(type) {
	case int:
		return v, nil
	case float64:
		return int(v), nil
	case string:
		var i int
		_, err := fmt.Sscan(v, &i)
		return i, err
	default:
		return 0, fmt.Errorf("value at '%s' cannot be converted to int", path)
	}
}

// GetBoolean retrieves the value at path as a boolean.
func (c Config) GetBoolean(path string) (bool, error) {
	val, err := c.getValue(path)
	if err != nil {
		return false, err
	}

	lit, ok := val.(*hocon.Literal)
	if !ok {
		return false, fmt.Errorf("value at '%s' is not a literal", path)
	}

	switch v := lit.Value.(type) {
	case bool:
		return v, nil
	case string:
		return strings.ToLower(v) == "true", nil
	default:
		return false, fmt.Errorf("value at '%s' is not a boolean", path)
	}
}

// GetDuration retrieves the value at path as a time.Duration.
// It supports HOCON duration formats such as "10s", "100ms", "1h".
func (c Config) GetDuration(path string) (time.Duration, error) {
	s, err := c.GetString(path)
	if err != nil {
		return 0, err
	}

	// Simple HOCON duration parsing
	// HOCON spec defines many units, but let's start with standard Go time.ParseDuration
	// but adding support for space like "10 s"
	s = strings.ReplaceAll(s, " ", "")

	// time.ParseDuration supports ns, us, ms, s, m, h
	return time.ParseDuration(s)
}

// Resolve scans the entire AST for substitutions (${path}) and resolves them.
func (c Config) Resolve() (Config, error) {
	if c.root == nil {
		return c, nil
	}
	resolvedRoot, err := hocon.Resolve(c.root)
	if err != nil {
		return Config{}, err
	}
	return newConfig(resolvedRoot), nil
}

// Unmarshal binds the configuration to a struct.
func (c Config) Unmarshal(v interface{}) error {
	return Unmarshal(c, v)
}

// root returns the internal root object.
func (c Config) getRoot() *hocon.Object {
	return c.root
}
