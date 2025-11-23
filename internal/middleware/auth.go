package middleware

import (
	"context"

	"github.com/go-kratos/kratos-layout/internal/pkg/jwt"
	"github.com/gofrs/uuid/v5"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// ContextKey for storing user info in context
type contextKey string

const (
	UserIDKey  contextKey = "user_id"
	UserEmailKey contextKey = "user_email"
	UserRoleKey contextKey = "user_role"
)

// AuthMiddleware validates JWT token and adds user info to context
func AuthMiddleware(jwtSecret []byte) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			// Extract token from HTTP request
			if httpReq, ok := http.RequestFromServerContext(ctx); ok {
				authHeader := httpReq.Header.Get("Authorization")
				if authHeader == "" {
					return nil, errors.Unauthorized("UNAUTHORIZED", "missing authorization header")
				}

				token, err := jwt.ExtractTokenFromHeader(authHeader)
				if err != nil {
					return nil, errors.Unauthorized("UNAUTHORIZED", err.Error())
				}

				// Validate token
				claims, err := jwt.ValidateToken(token, jwtSecret)
				if err != nil {
					if err == jwt.ErrExpiredToken {
						return nil, errors.Unauthorized("TOKEN_EXPIRED", "token has expired")
					}
					return nil, errors.Unauthorized("TOKEN_INVALID", "invalid token")
				}

				// Add user info to context
				ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
				ctx = context.WithValue(ctx, UserEmailKey, claims.Email)
				ctx = context.WithValue(ctx, UserRoleKey, claims.Role)

				return handler(ctx, req)
			}

			// For gRPC, token might be in metadata
			// This is a simplified version, can be extended for gRPC
			return handler(ctx, req)
		}
	}
}

// RequireRole middleware checks if user has required role
func RequireRole(roles ...string) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			userRole, ok := ctx.Value(UserRoleKey).(string)
			if !ok {
				return nil, errors.Unauthorized("UNAUTHORIZED", "user role not found")
			}

			// Check if user role is in allowed roles
			allowed := false
			for _, role := range roles {
				if userRole == role {
					allowed = true
					break
				}
			}

			if !allowed {
				return nil, errors.Forbidden("FORBIDDEN", "insufficient permissions")
			}

			return handler(ctx, req)
		}
	}
}

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(UserIDKey).(uuid.UUID)
	return userID, ok
}

// GetUserRoleFromContext extracts user role from context
func GetUserRoleFromContext(ctx context.Context) (string, bool) {
	role, ok := ctx.Value(UserRoleKey).(string)
	return role, ok
}

