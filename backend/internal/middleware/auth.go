package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JwtKey []byte

// MaxTokenAgeMonths defines the maximum age of a JWT token in months
const MaxTokenAgeMonths = 3

// maxTokenAge defines the maximum age of a JWT token as a duration
var maxTokenAge = time.Duration(MaxTokenAgeMonths) * 30 * 24 * time.Hour

// ContextKey is a custom type for context keys to avoid collisions
type ContextKey string

const UserIDKey ContextKey = "userId"

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// authMiddleware validates the JWT token and adds the user info to the request context.
// It returns a 401 error for invalid or expired tokens, but proceeds without user info for missing tokens.
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		if tokenString != "" {
			// The header is typically in the format "Bearer <token>".
			if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
				tokenString = tokenString[7:]

				// slog.Info(tokenString)

				// Initialize a new Claims struct.
				claims := &Claims{}

				// Parse the token.
				token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
					// Check the signing method.
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
					}
					return JwtKey, nil
				})

				

				// Check if the token is valid and not expired based on age.
				isValid := err == nil && token.Valid
				if isValid && claims.IssuedAt != nil {
					isValid = time.Since(claims.IssuedAt.Time) <= maxTokenAge
				}

				if isValid {
					ctx := r.Context()
					ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
					next.ServeHTTP(w, r.WithContext(ctx))
					return // End the middleware chain here.
				} else {
					// Token is invalid or expired, return 401.
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
			}
		}

		// If the token is missing, we simply proceed without adding the user ID to the context.
		next.ServeHTTP(w, r)
	})
}

func GetUserId(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}
