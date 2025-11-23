package middleware

import (
	"context"

	"github.com/gofrs/uuid/v5"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// RequestIDKey is the context key for request ID
type RequestIDKey struct{}

// RequestIDMiddleware generates a unique request ID for each request
func RequestIDMiddleware() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			// Generate unique request ID
			requestID := uuid.Must(uuid.NewV7()).String()

			// Add request ID to context
			ctx = context.WithValue(ctx, RequestIDKey{}, requestID)

			// Add request ID to HTTP response header if available
			if httpReq, ok := http.RequestFromServerContext(ctx); ok {
				httpReq.Header.Set("X-Request-ID", requestID)
			}

			return handler(ctx, req)
		}
	}
}

// GetRequestIDFromContext extracts request ID from context
func GetRequestIDFromContext(ctx context.Context) (string, bool) {
	requestID, ok := ctx.Value(RequestIDKey{}).(string)
	return requestID, ok
}

