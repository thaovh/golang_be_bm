package server

import (
	"context"

	authv1 "github.com/go-kratos/kratos-layout/api/auth/v1"
	countryv1 "github.com/go-kratos/kratos-layout/api/country/v1"
	helloworldv1 "github.com/go-kratos/kratos-layout/api/helloworld/v1"
	userv1 "github.com/go-kratos/kratos-layout/api/user/v1"
	"github.com/go-kratos/kratos-layout/internal/conf"
	"github.com/go-kratos/kratos-layout/internal/middleware"
	"github.com/go-kratos/kratos-layout/internal/service"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(c *conf.Server, greeter *service.GreeterService, user *service.UserService, auth *service.AuthService, country *service.CountryService, authConfig *conf.Auth, logger log.Logger) *http.Server {
	// Rate limiting for login endpoint
	loginRateLimit := middleware.LoginRateLimit()

	// Auth middleware for protected routes
	authMiddleware := middleware.AuthMiddleware([]byte(authConfig.JwtSecret))

	// Protected routes (require authentication)
	protectedPaths := []string{
		"/api/v1/users",
		"/api/v1/countries", // Country CRUD operations require authentication
		"/api/v1/auth/me",
		"/api/v1/auth/logout",
		"/api/v1/auth/revoke-all",
	}

	// Rate limited paths (login endpoint)
	rateLimitedPaths := []string{
		"/api/v1/auth/login",
		"/api/v1/auth/register",
	}

	var opts = []http.ServerOption{
		http.Middleware(
			recovery.Recovery(),
			// Request ID middleware (should be first to generate ID for all requests)
			middleware.RequestIDMiddleware(),
			// HTTP logging middleware (should be early to capture all requests)
			middleware.LoggingMiddleware(logger),
			// Apply rate limiting to login/register endpoints
			selector.Server(loginRateLimit).
				Match(func(ctx context.Context, operation string) bool {
					if req, ok := http.RequestFromServerContext(ctx); ok {
						path := req.URL.Path
						for _, p := range rateLimitedPaths {
							if path == p {
								return true
							}
						}
					}
					return false
				}).Build(),
			// Apply auth middleware to protected routes
			selector.Server(authMiddleware).
				Match(func(ctx context.Context, operation string) bool {
					if req, ok := http.RequestFromServerContext(ctx); ok {
						path := req.URL.Path
						for _, p := range protectedPaths {
							// Exact match or prefix match (for /api/v1/users/{id})
							if path == p || (len(path) > len(p) && path[:len(p)] == p && (path[len(p)] == '/' || path[len(p)] == '?')) {
								return true
							}
						}
					}
					return false
				}).Build(),
		),
	}
	if c.Http.Network != "" {
		opts = append(opts, http.Network(c.Http.Network))
	}
	if c.Http.Addr != "" {
		opts = append(opts, http.Address(c.Http.Addr))
	}
	if c.Http.Timeout != nil {
		opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
	}
	srv := http.NewServer(opts...)
	helloworldv1.RegisterGreeterHTTPServer(srv, greeter)
	userv1.RegisterUserServiceHTTPServer(srv, user)
	authv1.RegisterAuthServiceHTTPServer(srv, auth)
	countryv1.RegisterCountryServiceHTTPServer(srv, country)
	
	// Register Swagger UI
	RegisterSwaggerUI(srv)
	
	return srv
}
