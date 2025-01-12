// Defines Routes and Middleware
package server

import (
	"log/slog"

	"cloud.google.com/go/bigquery"
	"firebase.google.com/go/auth"
	"github.com/go-chi/chi/v5"
	"github.com/grid-stream-org/api/internal/app/handlers"
	"github.com/grid-stream-org/api/internal/app/middlewares"
	"github.com/grid-stream-org/api/internal/app/repositories"
)

func AddRoutes(
	r *chi.Mux,
	log *slog.Logger,
	bqClient *bigquery.Client,
	fbClient *auth.Client,
) {
	// initialize repositories
	projectRepo := repositories.NewProjectRepository(bqClient)

	// init handlers
	projectHandlers := handlers.NewProjectHandlers(projectRepo, log)
	healthHandler := handlers.NewHealthHandler(log)

	authMiddleware := middlewares.NewAuthMiddleware(fbClient)

	// Health check route
	r.Get("/health", middlewares.WrapHandler(healthHandler.HealthCheckHandler))

	// Projects routes
	r.Route("/projects", func(r chi.Router) {
		r.Use(authMiddleware.Handler) // JWT auth middleware for projects
		r.Get("/{id}", middlewares.WrapHandler(projectHandlers.GetProjectHandler))
		r.Put("/{id}", middlewares.WrapHandler(projectHandlers.UpdateProjectHandler))
	})

}
