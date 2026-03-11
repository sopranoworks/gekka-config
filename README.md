# gekka-config 💎

![Version](https://img.shields.io/badge/version-1.0.4-blue)
![License](https://img.shields.io/badge/license-MIT-green)
![CI Status](https://github.com/sopranoworks/gekka-config/actions/workflows/go.yml/badge.svg?branch=main)

gekka-config is a zero-dependency, pure Go implementation of HOCON (Human-Optimized Config Object Notation). Designed as the bedrock configuration engine for the gekka Actor System, it provides high compatibility with the Pekko/Akka ecosystem.

## Installation

```bash
go get github.com/sopranoworks/gekka-config
```

## Quick Start

The most powerful way to use gekka-config is mapping HOCON directly to Go structs.

```go
package main

import (
    "fmt"
    "time"
    "github.com/sopranoworks/gekka-config"
)

type AppConfig struct {
    Name    string        `hocon:"app.name"`
    Timeout time.Duration `hocon:"app.timeout"`
}

func main() {
    input := `app { name = "GekkaApp", timeout = 5s }`
    
    conf, err := config.ParseString(input)
    if err != nil {
        panic(err)
    }
    resolved, err := conf.Resolve()
    if err != nil {
        panic(err)
    }

    var cfg AppConfig
    resolved.Unmarshal(&cfg)

    fmt.Printf("App: %s, Timeout: %v\n", cfg.Name, cfg.Timeout)
}
```

## Examples

For more advanced usage, check the [/examples](https://github.com/sopranoworks/gekka-config/tree/main/examples) directory:

- **basic**: Basic key-value retrieval.
- **unmarshal**: Complex struct mapping with tags and nested objects.
- **merging**: Layering configurations using WithFallback.

## License

MIT License. Copyright (c) 2026 Sopranoworks, Osamu Takahashi.