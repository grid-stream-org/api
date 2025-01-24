package middlewares

import (
	"log/slog"
	"net/http"

	"github.com/grid-stream-org/api/internal/custom_error"
)

type HandlerFunc = func(w http.ResponseWriter, r *http.Request) error

// WrapHandler wraps a handler that returns an error to fit http.HandlerFunc.
func WrapHandler(handler HandlerFunc, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Call the handler
		err := handler(w, r)
		if err != nil {
			// Log error
			log.Error("Handler error", "error", err)

			// throw our error from custom_error.go
			if customErr, ok := err.(*custom_error.CustomError); ok {
				http.Error(w, customErr.Message, customErr.Code)
				return
			}

			// Default to internal server error
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
