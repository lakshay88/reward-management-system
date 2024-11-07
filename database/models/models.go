package models

import "time"

// transaction
type Transaction struct {
	ID                int       `json:"id"`
	TransactionID     string    `json:"transaction_id"`
	UserID            int       `json:"user_id"`
	TransactionAmount float64   `json:"transaction_amount"`
	Category          string    `json:"category"`
	TransactionDate   time.Time `json:"transaction_date"`
	ProductCode       string    `json:"product_code"`
	PointsEarned      int       `json:"points_earned"`
	CreatedOn         time.Time `json:"createdOn"`
}
