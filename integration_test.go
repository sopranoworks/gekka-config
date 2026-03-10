/*
 * integration_test.go
 * This file is part of the gekka-config project.
 *
 * Copyright (c) 2026 Sopranoworks, Osamu Takahashi
 * SPDX-License-Identifier: MIT
 */
package config

import (
	"testing"
	"time"
)

func TestIntegration_EndToEnd(t *testing.T) {
	hoconStr := `
        base {
            timeout = 5s
            host = "localhost"
        }
        app {
            name = "GekkaNode"
            port = ${?PORT}
            port = 8080
            addr = ${base.host} ":" ${app.port}
            connect-timeout = ${base.timeout}
        }
    `
	conf, err := ParseString(hoconStr)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	resolved, err := conf.Resolve()
	if err != nil {
		t.Fatalf("Failed to resolve: %v", err)
	}

	type Config struct {
		Name    string        `hocon:"app.name"`
		Addr    string        `hocon:"app.addr"`
		Timeout time.Duration `hocon:"app.connect-timeout"`
	}

	var cfg Config
	if err := resolved.Unmarshal(&cfg); err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Wait, if I don't support concatenation, addr will likely be just "localhost" or fail.
	// HOCON value concatenation: ${base.host} ":" ${app.port}
	// My current parser only takes the first value.

	if cfg.Addr != "localhost:8080" {
		t.Errorf("Expected addr localhost:8080, got %q", cfg.Addr)
	}
}

func TestIntegration_EnvVariable(t *testing.T) {
	t.Setenv("APP_PORT", "9090")
	hoconStr := `
        app {
            port = 8080
            port = ${?APP_PORT}
        }
    `
	conf, err := ParseString(hoconStr)
	if err != nil {
		t.Fatalf("Failed to parse: %v", err)
	}

	resolved, err := conf.Resolve()
	if err != nil {
		t.Fatalf("Failed to resolve: %v", err)
	}

	port, err := resolved.GetInt("app.port")
	if err != nil {
		t.Fatalf("Failed to get port: %v", err)
	}

	if port != 9090 {
		t.Errorf("Expected port 9090 from env, got %d", port)
	}
}
