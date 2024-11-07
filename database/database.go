package database

import (
	"github.com/lakshay88/reward-management-system/database/models"
)

type Database interface {
	// implement Database methods
	AddTransaction(*models.Transaction) error
	Close() error
}
