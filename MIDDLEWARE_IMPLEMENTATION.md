# Middleware Implementation Guide

## Tổng quan

Đã triển khai 2 middleware chính:
1. **Rate Limiting Middleware** - Giới hạn số lượng requests cho login/register
2. **Auth Middleware** - Bảo vệ các routes cần authentication

## 1. Rate Limiting Middleware

### Implementation
- **File**: `internal/middleware/ratelimit.go`
- **Type**: In-memory rate limiter
- **Algorithm**: Sliding window

### Configuration
- **Login/Register**: 5 requests per 15 minutes per IP
- **Storage**: In-memory (có thể mở rộng sang Redis)

### Usage
```go
// Tạo rate limiter
limiter := middleware.NewInMemoryRateLimiter(5, 15*time.Minute)

// Apply cho specific paths
selector.Server(middleware.LoginRateLimit()).
    Match(func(ctx context.Context, operation string) bool {
        // Match logic
    }).Build()
```

### Test Rate Limiting
```bash
# Test rate limiting (sẽ block sau 5 requests)
for i in {1..6}; do
  curl -X POST http://localhost:8000/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d '{"email":"test@test.com","password":"wrong"}'
  sleep 0.5
done
```

**Expected Result**: 
- Requests 1-5: `invalid email or password`
- Request 6: `rate limit exceeded` (429)

## 2. Auth Middleware

### Implementation
- **File**: `internal/middleware/auth.go`
- **Type**: JWT token validation
- **Storage**: Context values

### Protected Routes
Các routes sau yêu cầu authentication:
- `/api/v1/users` (tất cả methods)
- `/api/v1/auth/me`
- `/api/v1/auth/logout`
- `/api/v1/auth/revoke-all`

### Usage
```go
// Tạo auth middleware
authMiddleware := middleware.AuthMiddleware([]byte(jwtSecret))

// Apply cho protected routes
selector.Server(authMiddleware).
    Match(func(ctx context.Context, operation string) bool {
        // Match protected paths
    }).Build()
```

### Test Auth Middleware
```bash
# Test without token (should fail)
curl http://localhost:8000/api/v1/users \
  -H "Authorization: Bearer invalid"

# Expected: {"code":401, "reason":"TOKEN_INVALID", ...}

# Test with valid token
TOKEN="your-access-token"
curl http://localhost:8000/api/v1/users \
  -H "Authorization: Bearer $TOKEN"
```

## 3. Server Configuration

### File: `internal/server/http.go`

```go
func NewHTTPServer(...) *http.Server {
    // Rate limiting
    loginRateLimit := middleware.LoginRateLimit()
    
    // Auth middleware
    authMiddleware := middleware.AuthMiddleware([]byte(authConfig.JwtSecret))
    
    // Protected paths
    protectedPaths := []string{
        "/api/v1/users",
        "/api/v1/auth/me",
        "/api/v1/auth/logout",
        "/api/v1/auth/revoke-all",
    }
    
    // Rate limited paths
    rateLimitedPaths := []string{
        "/api/v1/auth/login",
        "/api/v1/auth/register",
    }
    
    var opts = []http.ServerOption{
        http.Middleware(
            recovery.Recovery(),
            // Rate limiting
            selector.Server(loginRateLimit).
                Match(func(ctx context.Context, operation string) bool {
                    // Match rate limited paths
                }).Build(),
            // Auth middleware
            selector.Server(authMiddleware).
                Match(func(ctx context.Context, operation string) bool {
                    // Match protected paths
                }).Build(),
        ),
    }
    // ...
}
```

## 4. Middleware Chain Order

Middleware được apply theo thứ tự:
1. **Recovery** - Catch panics
2. **Rate Limiting** - Limit requests (login/register only)
3. **Auth** - Validate tokens (protected routes only)
4. **Handler** - Business logic

## 5. Context Values

Auth middleware thêm các giá trị vào context:

```go
// Get user ID
userID, ok := middleware.GetUserIDFromContext(ctx)

// Get user role
role, ok := middleware.GetUserRoleFromContext(ctx)
```

## 6. Error Responses

### Rate Limit Exceeded
```json
{
  "code": 429,
  "reason": "RATE_LIMIT_EXCEEDED",
  "message": "rate limit exceeded"
}
```

### Unauthorized
```json
{
  "code": 401,
  "reason": "UNAUTHORIZED",
  "message": "missing authorization header"
}
```

### Token Invalid
```json
{
  "code": 401,
  "reason": "TOKEN_INVALID",
  "message": "invalid token"
}
```

### Token Expired
```json
{
  "code": 401,
  "reason": "TOKEN_EXPIRED",
  "message": "token has expired"
}
```

## 7. Testing

### Test Script
```bash
# Test auth service
./scripts/test_auth.sh

# Test rate limiting
for i in {1..6}; do
  curl -X POST http://localhost:8000/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d '{"email":"test@test.com","password":"wrong"}'
  sleep 0.5
done

# Test protected routes
curl http://localhost:8000/api/v1/users \
  -H "Authorization: Bearer <token>"
```

## 8. Customization

### Adjust Rate Limits
```go
// In internal/middleware/ratelimit.go
func LoginRateLimit() middleware.Middleware {
    limiter := NewInMemoryRateLimiter(
        5,              // limit
        15*time.Minute, // window
    )
    return IPBasedRateLimit(limiter)
}
```

### Add More Protected Routes
```go
protectedPaths := []string{
    "/api/v1/users",
    "/api/v1/users/{id}",
    "/api/v1/admin",
    // Add more...
}
```

### Redis-based Rate Limiting (Future)
```go
// Replace InMemoryRateLimiter with RedisRateLimiter
limiter := middleware.NewRedisRateLimiter(redisClient, 5, 15*time.Minute)
```

## 9. Best Practices

1. **Rate Limiting**:
   - Adjust limits based on traffic
   - Monitor rate limit hits
   - Consider per-user limits (not just IP)

2. **Auth Middleware**:
   - Always validate tokens
   - Check user status (active/inactive)
   - Log authentication failures

3. **Security**:
   - Use HTTPS in production
   - Rotate JWT secrets regularly
   - Monitor suspicious activity

## 10. Performance

- **Rate Limiter**: O(1) lookup, in-memory
- **Auth Middleware**: O(1) token validation
- **Cleanup**: Background goroutine (every 1 minute)

## Next Steps

1. ✅ Rate limiting implemented
2. ✅ Auth middleware implemented
3. ⏳ Add Redis-based rate limiting (for distributed systems)
4. ⏳ Add per-user rate limiting
5. ⏳ Add rate limit metrics
6. ⏳ Add IP whitelist/blacklist

