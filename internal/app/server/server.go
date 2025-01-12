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

	addMidleware(r, log, fbclient, bqclient)
    AddRoutes(r, log, bqclient, fbclient)

	// Configure and return the HTTP server
	// return &http.Server{
	// 	Addr:         fmt.Sprintf(":%d", cfg.Port),
	// 	Handler:      r,
	// 	ReadTimeout:  10 * time.Second,
	// 	WriteTimeout: 10 * time.Second,
	// }
    
    return r

}

func addMidleware(
	r *chi.Mux,
	log *slog.Logger,
	fbClient *auth.Client,
	bqClient *bigquery.Client,
) {
	// TODO should be configured with conf maybe
	corsOptions := cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // frontend
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
