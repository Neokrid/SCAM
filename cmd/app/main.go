package main

import (
	"log"
	"scam/config"
	"scam/internal/app"
)

// @title           SCAM API
// @version         1.0
// @description     This is SCAM api service.

const configDir = "./config/main.yaml"

func main() {
	cfg, err := config.NewConfig(configDir)

	if err != nil {
		log.Fatalf("Config error: %s", err)
	}
	// Run
	app.Run(cfg)
}
