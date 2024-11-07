package main

import (
	"log"

	"github.com/lakshay88/reward-management-system/config"
	"github.com/lakshay88/reward-management-system/database"
	"github.com/lakshay88/reward-management-system/gateway"
)

// Config Variable
var (
	db  database.Database
	cfg *config.AppConfig
)

func init() {
	// Initializing Config
	log.Println("Setting started, fetching configurations")
	var err error
	cfg, err = config.LoadConfiguration("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
		return
	}

	switch cfg.Database.Driver {
	case "postgres":
		db, err = database.ConnectionToPostgres(cfg.Database)
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
	}
}

func main() {

	defer db.Close()
	// Starting Gateway service
	// Instance of Gateway
	gatewayInstance := gateway.NewGateway()
	err := gatewayInstance.RegisterGateWayService(cfg, db)
	if err != nil {
		log.Fatalln("Wait to register routes -", err)
		return
	}
}
