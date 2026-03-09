/*
 * merging.go
 * This file is part of the gekka-config project.
 *
 * Copyright (c) 2026 Sopranoworks, Osamu Takahashi
 * SPDX-License-Identifier: MIT
 */
package hocon

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

	mergedRoot := mergeObjectsRecursive(c.root, fallback.root)
	return NewConfig(mergedRoot)
}

func mergeObjectsRecursive(primary, secondary *Object) *Object {
	// Create a new object for immutability
	merged := NewObject()

	// Copy secondary first
	for k, v := range secondary.Fields {
		merged.Fields[k] = v
	}

	// Overwrite with primary
	for k, v := range primary.Fields {
		if primaryVal, isObj := v.(*Object); isObj {
			if secondaryVal, wasObj := merged.Fields[k].(*Object); wasObj {
				// Both are objects, merge recursively
				merged.Fields[k] = mergeObjectsRecursive(primaryVal, secondaryVal)
				continue
			}
		}
		// Not both objects, or primary is literal/list: primary wins
		merged.Fields[k] = v
	}

	return merged
}
