package main

import (
	"fmt"
	"time"

	"github.com/sopranoworks/gekka-config"
)

type DatabaseConfig struct {
	Host     string `hocon:"host"`
	Port     int    `hocon:"port"`
	PoolSize int    `hocon:"pool-size"`
}

type ServerConfig struct {
	Address string        `hocon:"address"`
	Timeout time.Duration `hocon:"timeout"`
}

type AppConfig struct {
	Name     string         `hocon:"name"`
	Database DatabaseConfig `hocon:"database"`
	Server   ServerConfig   `hocon:"server"`
}

func main() {
	input := `
		app {
			name = "GekkaProduction"
			database {
				host = "db.example.com"
				port = 5432
				pool-size = 20
			}
			server {
				address = "0.0.0.0:8080"
				timeout = 30s
			}
		}
	`

	conf, err := config.ParseString(input)
	if err != nil {
		panic(err)
	}

	resolved, err := conf.Resolve()
	if err != nil {
		panic(err)
	}

	var cfg AppConfig
	// Unmarshal starting from the "app" path
	appConfig, err := resolved.GetConfig("app")
	if err != nil {
		panic(err)
	}
	err = appConfig.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}

	fmt.Printf("App Name: %s\n", cfg.Name)
	fmt.Printf("DB Host: %s, Port: %d\n", cfg.Database.Host, cfg.Database.Port)
	fmt.Printf("Server Timeout: %v\n", cfg.Server.Timeout)
}
