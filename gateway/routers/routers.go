package routers

import (
	"github.com/go-chi/chi/v5"
	"github.com/lakshay88/reward-management-system/config"
	"github.com/lakshay88/reward-management-system/database"
	"github.com/lakshay88/reward-management-system/handlers"
)

type Routers struct{}

func NewRouter() *Routers {
	return &Routers{}
}

func (r *Routers) RegisterRoutes(router *chi.Mux, cfg *config.AppConfig, db database.Database) {
	handlersInstance := handlers.NewHandlers()

	// Add Transaction
	router.Post("/transaction/add", handlersInstance.AddTransactions(cfg, db))

}
