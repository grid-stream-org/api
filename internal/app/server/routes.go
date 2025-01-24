// Defines Routes and Middleware
package server

import (
	"log/slog"

	"firebase.google.com/go/auth"
	"github.com/go-chi/chi/v5"
	"github.com/grid-stream-org/api/internal/app/handlers/v1"
	"github.com/grid-stream-org/api/internal/app/middlewares"
	"github.com/grid-stream-org/api/internal/app/repositories"
	"github.com/grid-stream-org/batcher/pkg/bqclient"
)

func AddRoutes(
	r *chi.Mux,
	log *slog.Logger,
	bqClient bqclient.BQClient,
	fbClient *auth.Client,
) {
	// initialize repositories
	projectRepo := repositories.NewProjectRepository(bqClient, log)

	// init handlers
	projectHandlers := handlers.NewProjectHandlers(projectRepo, log)
	healthHandler := handlers.NewHealthHandler(log)

	authMiddleware := middlewares.NewAuthMiddleware(fbClient, log)

	// Health check route
	r.Get("/health", middlewares.WrapHandler(healthHandler.HealthCheckHandler, log))

	r.Route("/v1", func(r chi.Router) {
		r.Route("/projects", func(r chi.Router) {
			r.Use(authMiddleware.Handler) // JWT auth middleware for projects
			r.Get("/{id}", middlewares.WrapHandler(projectHandlers.GetProjectHandler, log))
			r.Put("/{id}", middlewares.WrapHandler(projectHandlers.UpdateProjectHandler, log))
			r.Post("/", middlewares.WrapHandler(projectHandlers.CreateProjectHandler, log))
		})
	})

}
