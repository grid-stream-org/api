package middlewares

import (
	"log"
	"net/http"

	"github.com/grid-stream-org/api/internal/errors"
)

// WrapHandler wraps a handler that returns an error to fit http.HandlerFunc.
func WrapHandler(handler func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Call the handler
		err := handler(w, r)
		if err != nil {
			// Log error
			log.Printf("Handler error: %v", err)

			// throw our error from custom_error.go
			if customErr, ok := err.(*errors.CustomError); ok {
				http.Error(w, customErr.Message, customErr.Code)
				return
			}

			// Default to internal server error
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	}
}
