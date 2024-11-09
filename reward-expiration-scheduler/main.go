package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lakshay88/reward-management-system/config"
	"github.com/lakshay88/reward-management-system/database"
)

var (
	db  database.Database
	cfg *config.AppConfig
)

func init() {
	// Initializing Config
	log.Println("Setting started, fetching configurations")
	var err error
	cfg, err = config.LoadConfiguration("../config.yaml")
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
	log.Println("fetching configurations- Completed")
}

func main() {

	defer db.Close()

	timePeroid := time.Duration(cfg.SchedulerConfig.SchedulerRunnerTimeInMin) * time.Minute
	ticker := time.NewTicker(timePeroid)
	defer ticker.Stop()

	// Channel to catch OS signals for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan bool)

	go func() {
		for {
			select {
			case <-ticker.C:
				// Trigger the expiration job
				err := StartExpirationJob()
				if err != nil {
					fmt.Println("Error in expiration job:", err)
				}
			case <-done:
				fmt.Println("Expiration job stopped.")
				return
			}
		}
	}()

	// Wait for shutdown signal
	<-signalChan
	fmt.Println("Shutdown signal received. Initiating graceful shutdown...")

	// Stop ticker and finish ongoing job
	ticker.Stop()
	done <- true

	fmt.Println("Graceful shutdown completed.")
}

func StartExpirationJob() error {
	log.Println("Running expiration job...")

	now := time.Now()
	expirationTime := now.AddDate(-cfg.SchedulerConfig.ExpireTimeYear, cfg.SchedulerConfig.ExpireTimeMonth, cfg.SchedulerConfig.ExpireTimeDay)

	// Fetch transactions older than the expiration time
	transactions, err := db.TransactionOlderThanGivenTime(expirationTime)
	if err != nil {
		log.Fatal(err)
	}

	// Process each transaction to expire points
	for _, txn := range transactions {
		err := db.ExpirePoints(txn.UserID, txn.TransactionID, txn.PointsEarned, txn.TransactionDate)
		if err != nil {
			log.Printf("Error expiring points for user %d, transaction %s: %v", txn.UserID, txn.TransactionID, err)
		}
	}
	if len(transactions) > 0 {
		log.Println("Expiration job completed total transaction updated - %s.", len(transactions))
	} else {

		log.Println("Expiration job completed, Not transaction updated.")
	}
	return nil
}
