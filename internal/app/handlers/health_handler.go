package handlers

import (
	"net/http"
	"log/slog"
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
func (h *HealthHandler) HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	h.Log.Info("Health check endpoint hit")

	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		h.Log.Error("Failed to write response", "error", err)
	}
}
