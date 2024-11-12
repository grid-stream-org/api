// Defines Routes and Middleware
package server

import (
	"log/slog"

	"github.com/grid-stream-org/api/internal/app/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// NewRouter creates and configures a chi router
func NewRouter(log *slog.Logger) *chi.Mux {
	// Create a new chi router
	r := chi.NewRouter()

	// Apply global middleware
	r.Use(middleware.RequestID) // generate unique request id for each request
	r.Use(middleware.Logger) // midleware logger to log incoming requests
	r.Use(middleware.Recoverer) // prevents server from crashing and responds with 500 error

	// Define routes
	r.Get("/health", handlers.HealthCheckHandler(log))

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
