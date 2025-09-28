package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

// RequireRole middleware to check user roles
// This runs AFTER jwtauth.Verifier and jwtauth.Authenticator
func RequireRole(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get token and claims from context (set by jwtauth.Verifier)
			_, claims, err := jwtauth.FromContext(r.Context())
			if err != nil {
				http.Error(w, "Token error", http.StatusUnauthorized)
				return
			}

			// Check if user has required role
			userRole, ok := claims["role"].(string)
			if !ok {
				http.Error(w, "Role not found in token", http.StatusForbidden)
				return
			}

			// Allow admin to access everything, or check specific role
			if userRole != requiredRole && userRole != "admin" {
				http.Error(w, "Insufficient permissions", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Helper function to get user claims from context
func GetUserFromContext(ctx context.Context) (map[string]interface{}, bool) {
	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return nil, false
	}
	return claims, true
}

// Optional: Custom JSON error authenticator (alternative to jwtauth.Authenticator)
// Use this if you want JSON error responses instead of plain text
func JSONAuthenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())

		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "Invalid token", "message": "Token verification failed"}`))
			return
		}

		if token == nil || !token.Valid {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"error": "Unauthorized", "message": "Valid token required"}`))
			return
		}

		// Token is valid, proceed to next handler
		next.ServeHTTP(w, r)
	})
}
