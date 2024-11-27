package handlers

import (
	"log/slog"
	"net/http"

)

func HealthCheckHandler(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("Health check endpoint hit")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}
