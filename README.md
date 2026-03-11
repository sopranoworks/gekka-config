# gekka-config 💎

**gekka-config** is a zero-dependency, pure Go implementation of HOCON (Human-Optimized Config Object Notation), designed for the gekka Actor System and Pekko/Akka compatibility.

## Installation

```bash
go get github.com/sopranoworks/gekka-config
```

## Quick Start

The most powerful way to use **gekka-config** is mapping HOCON directly to Go structs.

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

For more advanced usage, check the `/examples` directory:

- **examples/basic/**: Basic key-value retrieval.
- **examples/unmarshal/**: Complex struct mapping with tags and nested objects.
- **examples/merging/**: Layering configurations using `WithFallback`.

## License

MIT License. Copyright (c) 2026 Sopranoworks, Osamu Takahashi.