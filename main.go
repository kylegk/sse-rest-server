package main

import (
	"fmt"
	"github.com/kylegk/sse-rest-server/app"
	"github.com/kylegk/sse-rest-server/config"
	"os"
)

func main() {
	conf := setupEnv()
	app.Init(conf)
}

// Set the configuration from environment variables
func setupEnv() config.Config {
	requiredEnvVars := []string {
		config.EnvURL,
		config.EnvPort,
	}

	for _, key := range requiredEnvVars {
		_, exists := os.LookupEnv(key)
		if !exists {
			fmt.Printf("Environment mismatch. Missing required: %s\n", key)
			os.Exit(1)
		}
	}

	clientURL := os.Getenv(config.EnvURL)
	port := os.Getenv(config.EnvPort)

	return config.Config{MemDBSchema: config.DBSchema, SSEServerUrl: clientURL, PORT: port}
}
