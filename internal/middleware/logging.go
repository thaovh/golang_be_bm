package middleware

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// LoggingMiddleware logs HTTP requests with detailed information
func LoggingMiddleware(logger log.Logger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			startTime := time.Now()

			// Extract HTTP request information
			var method, path, ip, userAgent string
			if httpReq, ok := http.RequestFromServerContext(ctx); ok {
				method = httpReq.Method
				path = httpReq.URL.Path
				ip = extractIP(httpReq.RemoteAddr)
				userAgent = httpReq.UserAgent()
			}

			// Get request ID
			requestID, _ := GetRequestIDFromContext(ctx)

			// Get user context if available
			userID, hasUserID := GetUserIDFromContext(ctx)
			userEmail, hasUserEmail := ctx.Value(UserEmailKey).(string)
			userRole, hasUserRole := GetUserRoleFromContext(ctx)

			// Execute handler
			resp, err := handler(ctx, req)

			// Calculate duration
			duration := time.Since(startTime)
			durationMs := duration.Milliseconds()

			// Determine status code
			statusCode := 200
			if err != nil {
				if kratosErr, ok := err.(interface{ Code() int }); ok {
					statusCode = kratosErr.Code()
				} else {
					statusCode = 500
				}
			}

			// Build log fields
			keyvals := []interface{}{
				"msg", "HTTP request",
				"method", method,
				"path", path,
				"ip", ip,
				"user_agent", userAgent,
				"request_id", requestID,
				"status_code", statusCode,
				"duration_ms", durationMs,
			}

			// Add user context if available
			if hasUserID {
				keyvals = append(keyvals, "user_id", userID.String())
			}
			if hasUserEmail {
				keyvals = append(keyvals, "user_email", userEmail)
			}
			if hasUserRole {
				keyvals = append(keyvals, "user_role", userRole)
			}

			// Log based on status code
			if statusCode >= 500 {
				// Error level for server errors
				log.NewHelper(logger).WithContext(ctx).Log(log.LevelError, keyvals...)
			} else if statusCode >= 400 {
				// Warning level for client errors
				log.NewHelper(logger).WithContext(ctx).Log(log.LevelWarn, keyvals...)
			} else {
				// Info level for successful requests
				log.NewHelper(logger).WithContext(ctx).Log(log.LevelInfo, keyvals...)
			}

			return resp, err
		}
	}
}

// extractIP extracts IP address from RemoteAddr
func extractIP(remoteAddr string) string {
	if remoteAddr == "" {
		return "unknown"
	}

	// Handle IPv6 format [::1]:port
	if len(remoteAddr) > 0 && remoteAddr[0] == '[' {
		endIdx := -1
		for i := 1; i < len(remoteAddr); i++ {
			if remoteAddr[i] == ']' {
				endIdx = i
				break
			}
		}
		if endIdx > 0 {
			return remoteAddr[1:endIdx]
		}
	}

	// Handle IPv4 format 127.0.0.1:port
	for i := len(remoteAddr) - 1; i >= 0; i-- {
		if remoteAddr[i] == ':' {
			return remoteAddr[:i]
		}
	}

	return remoteAddr
}

