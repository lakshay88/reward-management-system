package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	utils "github.com/lakshay88/reward-management-system/Utils"
	auth "github.com/lakshay88/reward-management-system/authentation"
	"github.com/lakshay88/reward-management-system/config"
	"github.com/lakshay88/reward-management-system/database"
	"github.com/lakshay88/reward-management-system/database/models"
	"github.com/lakshay88/reward-management-system/handlers/validations"
	"golang.org/x/crypto/bcrypt"
)

type Handlers struct{}

func NewHandlers() *Handlers {
	return &Handlers{}
}

func (h *Handlers) LoginRequest(cfg *config.AppConfig, db database.Database) (handlerFn http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request) {
		var loginRequest models.LoginRequest

		err := json.NewDecoder(r.Body).Decode(&loginRequest)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid input format"})
			return
		}

		var user *models.User
		// checking user exist or not
		user, err = db.GetUserByEmail(loginRequest.UserEmail, user)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid username or password"})
			return
		}

		// Password validation
		err = bcrypt.CompareHashAndPassword([]byte(user.UserPassword), []byte(loginRequest.Password))
		if err != nil {
			utils.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid username or password"})
			return
		}

		// token generation
		accessToken, refreshToken, err := auth.GenerateTokens(user.Email)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to generate tokens: %v", err), http.StatusInternalServerError)
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		})

	}
}

func (h *Handlers) RefreshToken(cfg *config.AppConfig, db database.Database) (handlerFn http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request) {
		var refreshRequest struct {
			RefreshToken string `json:"refresh_token"`
		}
		err := json.NewDecoder(r.Body).Decode(&refreshRequest)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid input format"})
			return
		}

		claims, err := auth.ValidateToken(refreshRequest.RefreshToken)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid or expired refresh token"})
			return
		}

		accessToken, _, err := auth.GenerateTokens(claims.Username)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"Failed to generate new access token:": err.Error()})
			http.Error(w, fmt.Sprintf("Failed to generate new access token: %v", err), http.StatusInternalServerError)
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
			"access_token": accessToken,
		})
	}
}

// CreateUser handles the creation of a new user
func (h *Handlers) CreateUser(cfg *config.AppConfig, db database.Database) (handlerFn http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid input format"})
			return
		}

		if err := validations.UserValidation(user); err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"Validation Error": err.Error()})
			return
		}

		user.UserPassword = utils.HashPassword(user.UserPassword)

		resp, err := db.CreateUser(&user)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"Failed to create user error": err.Error()})
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, resp)
	}
}

func (h *Handlers) GetUserByID(cfg *config.AppConfig, db database.Database) (handlerFn http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request) {
		var input models.GetUserInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid input"})
			return
		}

		var user *models.User
		var err error

		if input.UserID != 0 {
			user, err = db.GetUserByID(input.UserID, user)
			if err != nil {
				utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"Failed to get user error": err.Error()})
				return
			}
		} else if input.UserEmail != "" {
			user, err = db.GetUserByEmail(input.UserEmail, user)
			if err != nil {
				utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"Failed to get user by UserEmail error": err.Error()})
				return
			}
			user.UserPassword = ""
		} else {
			utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "id or username must be provided"})
			return
		}
		utils.RespondWithJSON(w, http.StatusOK, user)
	}
}

func (h *Handlers) AddTransactions(cfg *config.AppConfig, db database.Database) (handlerFn http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request) {

		var txn models.Transaction

		// incoming JSON
		err := json.NewDecoder(r.Body).Decode(&txn)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid input format"})
			return
		}

		// payload validation
		err = validations.ValidateTransaction(txn)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		txnCreated, err := db.AddTransaction(&txn)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
			"message":          "Transaction added successfully",
			"transaction_id":   txnCreated.TransactionID,
			"points_earned":    txnCreated.PointsEarned,
			"transaction_date": txnCreated.TransactionDate,
		})

	}
}

func (h *Handlers) PointBalance(cfg *config.AppConfig, db database.Database) (handlerFn http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request) {

		var pointBalanceRequest models.PointBalanceRequest
		if err := json.NewDecoder(r.Body).Decode(&pointBalanceRequest); err != nil {
			utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid input format"})
			return
		}

		if pointBalanceRequest.Page <= 0 {
			pointBalanceRequest.Page = 1
		}
		if pointBalanceRequest.Limit <= 0 {
			pointBalanceRequest.Limit = 10
		}

		balance, err := db.GetPointsBalance(pointBalanceRequest.UserID)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		// Fetch paginated points history
		history, err := db.GetPointsHistory(pointBalanceRequest.UserID, pointBalanceRequest.Page, pointBalanceRequest.Limit, "", "", "")
		if err != nil {
			utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		response := map[string]interface{}{
			"balance": balance,
			"history": history,
		}
		utils.RespondWithJSON(w, http.StatusOK, response)

	}
}

func (h *Handlers) GetPointsHistory(cfg *config.AppConfig, db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request models.PointsHistoryRequest

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
			return
		}

		if request.Page < 1 {
			request.Page = 1
		}
		if request.PageSize < 1 || request.PageSize > 100 {
			request.PageSize = 10
		}

		transactions, err := db.GetPointsHistory(request.UserID, request.Page, request.PageSize, request.StartDate, request.EndDate, request.TransactionType)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		response := map[string]interface{}{
			"transactions":       transactions,
			"total_transactions": len(transactions),
			"page":               request.Page,
			"page_size":          request.PageSize,
		}
		utils.RespondWithJSON(w, http.StatusOK, response)
	}
}

func (h *Handlers) RedeemPoints(cfg *config.AppConfig, db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var request models.RedeemPointsRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid input format"})
			return
		}

		// validation for null check
		if request.PointsToRedeem <= 0 {
			utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Points to redeem must be greater than zero"})
			return
		}

		// getting points
		totalAvailablePoints, err := db.GetAvailablePoints(request.UserID)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch points balance"})
			return
		}
		if request.PointsToRedeem > totalAvailablePoints {
			utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Insufficient points for redemption"})
			return
		}

		remainingBalance, err := db.DeductPoints(request.UserID, request.PointsToRedeem)
		if err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		err = db.LogPointsHistory(request.UserID, request.PointsToRedeem, "redeem", "Points redeemed for discount")
		if err != nil {
			utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to log points history"})
			return
		}

		utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
			"message":           "Points redeemed successfully",
			"points_redeemed":   request.PointsToRedeem,
			"remaining_balance": remainingBalance,
		})

	}
}
