package validations

import (
	"errors"

	"github.com/lakshay88/reward-management-system/database/models"
)

func ValidateTransaction(t models.Transaction) error {
	// Check each required field and return an error if it's missing
	if t.UserID <= 0 {
		return errors.New("user ID is required and must be greater than 0")
	}
	if t.TransactionAmount <= 0 {
		return errors.New("transaction amount is required and must be greater than 0")
	}
	if t.Category == "" {
		return errors.New("category is required")
	}
	if t.ProductCode == "" {
		return errors.New("product code is required")
	}
	if t.TransactionDate.IsZero() {
		return errors.New("transaction date is required")
	}

	return nil
}
