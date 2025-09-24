package interceptors

import (
	"context"
	"net/http"

	"connectrpc.com/connect"
)

// NewCORSInterceptor creates a Connect interceptor that handles CORS headers
func NewCORSInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			// Call the next handler first
			res, err := next(ctx, req)

			// Add CORS headers to the response only if response is not nil
			if res != nil && res.Header() != nil {
				res.Header().Set("Access-Control-Allow-Origin", "*")
				res.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				res.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				res.Header().Set("Access-Control-Max-Age", "86400")
			}

			return res, err
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}

// CORSHandler wraps an http.Handler to add CORS headers for preflight requests
func CORSHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		h.ServeHTTP(w, r)
	})
}
