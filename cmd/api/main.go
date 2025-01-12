package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grid-stream-org/api/internal/config"

	"github.com/grid-stream-org/api/internal/app/server"
	"github.com/grid-stream-org/api/pkg/database"
	"github.com/grid-stream-org/api/pkg/firebase"
	"github.com/grid-stream-org/api/pkg/logger"
	"github.com/pkg/errors"
)

func main() {
    ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, w io.Writer, args []string) error {
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
	bqClient, err := database.InitializeBigQueryClient(ctx, cfg, log)
	if err != nil {
		return errors.Wrap(err, "Failed init Big Queryu client")
	}
	defer database.CloseBigQueryClient(log)

	// Initialize Firebase Auth client
	firebaseAuth, err := firebase.InitializeFirebaseClient(ctx, cfg, log)
	if err != nil {
		return errors.Wrap(err, "failed to initialize Firebase Auth client")
	}

	// setup server handler
	handler := server.NewServer(cfg, bqClient, firebaseAuth, log)

    // Create the HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: handler,
	}

	// Start the server in a goroutine
	serverErrChan := make(chan error, 1)
	go func() {
		log.Info("Starting API server...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErrChan <- errors.Wrap(err, "server failed")
		}
		close(serverErrChan)
	}()

	log.Info("Application is running...")
	<-ctx.Done()

	log.Info("Shutting down...")
	// call shutdowns

	return nil
}
