package main

import (
	"fmt"

	"github.com/sopranoworks/gekka-config"
)

func main() {
	// Reference configuration (defaults)
	defaults := `
		app {
			name = "GekkaApp"
			port = 8080
			features {
				logging = true
				monitoring = false
			}
		}
	`

	// Local overrides
	overrides := `
		app {
			port = 9090
			features {
				monitoring = true
			}
		}
	`

	// Parse both (ParseString returns *Config)
	defaultConf, _ := config.ParseString(defaults)
	overridesConf, _ := config.ParseString(overrides)

	// Merge: overrides with defaults as fallback
	// Use dereference since WithFallback takes Config (not *Config)
	merged := overridesConf.WithFallback(*defaultConf)

	// Resolve substitutions (if any)
	resolved, _ := merged.Resolve()

	// Verify values
	name, _ := resolved.GetString("app.name")           // from defaults
	port, _ := resolved.GetInt("app.port")             // from overrides
	logging, _ := resolved.GetBoolean("app.features.logging")    // from defaults
	monitoring, _ := resolved.GetBoolean("app.features.monitoring") // from overrides

	fmt.Printf("App: %s\n", name)
	fmt.Printf("Port: %d\n", port)
	fmt.Printf("Logging: %v\n", logging)
	fmt.Printf("Monitoring: %v\n", monitoring)
}
