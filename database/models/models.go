package models

import "time"

// User
type User struct {
	ID           int       `json:"id,omitempty"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	UserPassword string    `json:"user_password,omitempty"`
	CreatedOn    time.Time `json:"created_on,omitempty"`
}

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
	CreatedOn         time.Time `json:"created_on"`
}

type PointBalanceRequest struct {
	UserID int `json:"user_id"`
	Page   int `json:"page"`
	Limit  int `json:"limit"`
}

type PointsHistoryRequest struct {
	UserID          int    `json:"user_id"`
	StartDate       string `json:"start_date"`
	EndDate         string `json:"end_date"`
	TransactionType string `json:"transaction_type"`
	Page            int    `json:"page"`
	PageSize        int    `json:"page_size"`
}

type GetUserInput struct {
	UserID    int    `json:"userId,omitempty"`
	UserEmail string `json:"userEmail,omitempty"`
}

type PointsBalance struct {
	TotalPoints    int `json:"total_points"`
	PointsRedeemed int `json:"points_redeemed"`
}

type PointsHistory struct {
	Points     int       `json:"points"`
	PointsType string    `json:"points_type"` // earn, redeem, expire
	Reason     string    `json:"reason"`
	Date       time.Time `json:"date"`
}

type RedeemPointsRequest struct {
	UserID         int `json:"user_id"`
	PointsToRedeem int `json:"points_to_redeem"`
}

type LoginRequest struct {
	UserEmail string `json:"userEmail"`
	Password  string `json:"password"`
}
