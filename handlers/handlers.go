package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/lakshay88/reward-management-system/config"
	"github.com/lakshay88/reward-management-system/database"
	"github.com/lakshay88/reward-management-system/database/models"
	"github.com/lakshay88/reward-management-system/handlers/validations"
)

type Handlers struct{}

func NewHandlers() *Handlers {
	return &Handlers{}
}

func (h *Handlers) AddTransactions(cfg *config.AppConfig, db database.Database) (handlerFn http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request) {

		var txn models.Transaction

		// incoming JSON
		err := json.NewDecoder(r.Body).Decode(&txn)
		if err != nil {
			respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid input format"})
			return
		}

		// payload validation
		err = validations.ValidateTransaction(txn)
		if err != nil {
			respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		if err := db.AddTransaction(&txn); err != nil {
			respondWithJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		respondWithJSON(w, http.StatusOK, map[string]string{"message": "Transaction added successfully"})

	}
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}
