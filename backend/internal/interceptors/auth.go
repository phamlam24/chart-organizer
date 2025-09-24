package interceptors

import (
	"context"
	"fmt"
	"time"

	"connectrpc.com/connect"
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

// NewAuthInterceptor creates a Connect interceptor that validates JWT tokens
// and adds user info to the request context.
func NewAuthInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			// Extract Authorization header
			tokenString := req.Header().Get("Authorization")

			if tokenString != "" && len(tokenString) > 7 && tokenString[:7] == "Bearer " {
				tokenString = tokenString[7:]

				claims := &Claims{}

				token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
					}
					return JwtKey, nil
				})

				isValid := err == nil && token.Valid
				if isValid && claims.IssuedAt != nil {
					isValid = time.Since(claims.IssuedAt.Time) <= maxTokenAge
				}

				if isValid {
					ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
				}
			}

			return next(ctx, req)
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}

func GetUserId(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value(UserIDKey).(string)
	return userID, ok
}
