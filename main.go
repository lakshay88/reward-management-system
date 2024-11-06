package main

import (
	"log"

	"github.com/lakshay88/reward-managment-system/config"
	"github.com/lakshay88/reward-managment-system/gateway"
)

// Config Variable
var cfg *config.AppConfig

func init() {
	// Initializing Config
	log.Println("Setting started, fetching configurations")
	var err error
	cfg, err = config.LoadConfiguration("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
		return
	}
}

func main() {
	// Starting Gateway service
	// Instance of Gateway
	gate := gateway.NewGateway()
	err := gate.RegisterGateWayService(cfg.ServerConfig)
	if err != nil {
		log.Fatalln("Wait to register routes -", err)
		return
	}
}
