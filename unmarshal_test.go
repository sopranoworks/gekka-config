/*
 * unmarshal_test.go
 * This file is part of the gekka-config project.
 *
 * Copyright (c) 2026 Sopranoworks, Osamu Takahashi
 * SPDX-License-Identifier: MIT
 */
package config

import (
	"testing"
	"time"

	"github.com/sopranoworks/gekka-config/internal/hocon"
)

func TestUnmarshal_Basic(t *testing.T) {
	input := `
		name : "Aoi"
		version : 1
		enabled : true
		timeout : 5s
	`
	scanner := hocon.NewScanner(input)
	parser := hocon.NewParser(scanner)
	obj, _ := parser.Parse()
	conf := newConfig(obj)

	type ConfigStruct struct {
		Name    string
		Version int
		Enabled bool
		Timeout time.Duration
	}

	var cfg ConfigStruct
	err := conf.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if cfg.Name != "Aoi" {
		t.Errorf("Name = %v, want Aoi", cfg.Name)
	}
	if cfg.Version != 1 {
		t.Errorf("Version = %v, want 1", cfg.Version)
	}
	if cfg.Enabled != true {
		t.Errorf("Enabled = %v, want true", cfg.Enabled)
	}
	if cfg.Timeout != 5*time.Second {
		t.Errorf("Timeout = %v, want 5s", cfg.Timeout)
	}
}

func TestUnmarshal_Tags(t *testing.T) {
	input := `
		nested {
			path : "tagged-value"
		}
	`
	scanner := hocon.NewScanner(input)
	parser := hocon.NewParser(scanner)
	obj, _ := parser.Parse()
	conf := newConfig(obj)

	type ConfigStruct struct {
		Mapped string `hocon:"nested.path"`
	}

	var cfg ConfigStruct
	err := conf.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if cfg.Mapped != "tagged-value" {
		t.Errorf("Mapped = %v, want tagged-value", cfg.Mapped)
	}
}

func TestUnmarshal_Nested(t *testing.T) {
	input := `
		server {
			host : "localhost"
			port : 8080
		}
	`
	scanner := hocon.NewScanner(input)
	parser := hocon.NewParser(scanner)
	obj, _ := parser.Parse()
	conf := newConfig(obj)

	type ServerConfig struct {
		Host string
		Port int
	}
	type AppConfig struct {
		Server ServerConfig
	}

	var cfg AppConfig
	err := conf.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if cfg.Server.Host != "localhost" {
		t.Errorf("Host = %v, want localhost", cfg.Server.Host)
	}
	if cfg.Server.Port != 8080 {
		t.Errorf("Port = %v, want 8080", cfg.Server.Port)
	}
}

func TestUnmarshal_Slice(t *testing.T) {
	input := `
		items : ["a", "b", "c"]
		nums : [1, 2, 3]
	`
	scanner := hocon.NewScanner(input)
	parser := hocon.NewParser(scanner)
	obj, _ := parser.Parse()
	conf := newConfig(obj)

	type ConfigStruct struct {
		Items []string
		Nums  []int
	}

	var cfg ConfigStruct
	err := conf.Unmarshal(&cfg)
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	if len(cfg.Items) != 3 || cfg.Items[1] != "b" {
		t.Errorf("Items = %v, want [a b c]", cfg.Items)
	}
	if len(cfg.Nums) != 3 || cfg.Nums[2] != 3 {
		t.Errorf("Nums = %v, want [1 2 3]", cfg.Nums)
	}
}
