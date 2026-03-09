/*
 * api.go
 * This file is part of the gekka-config project.
 *
 * Copyright (c) 2026 Sopranoworks, Osamu Takahashi
 * SPDX-License-Identifier: MIT
 */
package hocon

// ParseString is a convenience function that tokenizes, parses, and wraps
// a HOCON input string into a Config object.
func ParseString(input string) (*Config, error) {
	scanner := NewScanner(input)
	parser := NewParser(scanner)
	obj, err := parser.Parse()
	if err != nil {
		return nil, err
	}
	conf := NewConfig(obj)
	return &conf, nil
}
