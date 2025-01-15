package server

import (
	"log/slog"
	"net/http"

	"cloud.google.com/go/bigquery"
	"firebase.google.com/go/auth"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/grid-stream-org/api/internal/config"
)

// NewServer sets up and returns an HTTP server
func NewServer(
	cfg *config.Config,
	bqclient *bigquery.Client,
	fbclient *auth.Client,
	log *slog.Logger,
) http.Handler {
	// r := NewRouter(log, bqclient, fbclient)
	r := chi.NewRouter()

	addMidleware(r, log, fbclient, bqclient, cfg)
	AddRoutes(r, log, bqclient, fbclient)

	return r

}

func addMidleware(
	r *chi.Mux,
	log *slog.Logger,
	fbClient *auth.Client,
	bqClient *bigquery.Client,
	cfg *config.Config,
) {
	// TODO should be configured with conf maybe
	corsOptions := cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}

	r.Use(cors.Handler(corsOptions))
	r.Use(middleware.RequestID) // generate unique request id for each request
	r.Use(middleware.Logger)    // midleware logger to log incoming requests
	r.Use(middleware.Recoverer) // prevents server from crashing and responds with 500 error
}
