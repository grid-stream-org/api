// auth middleware for service to serevice communication
package middlewares

import (
	"context"
	"log"
	"net/http"
	"strings"

	"google.golang.org/api/idtoken"
)

// Expected service account email, CHANGE EMAIL WHEN SA CREATED
const expectedServiceAccount = "validator-service@my-project.iam.gserviceaccount.com"

// Expected audience (Cloud Run URL)
const expectedAudience = "https://my-api-url.run.app"

// ServiceAuthMiddleware verifies requests from the Validator service or any other service that may want to communicate with api
func ServiceAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := tokenParts[1]

		// Validate the ID token
		payload, err := idtoken.Validate(context.Background(), tokenString, expectedAudience)
		if err != nil {
			log.Println("Token validation failed:", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Ensure the request comes from the correct service account
		if payload.Claims["email"] != expectedServiceAccount {
			log.Println("Unauthorized service account:", payload.Claims["email"])
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// Token is valid, continue request
		next.ServeHTTP(w, r)
	})
}
