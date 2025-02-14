package middlewares

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go/auth"
	"github.com/grid-stream-org/api/pkg/firebase"
)

type contextKey string

const userKey contextKey = "users"

type AuthMiddleware struct {
	FirebaseAuth *auth.Client
	Firestore    *firestore.Client
}

func NewAuthMiddleware(firebaseClient firebase.FirebaseClient, log *slog.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		FirebaseAuth: firebaseClient.Auth(),
		Firestore:    firebaseClient.Firestore(),
	}
}

func (am *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
// RequireAuth is used when you only need authentication, no role check
func (am *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return am.RequireRole()(next)
}

func getRoleFromInt(value int) string {
	switch value {
	case 0:
		return "Utility"
	case 1:
		return "Residential"
	case 2:
		return "Technician"
	default:
		return ""
	}
}

// RequireRole checks authentication and role(s)
func (am *AuthMiddleware) RequireRole(requiredRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			// Get and verify token
			token, err := am.FirebaseAuth.VerifyIDToken(r.Context(), getTokenFromHeader(r))
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// If no roles required, proceed after auth check
			if len(requiredRoles) == 0 {
				next.ServeHTTP(w, r)
				return
			}

			// Get user from Firestore
			userDoc, err := am.Firestore.Collection(string(userKey)).Doc(token.UID).Get(r.Context())
			if err != nil {
				fmt.Printf("error fetching user doc: %v\n", err)
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			// Get user's role
			userRole, ok := userDoc.Data()["role"]
			if !ok {
				fmt.Printf("User role %d\n", userDoc.Data()["role"])
				fmt.Printf("user has no role field: %v\n", userDoc.Data())
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			val, ok := userRole.(int64)
			if !ok {
				http.Error(w, "Failed to resolve user role", http.StatusInternalServerError)
				return
			}

			ro := getRoleFromInt(int(val))

			// Check if user's role matches any required role
			hasRequiredRole := false
			for _, requiredRole := range requiredRoles {
				if ro == requiredRole {
					hasRequiredRole = true
					break
				}
			}

			if !hasRequiredRole {
				fmt.Printf("user lacks required role. Has: %v, Needs one of: %v\n", ro, requiredRoles)
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			// Add user data to context
			ctx := context.WithValue(r.Context(), userKey, userDoc.Data())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func getTokenFromHeader(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		return authHeader[7:]
	}
	return ""
}
