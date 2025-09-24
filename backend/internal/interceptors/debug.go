package interceptors

import (
	"context"
	"log/slog"

	"connectrpc.com/connect"
)

// NewDebugInterceptor creates a Connect interceptor for logging requests
func NewDebugInterceptor() connect.UnaryInterceptorFunc {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(
			ctx context.Context,
			req connect.AnyRequest,
		) (connect.AnyResponse, error) {
			// Call the next handler
			res, err := next(ctx, req)

			// Log the request details
			status := "success"
			if err != nil {
				status = "error"
			}

			slog.Info("Request handled",
				"procedure", req.Spec().Procedure,
				"status", status,
			)

			return res, err
		})
	}
	return connect.UnaryInterceptorFunc(interceptor)
}
