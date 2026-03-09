# gekka-config

![Version](https://img.shields.io/badge/version-1.0.0-blue)
![License](https://img.shields.io/badge/license-MIT-green)
![CI Status](https://github.com/takahashi/gekka-config/actions/workflows/go.yml/badge.svg?branch=master)

Pure Go HOCON (Human-Optimized Config Object Notation) implementation for Pekko/Akka compatibility.

## Features

- Standard HOCON Support: Handles nested objects, lists, and primitives.
- Substitution Engine: Full support for ${path} and ${?path}.
- Config Merging: Recursive merging of configurations via layered fallbacks.
- Struct Binding: Map HOCON directly to Go structs with native type support.
- Zero Dependencies: Optimized for clean, dependency-free integration.

## Installation

To install the library, run:

go get github.com/takahashi/gekka-config

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/takahashi/gekka-config/pkg/hocon"
)

func main() {
    input := `
        app {
            name = "GekkaApp"
            timeout = 5s
            port = ${?PORT}
        }
    `
    // Parse, Resolve, and Access
    conf, _ := hocon.ParseString(input)
    resolved, _ := conf.Resolve()
    
    name, _ := resolved.GetString("app.name")
    fmt.Printf("App: %s\n", name)
}
``` 

## Features Deep Dive

1. Parsing: Robust recursive descent parser with coordinate tracking.
2. Merging: Seamlessly layer environmental, local, and default configs.
3. Substitution: Built-in environment variable fallback and circularity detection.
4. Struct Binding: High-performance reflection-based unmarshalling with struct tags.

## License

MIT License.
