package validations

import (
	"errors"
	"regexp"

	"github.com/lakshay88/reward-management-system/database/models"
)

func UserValidation(usr models.User) error {
	if usr.Username == "" {
		return errors.New("username is required")
	}

	if !isValidEmail(usr.Email) {
		return errors.New("invalid email format")
	}

	if len(usr.UserPassword) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	return nil
}

func isValidEmail(email string) bool {
	const emailRegex = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

func ValidateTransaction(t models.Transaction) error {
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

	return nil
}
