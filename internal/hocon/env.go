/*
 * env.go
 * This file is part of the gekka-config project.
 *
 * Copyright (c) 2026 Sopranoworks, Osamu Takahashi
 * SPDX-License-Identifier: MIT
 */
package hocon

import (
	"os"
	"strings"
)

// LookupEnv checks for an environment variable matching the HOCON path.
// It converts dot-notation paths like "a.b.c" to "A_B_C".
func LookupEnv(path string) (string, bool) {
	envName := strings.ToUpper(strings.ReplaceAll(path, ".", "_"))
	val, ok := os.LookupEnv(envName)
	return val, ok
}
