// Defines Routes and Middleware
package server

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/grid-stream-org/api/internal/app/handlers/v1"
	"github.com/grid-stream-org/api/internal/app/middlewares"
	"github.com/grid-stream-org/api/internal/app/repositories"
	"github.com/grid-stream-org/api/pkg/firebase"
	"github.com/grid-stream-org/batcher/pkg/bqclient"
)

func AddRoutes(
	r *chi.Mux,
	log *slog.Logger,
	bqClient bqclient.BQClient,
	fbClient firebase.FirebaseClient,
) {
	// initialize repositories
	projectRepo := repositories.NewProjectRepository(bqClient, log)
	utilRepo := repositories.NewUtilityRepository(bqClient, log)
	contractRepo := repositories.NewContractRepository(bqClient, log)
	derMetaRepo := repositories.NewDERMetadataRepository(bqClient, log)

	// init handlers
	projectHandlers := handlers.NewProjectHandlers(projectRepo, log)
	utilHandlers := handlers.NewUtilityRepository(utilRepo, log)
	contractHandlers := handlers.NewContractHandlers(contractRepo, log)
	healthHandler := handlers.NewHealthHandler(log)
	derHandler := handlers.NewDERMetadataHandlers(derMetaRepo, log)

	// init middlewares
	authMiddleware := middlewares.NewAuthMiddleware(fbClient, log)
	r.Use(middlewares.PerClientRateLimiter)

	// Health check route
	r.Get("/health", middlewares.WrapHandler(healthHandler.HealthCheckHandler, log))

	r.Route("/v1", func(r chi.Router) {
		r.Route("/projects", func(r chi.Router) {
			// GET and PUT: only need "Residential"
			r.With(authMiddleware.RequireAuth).Get("/{id}", middlewares.WrapHandler(projectHandlers.GetProjectHandler, log))
			r.With(authMiddleware.RequireRole("Residential")).Put("/{id}", middlewares.WrapHandler(projectHandlers.UpdateProjectHandler, log))

			// POST and DELETE: only "Utility"
			r.With(authMiddleware.RequireRole("Technician")).Post("/", middlewares.WrapHandler(projectHandlers.CreateProjectHandler, log))
			r.With(authMiddleware.RequireRole("Technician")).Delete("/{id}", middlewares.WrapHandler(projectHandlers.DeleteProjectHandler, log))

		})

		r.Route("/utilities", func(r chi.Router) {
			r.With(authMiddleware.RequireRole("Technician", "Residential")).Get("/{id}", middlewares.WrapHandler(utilHandlers.GetUtilityHandler, log))
			r.With(authMiddleware.RequireRole("Technician")).Post("/", middlewares.WrapHandler(utilHandlers.CreateUtilityHandler, log))
			r.With(authMiddleware.RequireRole("Utility")).Put("/{id}", middlewares.WrapHandler(utilHandlers.UpdateUtilityHandler, log))
			r.With(authMiddleware.RequireRole("Utility")).Delete("/{id}", middlewares.WrapHandler(utilHandlers.DeleteUtilityHandler, log))
		})

		r.Route("/contracts", func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth)
			r.Get("/{id}", middlewares.WrapHandler(contractHandlers.GetContractHandler, log))
			r.Put("/{id}", middlewares.WrapHandler(contractHandlers.UpdateContractHandler, log))
			r.Delete("/{id}", middlewares.WrapHandler(contractHandlers.DeleteContractHandler, log))
			r.Post("/", middlewares.WrapHandler(contractHandlers.CreateContractHandler, log))
		})

		r.Route("/der-metadata", func(r chi.Router) {
			r.Use(authMiddleware.RequireAuth)
			r.Get("/{id}", middlewares.WrapHandler(derHandler.GetDERMetadataHandler, log))
			r.Put("/{id}", middlewares.WrapHandler(derHandler.UpdateDERMetadataHandler, log))
			r.Delete("/{id}", middlewares.WrapHandler(derHandler.DeleteDERMetadataHandler, log))
			r.Post("/", middlewares.WrapHandler(derHandler.CreateDERMetadataHandler, log))
		})
	})

}
