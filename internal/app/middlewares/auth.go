package middlewares

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"firebase.google.com/go/auth"
)

type AuthMiddleware struct {
	FirebaseAuth *auth.Client
}

func NewAuthMiddleware(firebaseAuth *auth.Client, log *slog.Logger) *AuthMiddleware {
	return &AuthMiddleware{FirebaseAuth: firebaseAuth}
}

func (am *AuthMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// Verify the Firebase JWT token
		_, err := am.FirebaseAuth.VerifyIDToken(r.Context(), tokenString)
		if err != nil {
			fmt.Printf("err verrifying token: %v\n", err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
