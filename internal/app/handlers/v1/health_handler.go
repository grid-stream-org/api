package handlers

import (
	"log/slog"
	"net/http"

	"github.com/grid-stream-org/api/internal/errors"
)

// HealthHandler contains the logger for health check handling.
type HealthHandler struct {
	Log *slog.Logger
}

// NewHealthHandler creates a new instance of HealthHandler.
func NewHealthHandler(log *slog.Logger) *HealthHandler {
	return &HealthHandler{Log: log}
}

// HealthCheckHandler handles the health check endpoint.
// Does not call a repository because it's simple enough
func (h *HealthHandler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) error {
	h.Log.Info("Health check endpoint hit")

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		h.Log.Error("Failed to write response", "error", err)
        return errors.New(http.StatusInternalServerError, "Unexpected error hitting health endpoint")
	}
    return nil
}
