/*
 * resolver.go
 * This file is part of the gekka-config project.
 *
 * Copyright (c) 2026 Sopranoworks, Osamu Takahashi
 * SPDX-License-Identifier: MIT
 */
package hocon

import (
	"fmt"
	"strings"
)

type resolver struct {
	root        *Object
	activePaths map[string]bool
}

func newResolver(root *Object) *resolver {
	return &resolver{
		root:        root,
		activePaths: make(map[string]bool),
	}
}

// Resolve is the entry point for resolving substitutions in an Object.
func Resolve(root *Object) (*Object, error) {
	r := newResolver(root)
	return r.resolve()
}

func (r *resolver) resolve() (*Object, error) {
	resolved, err := r.resolveValue(r.root)
	if err != nil {
		return nil, err
	}
	return resolved.(*Object), nil
}

func (r *resolver) resolveValue(v Value) (Value, error) {
	switch val := v.(type) {
	case *Object:
		newObj := NewObject()
		for k, v := range val.Fields {
			resolvedV, err := r.resolveValue(v)
			if err != nil {
				return nil, err
			}
			if resolvedV != nil {
				newObj.Fields[k] = resolvedV
			}
		}
		return newObj, nil
	case *List:
		newList := &List{}
		for _, v := range val.Elements {
			resolvedV, err := r.resolveValue(v)
			if err != nil {
				return nil, err
			}
			if resolvedV != nil {
				newList.Elements = append(newList.Elements, resolvedV)
			}
		}
		return newList, nil
	case *Substitution:
		return r.resolveSubstitution(val)
	case *Concatenation:
		return r.resolveConcatenation(val)
	case *Literal:
		return val, nil
	default:
		return v, nil
	}
}

func (r *resolver) resolveConcatenation(c *Concatenation) (Value, error) {
	var sb strings.Builder
	allObjects := true
	var objects []*Object

	for _, part := range c.Parts {
		resolved, err := r.resolveValue(part)
		if err != nil {
			return nil, err
		}
		if resolved == nil {
			continue // Skip missing optional substitutions
		}
		if resolved.Type() != ObjectType {
			allObjects = false
		} else {
			objects = append(objects, resolved.(*Object))
		}
		sb.WriteString(resolved.String())
	}

	if allObjects && len(objects) > 0 {
		// HOCON: Objects merge when concatenated
		result := objects[0]
		for i := 1; i < len(objects); i++ {
			result = MergeObjectsRecursive(result, objects[i])
		}
		return result, nil
	}

	// Otherwise, it's a string concatenation
	return &Literal{Value: sb.String()}, nil
}

func (r *resolver) resolveSubstitution(s *Substitution) (Value, error) {
	if r.activePaths[s.Path] {
		return nil, fmt.Errorf("circular dependency detected at path: %s", s.Path)
	}

	r.activePaths[s.Path] = true
	defer delete(r.activePaths, s.Path)

	// 1. Try to look up in the root object
	val, err := r.lookupPath(s.Path)
	if err == nil {
		// Found it, now resolve it (it might be another substitution)
		return r.resolveValue(val)
	}

	// 2. Try environment variables
	if envVal, ok := LookupEnv(s.Path); ok {
		return &Literal{Value: envVal}, nil
	}

	// 3. Optional substitution
	if s.Optional {
		return nil, nil // Return nil so it can be filtered out
	}

	return nil, fmt.Errorf("could not resolve mandatory substitution: ${%s}", s.Path)
}

func (r *resolver) lookupPath(path string) (Value, error) {
	if path == "" {
		return r.root, nil
	}

	parts := strings.Split(path, ".")
	var current Value = r.root

	for _, part := range parts {
		obj, ok := current.(*Object)
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
