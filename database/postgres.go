package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/lakshay88/reward-management-system/config"
	"github.com/lakshay88/reward-management-system/database/models"
	_ "github.com/lib/pq"
)

type PostgresDB struct {
	connection *sql.DB
}

func ConnectionToPostgres(cfg config.DatabaseConfig) (Database, error) {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	connection, err := sql.Open(cfg.Driver, connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}

	// Set connection limits
	connection.SetMaxOpenConns(25)
	connection.SetMaxIdleConns(25)
	connection.SetConnMaxLifetime(5 * time.Minute)

	if err := connection.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresDB{connection: connection}, nil
}

func (db *PostgresDB) Close() error {
	return db.connection.Close()
}

func (db *PostgresDB) AddTransaction(txn *models.Transaction) error {
	txn.TransactionID = uuid.New().String()
	if txn.TransactionDate.IsZero() {
		txn.TransactionDate = time.Now()
	}
	return nil
}
