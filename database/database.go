package database

import (
	"time"

	"github.com/lakshay88/reward-management-system/database/models"
)

type Database interface {
	// implement Database methods

	// User Functions
	CreateUser(*models.User) (*models.User, error)
	GetUserByID(int, *models.User) (*models.User, error)
	GetUserByEmail(string, *models.User) (*models.User, error)

	// Add Transaction
	AddTransaction(*models.Transaction) (*models.Transaction, error)

	// Points Balance
	GetPointsBalance(int) (models.PointsBalance, error)
	GetPointsHistory(int, int, int, string, string, string) ([]models.PointsHistory, error)

	// Reward redeem
	GetAvailablePoints(int) (int, error)
	DeductPoints(int, int) (int, error)
	LogPointsHistory(int, int, string, string) error

	// Exprite
	TransactionOlderThanGivenTime(time.Time) ([]models.Transaction, error)
	ExpirePoints(int, string, int, time.Time) error

	Close() error
}
