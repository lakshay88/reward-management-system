package gateway

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/lakshay88/reward-managment-system/config"
	"github.com/lakshay88/reward-managment-system/gateway/routers"
)

type Gateway struct{}

func NewGateway() *Gateway {
	return &Gateway{}
}

func (g *Gateway) RegisterGateWayService(cfg config.RestServerConfig) error {
	log.Println("Setting routers")

	// Initialize the chi engine
	r := chi.NewRouter()
	apiRouter := routers.NewRouter()
	apiRouter.RegisterRoutes(r)

	// http server
	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", cfg.Port),
		Handler:      r,
		ReadTimeout:  2 * time.Minute,
		WriteTimeout: 0,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Println("Server started successfully on port", cfg.Port)

	select {} // This will block indefinitely. You can replace this with other logic if needed.

	return nil
}
