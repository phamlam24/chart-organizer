package middleware

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("JWT_KEY"))

// MaxTokenAgeMonths defines the maximum age of a JWT token in months
const MaxTokenAgeMonths = 3

// ContextKey is a custom type for context keys to avoid collisions
type ContextKey string

const UserIDKey ContextKey = "userID"

type claims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// AuthMiddleware is a standard Go HTTP middleware
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		tokenString = tokenString[7:] // "Bearer " is 7 characters

		claims := &claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if token is older than MaxTokenAgeMonths
		if claims.IssuedAt != nil {
			tokenAge := time.Since(claims.IssuedAt.Time)
			maxAge := time.Duration(MaxTokenAgeMonths) * 30 * 24 * time.Hour // Approximate months to hours

			if tokenAge > maxAge {
				http.Error(w, "Token expired - too old", http.StatusUnauthorized)
				return
			}
		}

		// Attach user ID to the request context
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
