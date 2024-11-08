package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/grid-stream-org/api/internal/app/handlers"
	"api/internal/app/middlewares"
	"api/internal/config"
	"api/pkg/database"
	

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	ctx := context.Background()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logging.Error(ctx, "failed to load configuration: %v", err)
		os.Exit(1)
	}

	// Set up database connection
	db, err := database.NewConnection(ctx, cfg.DatabaseURL)
	if err != nil {
		logging.Error(ctx, "failed to connect to database: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	// Create the router
	r := chi.NewRouter()

	// Apply middleware
	r.Use(middleware.Logger)
	r.Use(middlewares.RequestID)
	r.Use(middlewares.Recoverer)

	// Register routes
	r.Post("/posts", handlers.CreatePostHandler(db))

	// Set up the server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start the server
	go func() {
		logging.Info(ctx, "starting API server on port %d", cfg.Port)
		if err := srv.ListenAndServe(); err != nil {
			logging.Error(ctx, "API server failed: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	logging.Info(ctx, "shutting down API server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logging.Error(ctx, "failed to gracefully shut down server: %v", err)
	}

	logging.Info(ctx, "API server stopped")
}

// getenv := func(key string) string {
// 	switch key {
// 	case "MYAPP_FORMAT":
// 		return "markdown"
// 	case "MYAPP_TIMEOUT":
// 		return "5s"
// 	default:
// 		return ""
// }