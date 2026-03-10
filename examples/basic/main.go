/*
 * main.go
 * This file is part of the gekka-config project.
 *
 * Copyright (c) 2026 Sopranoworks, Osamu Takahashi
 * SPDX-License-Identifier: MIT
 */
package main

import (
	"fmt"

	config "github.com/sopranoworks/gekka-config"
)

func main() {
	input := `
        pekko {
            remote.artery.canonical.port = 25520
            cluster.seed-nodes = ["pekko://system@127.0.0.1:25520"]
        }
    `
	conf, err := config.ParseString(input)
	if err != nil {
		fmt.Printf("Error parsing: %v\n", err)
		return
	}
	resolved, err := conf.Resolve()
	if err != nil {
		fmt.Printf("Error resolving: %v\n", err)
		return
	}

	port, _ := resolved.GetInt("pekko.remote.artery.canonical.port")
	fmt.Printf("Pekko Port: %d\n", port)
}
