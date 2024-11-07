package gateway

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/lakshay88/reward-management-system/config"
	"github.com/lakshay88/reward-management-system/database"
	"github.com/lakshay88/reward-management-system/gateway/routers"
)

type Gateway struct{}

func NewGateway() *Gateway {
	return &Gateway{}
}

func (g *Gateway) RegisterGateWayService(cfg *config.AppConfig, db database.Database) error {
	log.Println("Setting routers")

	// Initialize the chi engine
	r := chi.NewRouter()
	apiRouter := routers.NewRouter()
	apiRouter.RegisterRoutes(r, cfg, db)

	// http server
	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", cfg.ServerConfig.Port),
		Handler:      r,
		ReadTimeout:  2 * time.Minute,
		WriteTimeout: 0,
		IdleTimeout:  60 * time.Second,
	}

	// channel to listen for OS signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("Server started successfully on port", cfg.ServerConfig.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	<-signalChan
	log.Println("starting shutdown...")

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server Shutdown Error: %v", err)
	}
	log.Println("Server gracefully shut down")

	return nil
}
