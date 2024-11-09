package routers

import (
	"github.com/go-chi/chi/v5"
	auth "github.com/lakshay88/reward-management-system/authentation"
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

	// User routes
	// User registration and login routes
	router.Post("/createUser", handlersInstance.CreateUser(cfg, db))
	router.Post("/login", handlersInstance.LoginRequest(cfg, db))
	router.Post("/refresh-token", handlersInstance.RefreshToken(cfg, db))

	protected := chi.NewRouter()
	protected.Use(auth.AuthMiddleware)
	router.Get("/user", handlersInstance.GetUserByID(cfg, db))
	// Add Transaction
	router.Post("/transaction/add", handlersInstance.AddTransactions(cfg, db))

	// Get Points balance
	router.Get("/points/balance", handlersInstance.PointBalance(cfg, db))

	// Redeem Point API
	router.Post("/points/redeem", handlersInstance.RedeemPoints(cfg, db))

	// Get Point History
	router.Post("/points/history", handlersInstance.GetPointsHistory(cfg, db))

	router.Mount("/", protected)
}
