// Defines Routes and Middleware
package server

import (
	"log/slog"

	"cloud.google.com/go/bigquery"
	"firebase.google.com/go/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/grid-stream-org/api/internal/app/handlers"
	"github.com/grid-stream-org/api/internal/app/middlewares"
	"github.com/grid-stream-org/api/internal/app/repositories"
)

// NewRouter creates and configures a chi router
func NewRouter(log *slog.Logger, bqclient *bigquery.Client, fbclient *auth.Client) *chi.Mux {
	// Create a new chi router
	r := chi.NewRouter()

    // Coors config, I hate coors
    corsOptions := cors.Options{
        AllowedOrigins:     []string{"http://localhost:5173"}, // frontend
        AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:     []string{"Accept", "Authorization", "Content-Type"},
        ExposedHeaders:     []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
    }

	// Apply global middleware
    r.Use(cors.Handler(corsOptions))
	r.Use(middleware.RequestID) // generate unique request id for each request
	r.Use(middleware.Logger)    // midleware logger to log incoming requests
	r.Use(middleware.Recoverer) // prevents server from crashing and responds with 500 error
    // initialize auth middleware
    authMiddleware := middlewares.NewAuthMiddleware(fbclient)
    
    // initialize repository and handler for Projects
	projectRepo := repositories.NewProjectRepository(bqclient)
	projectHandlers := handlers.NewProjectHandlers(projectRepo, log)

	// Define routes
	r.Get("/health", handlers.HealthCheckHandler(log))
	r.Route("/projects", func(r chi.Router) {
        r.Use(authMiddleware.Handler) // Apply AuthMiddleware to verify jwt token
		r.Get("/{id}", projectHandlers.GetProjectHandler)
		r.Put("/{id}", projectHandlers.UpdateProjectHandler)
	})

	// Example route handler for /posts
	// this gives an example for a Post, Getters, Put and Delete
	// r.Route("/posts", func(r chi.Router) {
	// 	r.Post("/", handlers.CreatePostHandler)
	// 	r.Get("/", handlers.ListPostsHandler)
	// 	r.Get("/{id}", handlers.GetPostHandler)
	// 	r.Put("/{id}", handlers.UpdatePostHandler)
	// 	r.Delete("/{id}", handlers.DeletePostHandler)
	// })

	return r
}
