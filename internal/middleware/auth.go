package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/CB-AutoStack/AutoStack/apps/api-inventory/internal/auth"
	"github.com/sirupsen/logrus"
)

// ContextKey is a type for context keys
type ContextKey string

const (
	// UserIDKey is the context key for user ID
	UserIDKey ContextKey = "userId"
	// EmailKey is the context key for email
	EmailKey ContextKey = "email"
)

// AuthMiddleware creates authentication middleware
func AuthMiddleware(jwtManager *auth.JWTManager, logger *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				logger.Warn("Missing authorization header")
				http.Error(w, "Unauthorized: missing authorization header", http.StatusUnauthorized)
				return
			}

			// Extract Bearer token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				logger.Warn("Invalid authorization header format")
				http.Error(w, "Unauthorized: invalid authorization header format", http.StatusUnauthorized)
				return
			}

			token := parts[1]
			claims, err := jwtManager.ValidateToken(token)
			if err != nil {
				logger.WithError(err).Warn("Invalid token")
				http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
				return
			}

			// Add user info to context
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, EmailKey, claims.Email)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
