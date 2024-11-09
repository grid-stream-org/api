package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// "github.com/grid-stream-org/api/internal/app/middlewares"

	"github.com/grid-stream-org/api/internal/config"

	// "github.com/grid-stream-org/api/internal/app/handlers"
	"github.com/grid-stream-org/api/pkg/database"
	"github.com/grid-stream-org/api/pkg/logger"
	"github.com/pkg/errors"

	"github.com/go-chi/chi/v5"
	//"github.com/go-chi/chi/v5/middleware"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		return errors.Wrap(err, "loading conf")
	}

	log, err := logger.Init(&cfg.Logger, nil)
	if err != nil {
		return errors.Wrap(err, "init logger")
	}

	// Set up big query connection
	if err := database.InitializeBigQueryClient(ctx, cfg, log); err != nil {
		log.Error("Failed to initialize BigQuery client")
	}
	defer database.CloseBigQueryClient(log)

	// Create the router
	r := chi.NewRouter()

	// Apply middleware
	// r.Use(middleware.Logger)
	// r.Use(middlewares.RequestID)
	// r.Use(middlewares.Recoverer)

	// Register routes
	// r.Post("/posts", handlers.CreatePostHandler(db))

	// Set up the server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// Start the server
	go func() {
		log.Info("starting API server on port %d", cfg.Port)
		if err := srv.ListenAndServe(); err != nil {
			log.Error("API server failed: %v", err)
		}
	}()

	log.Info("Application is running...")
	<-ctx.Done()

	log.Info("Shutting down...")
	// call shutdowns

	return nil
}