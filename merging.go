/*
 * merging.go
 * This file is part of the gekka-config project.
 *
 * Copyright (c) 2026 Sopranoworks, Osamu Takahashi
 * SPDX-License-Identifier: MIT
 */
package config

import "github.com/sopranoworks/gekka-config/internal/hocon"

// WithFallback returns a new Config where the current config takes precedence,
// and the provided fallback config provides values for missing keys.
// If a key exists as an object in both, they are merged recursively.
func (c Config) WithFallback(fallback Config) Config {
	if c.root == nil {
		return fallback
	}
	if fallback.root == nil {
		return c
	}

	mergedRoot := hocon.MergeObjectsRecursive(c.root, fallback.root)
	return newConfig(mergedRoot)
}
