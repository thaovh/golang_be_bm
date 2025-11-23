# Auth Service Documentation

Hướng dẫn sử dụng Auth Service với JWT authentication.

## Tổng quan

Auth Service cung cấp:
- **Login**: Đăng nhập với email/username và password
- **Register**: Đăng ký user mới và tự động login
- **RefreshToken**: Làm mới access token
- **Logout**: Thu hồi token
- **GetCurrentUser**: Lấy thông tin user hiện tại từ token
- **RevokeAllTokens**: Thu hồi tất cả tokens của user

## Authentication Flow

```
1. Client → POST /api/v1/auth/login (email, password)
2. Server → Validate credentials
3. Server → Generate JWT access token + refresh token
4. Server → Save refresh token to database
5. Server → Return tokens + user info
6. Client → Store tokens
7. Client → Use access token in Authorization header for protected routes
8. When access token expires → Use refresh token to get new access token
```

## API Endpoints

### 1. Login

**POST** `/api/v1/auth/login`

```bash
curl -X POST http://localhost:8000/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "Admin@123"
  }'
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "019ab143-5427-74b2-89e8-fb6f03168236",
  "expires_in": 3600,
  "token_type": "Bearer",
  "user": {
    "id": "019ab143-5427-74b2-89e8-fb6f03168236",
    "email": "admin@example.com",
    "username": "admin",
    "full_name": "Administrator",
    "role": "admin",
    "status": "active"
  }
}
```

### 2. Register

**POST** `/api/v1/auth/register`

```bash
curl -X POST http://localhost:8000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newuser@example.com",
    "username": "newuser",
    "password": "Secure@1234",
    "full_name": "New User"
  }'
```

**Response:** Tương tự Login response

### 3. Refresh Token

**POST** `/api/v1/auth/refresh`

```bash
curl -X POST http://localhost:8000/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "019ab143-5427-74b2-89e8-fb6f03168236"
  }'
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "019ab143-5427-74b2-89e8-fb6f03168236",
  "expires_in": 3600,
  "token_type": "Bearer"
}
```

### 4. Get Current User

**GET** `/api/v1/auth/me`

```bash
curl http://localhost:8000/api/v1/auth/me \
  -H "Authorization: Bearer <access_token>"
```

**Response:**
```json
{
  "user": {
    "id": "019ab143-5427-74b2-89e8-fb6f03168236",
    "email": "admin@example.com",
    "username": "admin",
    "full_name": "Administrator",
    "role": "admin",
    "status": "active"
  }
}
```

### 5. Logout

**POST** `/api/v1/auth/logout`

```bash
curl -X POST http://localhost:8000/api/v1/auth/logout \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "019ab143-5427-74b2-89e8-fb6f03168236"
  }'
```

**Response:**
```json
{
  "success": true
}
```

### 6. Revoke All Tokens

**POST** `/api/v1/auth/revoke-all`

```bash
curl -X POST http://localhost:8000/api/v1/auth/revoke-all \
  -H "Authorization: Bearer <access_token>"
```

**Response:**
```json
{
  "success": true
}
```

## Sử dụng Token

### Trong HTTP Requests

```bash
curl http://localhost:8000/api/v1/users/me \
  -H "Authorization: Bearer <access_token>"
```

### Trong JavaScript/Frontend

```javascript
fetch('http://localhost:8000/api/v1/users/me', {
  headers: {
    'Authorization': `Bearer ${accessToken}`
  }
})
```

## Token Details

### Access Token
- **Type**: JWT (JSON Web Token)
- **Algorithm**: HS256
- **Expiry**: 1 hour (configurable)
- **Contains**: User ID, Email, Role
- **Storage**: Client-side (memory/localStorage)

### Refresh Token
- **Type**: UUID v7
- **Expiry**: 7 days (configurable)
- **Storage**: Database (auth_tokens table)
- **Purpose**: Generate new access tokens

## Security Features

1. ✅ JWT với secret key
2. ✅ Password hashing (bcrypt)
3. ✅ Token expiry
4. ✅ Token revocation (logout)
5. ✅ Refresh token rotation
6. ✅ IP và User-Agent tracking
7. ✅ User status validation (active/inactive)

## Error Responses

### Invalid Credentials
```json
{
  "code": 401,
  "reason": "INVALID_CREDENTIALS",
  "message": "invalid email or password"
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

### Token Invalid
```json
{
  "code": 401,
  "reason": "TOKEN_INVALID",
  "message": "invalid token"
}
```

### Unauthorized
```json
{
  "code": 401,
  "reason": "UNAUTHORIZED",
  "message": "missing or invalid authorization header"
}
```

## Testing

Chạy test script:

```bash
./scripts/test_auth.sh
```

## Configuration

**File: `configs/config.yaml`**

```yaml
auth:
  jwt_secret: "your-secret-key-change-in-production-min-32-chars"
  access_token_expiry: 3600    # 1 hour in seconds
  refresh_token_expiry: 604800 # 7 days in seconds
```

**Lưu ý**: Đổi `jwt_secret` trong production!

## Best Practices

1. **Store tokens securely**: 
   - Access token: Memory (JavaScript) hoặc secure storage
   - Refresh token: HttpOnly cookie (recommended) hoặc secure storage

2. **Token refresh strategy**:
   - Refresh access token trước khi expire (ví dụ: 5 phút trước)
   - Handle token expiry gracefully

3. **Logout**:
   - Luôn revoke refresh token khi logout
   - Clear tokens ở client-side

4. **Security**:
   - Sử dụng HTTPS trong production
   - Rotate JWT secret định kỳ
   - Monitor token usage

## Integration với User Service

Auth service tự động:
- Tạo user mới khi register (reuse UserUsecase)
- Update last login khi login thành công
- Validate user status (active/inactive)

## Middleware Usage

Để protect routes, sử dụng AuthMiddleware:

```go
// In server setup
http.Middleware(
    recovery.Recovery(),
    middleware.AuthMiddleware([]byte(jwtSecret)),
    // Other middlewares
)
```

## Next Steps

1. Implement role-based access control (RBAC)
2. Add rate limiting cho login attempts
3. Implement 2FA (Two-Factor Authentication)
4. Add session management
5. Implement token blacklist (nếu cần)

