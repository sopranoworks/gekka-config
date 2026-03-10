/*
 * merging.go
 * This file is part of the gekka-config project.
 *
 * Copyright (c) 2026 Sopranoworks, Osamu Takahashi
 * SPDX-License-Identifier: MIT
 */
package hocon

// MergeObjectsRecursive merges two objects recursively.
// The primary object's values take precedence over the secondary object's.
func MergeObjectsRecursive(primary, secondary *Object) *Object {
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
				merged.Fields[k] = MergeObjectsRecursive(primaryVal, secondaryVal)
				continue
			}
		}
		// Not both objects, or primary is literal/list: primary wins
		merged.Fields[k] = v
	}

	return merged
}
