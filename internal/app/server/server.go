package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/grid-stream-org/api/internal/config"
)

// NewServer sets up and returns an HTTP server
func NewServer(cfg *config.Config, log *slog.Logger) *http.Server {
	r := NewRouter(log)

	// Configure and return the HTTP server
	return &http.Server{
		Addr:        fmt.Sprintf(":%d", cfg.Port),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}
