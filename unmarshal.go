/*
 * unmarshal.go
 * This file is part of the gekka-config project.
 *
 * Copyright (c) 2026 Sopranoworks, Osamu Takahashi
 * SPDX-License-Identifier: MIT
 */
package config

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/sopranoworks/gekka-config/internal/hocon"
)

// Unmarshal binds the configuration values to the fields of the provided struct pointer.
// It supports `hocon` struct tags for explicit path mapping.
func Unmarshal(c Config, v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("unmarshal requires a pointer to a struct")
	}

	return unmarshalObject(c, val.Elem())
}

func unmarshalObject(c Config, structVal reflect.Value) error {
	structType := structVal.Type()

	for i := 0; i < structVal.NumField(); i++ {
		field := structType.Field(i)
		fieldVal := structVal.Field(i)

		if !fieldVal.CanSet() {
			continue
		}

		path := getHoconPath(field)
		if path == "" {
			continue
		}

		err := unmarshalField(c, path, fieldVal)
		if err != nil {
			return fmt.Errorf("failed to unmarshal field %s: %w", field.Name, err)
		}
	}

	return nil
}

func getHoconPath(field reflect.StructField) string {
	tag := field.Tag.Get("hocon")
	if tag != "" {
		return tag
	}

	// Fallback to case-insensitive field name match
	// For simplicity, we'll try lowercase version of the field name
	return strings.ToLower(field.Name)
}

func unmarshalField(c Config, path string, fieldVal reflect.Value) error {
	// Specialized handling for time.Duration
	if fieldVal.Type() == reflect.TypeOf(time.Duration(0)) {
		d, err := c.GetDuration(path)
		if err != nil {
			return nil // Skip if not found, or should we return error?
		}
		fieldVal.Set(reflect.ValueOf(d))
		return nil
	}

	switch fieldVal.Kind() {
	case reflect.String:
		s, err := c.GetString(path)
		if err == nil {
			fieldVal.SetString(s)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := c.GetInt(path)
		if err == nil {
			fieldVal.SetInt(int64(i))
		}
	case reflect.Bool:
		b, err := c.GetBoolean(path)
		if err == nil {
			fieldVal.SetBool(b)
		}
	case reflect.Float32, reflect.Float64:
		// We don't have GetFloat, but we can use getValue and check Literal
		val, err := c.getValue(path)
		if err == nil {
			if lit, ok := val.(*hocon.Literal); ok {
				switch f := lit.Value.(type) {
				case float64:
					fieldVal.SetFloat(f)
				case int:
					fieldVal.SetFloat(float64(f))
				}
			}
		}
	case reflect.Struct:
		subConfig, err := c.GetConfig(path)
		if err == nil {
			return unmarshalObject(subConfig, fieldVal)
		}
	case reflect.Ptr:
		if fieldVal.Type().Elem().Kind() == reflect.Struct {
			subConfig, err := c.GetConfig(path)
			if err == nil {
				if fieldVal.IsNil() {
					fieldVal.Set(reflect.New(fieldVal.Type().Elem()))
				}
				return unmarshalObject(subConfig, fieldVal.Elem())
			}
		}
	case reflect.Slice:
		val, err := c.getValue(path)
		if err == nil {
			if list, ok := val.(*hocon.List); ok {
				slice := reflect.MakeSlice(fieldVal.Type(), len(list.Elements), len(list.Elements))
				for i, elem := range list.Elements {
					err := setReflectValue(slice.Index(i), elem)
					if err != nil {
						return err
					}
				}
				fieldVal.Set(slice)
			}
		}
	}

	return nil
}

func setReflectValue(v reflect.Value, hValue hocon.Value) error {
	switch v.Kind() {
	case reflect.String:
		if lit, ok := hValue.(*hocon.Literal); ok {
			v.SetString(fmt.Sprint(lit.Value))
		}
	case reflect.Int, reflect.Int64:
		if lit, ok := hValue.(*hocon.Literal); ok {
			if i, ok := lit.Value.(int); ok {
				v.SetInt(int64(i))
			}
		}
	case reflect.Bool:
		if lit, ok := hValue.(*hocon.Literal); ok {
			if b, ok := lit.Value.(bool); ok {
				v.SetBool(b)
			}
		}
	}
	return nil
}
