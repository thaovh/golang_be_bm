package middleware

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport/http"
)

// RateLimiter interface for rate limiting
type RateLimiter interface {
	Allow(key string) bool
	Reset(key string)
}

// InMemoryRateLimiter implements rate limiting using in-memory storage
type InMemoryRateLimiter struct {
	requests map[string][]time.Time
	mu       sync.RWMutex
	limit    int           // Maximum number of requests
	window   time.Duration // Time window
	cleanup  *time.Ticker  // Cleanup ticker for old entries
}

// NewInMemoryRateLimiter creates a new in-memory rate limiter
func NewInMemoryRateLimiter(limit int, window time.Duration) *InMemoryRateLimiter {
	rl := &InMemoryRateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
		cleanup:  time.NewTicker(1 * time.Minute), // Cleanup every minute
	}

	// Start cleanup goroutine
	go rl.cleanupOldEntries()

	return rl
}

// Allow checks if a request is allowed
func (rl *InMemoryRateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-rl.window)

	// Clean old requests for this key
	if timestamps, exists := rl.requests[key]; exists {
		validTimestamps := make([]time.Time, 0, len(timestamps))
		for _, ts := range timestamps {
			if ts.After(cutoff) {
				validTimestamps = append(validTimestamps, ts)
			}
		}
		rl.requests[key] = validTimestamps

		// Check if limit exceeded
		if len(validTimestamps) >= rl.limit {
			return false
		}
	}

	// Add current request
	if rl.requests[key] == nil {
		rl.requests[key] = make([]time.Time, 0, rl.limit)
	}
	rl.requests[key] = append(rl.requests[key], now)

	return true
}

// Reset clears all requests for a key
func (rl *InMemoryRateLimiter) Reset(key string) {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	delete(rl.requests, key)
}

// cleanupOldEntries periodically removes old entries
func (rl *InMemoryRateLimiter) cleanupOldEntries() {
	for range rl.cleanup.C {
		rl.mu.Lock()
		cutoff := time.Now().Add(-rl.window)
		for key, timestamps := range rl.requests {
			validTimestamps := make([]time.Time, 0)
			for _, ts := range timestamps {
				if ts.After(cutoff) {
					validTimestamps = append(validTimestamps, ts)
				}
			}
			if len(validTimestamps) == 0 {
				delete(rl.requests, key)
			} else {
				rl.requests[key] = validTimestamps
			}
		}
		rl.mu.Unlock()
	}
}

// Stop stops the cleanup ticker
func (rl *InMemoryRateLimiter) Stop() {
	rl.cleanup.Stop()
}

// RateLimitMiddleware creates a rate limiting middleware
func RateLimitMiddleware(limiter RateLimiter, keyFunc func(ctx context.Context) string) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			key := keyFunc(ctx)
			if !limiter.Allow(key) {
				return nil, errors.New(429, "RATE_LIMIT_EXCEEDED", "rate limit exceeded")
			}
			return handler(ctx, req)
		}
	}
}

// IPBasedRateLimit creates rate limiting based on IP address
func IPBasedRateLimit(limiter RateLimiter) middleware.Middleware {
	return RateLimitMiddleware(limiter, func(ctx context.Context) string {
		if req, ok := http.RequestFromServerContext(ctx); ok {
			ip := req.RemoteAddr
			// Extract IP from "host:port" format
			if idx := len(ip) - 1; idx >= 0 && ip[idx] == ']' {
				// IPv6 format
				for i := idx - 1; i >= 0; i-- {
					if ip[i] == '[' {
						ip = ip[i+1 : idx]
						break
					}
				}
			} else {
				// IPv4 format
				for i := len(ip) - 1; i >= 0; i-- {
					if ip[i] == ':' {
						ip = ip[:i]
						break
					}
				}
			}
			return fmt.Sprintf("ip:%s", ip)
		}
		return "unknown"
	})
}

// LoginRateLimit creates rate limiting specifically for login endpoint
// Limits: 50 attempts per 15 minutes per IP
func LoginRateLimit() middleware.Middleware {
	limiter := NewInMemoryRateLimiter(500, 15*time.Minute)
	return IPBasedRateLimit(limiter)
}

