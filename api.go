/*
 * api.go
 * This file is part of the gekka-config project.
 *
 * Copyright (c) 2026 Sopranoworks, Osamu Takahashi
 * SPDX-License-Identifier: MIT
 */
package config

import "github.com/sopranoworks/gekka-config/internal/hocon"

// ParseString is a convenience function that tokenizes, parses, and wraps
// a HOCON input string into a Config object.
func ParseString(input string) (*Config, error) {
	scanner := hocon.NewScanner(input)
	parser := hocon.NewParser(scanner)
	obj, err := parser.Parse()
	if err != nil {
		return nil, err
	}
	conf := newConfig(obj)
	return &conf, nil
}
