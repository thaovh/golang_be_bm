# API Testing Guide

Hướng dẫn test User Service API endpoints.

## Prerequisites

1. Service đang chạy trên `http://localhost:8000`
2. Database migration đã được chạy
3. Có `curl` và `jq` (optional, để format JSON)

## API Endpoints

### 1. Create User

**POST** `/api/v1/users`

```bash
curl -X POST http://localhost:8000/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "johndoe",
    "password": "Secure@1234",
    "full_name": "John Doe",
    "gender": "male",
    "role": "user"
  }'
```

**Password Requirements:**
- Minimum 8 characters
- At least 1 uppercase letter
- At least 1 lowercase letter
- At least 1 digit
- At least 1 special character

**Response:**
```json
{
  "user": {
    "id": "018f1234-5678-9abc-def0-123456789abc",
    "email": "user@example.com",
    "username": "johndoe",
    "full_name": "John Doe",
    "gender": "male",
    "role": "user",
    "status": "active",
    "created_at": "2025-11-23T10:00:00Z",
    "updated_at": "2025-11-23T10:00:00Z"
  }
}
```

### 2. Get User by ID

**GET** `/api/v1/users/{id}`

```bash
curl http://localhost:8000/api/v1/users/018f1234-5678-9abc-def0-123456789abc
```

### 3. Get User by Email

**GET** `/api/v1/users/email/{email}`

```bash
curl http://localhost:8000/api/v1/users/email/user@example.com
```

### 4. Get User by Username

**GET** `/api/v1/users/username/{username}`

```bash
curl http://localhost:8000/api/v1/users/username/johndoe
```

### 5. List Users

**GET** `/api/v1/users?page=1&page_size=10&search=john&role=user&status=active`

```bash
curl "http://localhost:8000/api/v1/users?page=1&page_size=10"
```

**Query Parameters:**
- `page`: Page number (default: 1)
- `page_size`: Items per page (default: 10, max: 100)
- `search`: Search in email, username, full_name
- `role`: Filter by role (user, admin, moderator)
- `status`: Filter by status (active, inactive, archived)

### 6. Update User

**PUT** `/api/v1/users/{id}`

```bash
curl -X PUT http://localhost:8000/api/v1/users/018f1234-5678-9abc-def0-123456789abc \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "John Updated",
    "gender": "female",
    "date_of_birth": "1990-01-01"
  }'
```

### 7. Change Password

**POST** `/api/v1/users/{id}/change-password`

```bash
curl -X POST http://localhost:8000/api/v1/users/018f1234-5678-9abc-def0-123456789abc/change-password \
  -H "Content-Type: application/json" \
  -d '{
    "new_password": "NewSecure@1234"
  }'
```

### 8. Delete User

**DELETE** `/api/v1/users/{id}`

```bash
curl -X DELETE http://localhost:8000/api/v1/users/018f1234-5678-9abc-def0-123456789abc
```

## Validation Examples

### Invalid Email
```bash
curl -X POST http://localhost:8000/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "invalid-email",
    "username": "test",
    "password": "Test@1234"
  }'
```

**Expected Error:**
```json
{
  "code": 400,
  "reason": "INVALID_EMAIL",
  "message": "invalid email format"
}
```

### Weak Password
```bash
curl -X POST http://localhost:8000/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "username": "test",
    "password": "weak"
  }'
```

**Expected Error:**
```json
{
  "code": 400,
  "reason": "INVALID_PASSWORD",
  "message": "password must be at least 8 characters long"
}
```

### Invalid Username
```bash
curl -X POST http://localhost:8000/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "username": "ab",
    "password": "Test@1234"
  }'
```

**Expected Error:**
```json
{
  "code": 400,
  "reason": "INVALID_USERNAME",
  "message": "username must be at least 3 characters long"
}
```

## Using Test Script

Chạy script tự động để test tất cả endpoints:

```bash
./scripts/test_api.sh
```

Hoặc với custom base URL:

```bash
BASE_URL=http://localhost:8000 ./scripts/test_api.sh
```

## Testing with Postman/Insomnia

Import các requests vào Postman hoặc Insomnia:

1. Base URL: `http://localhost:8000`
2. Headers: `Content-Type: application/json`
3. Sử dụng các JSON payloads từ examples trên

## Notes

- Passwords được hash tự động bằng bcrypt
- User ID là UUID v7 (time-ordered)
- Soft delete được sử dụng (DeletedAt field)
- Tất cả timestamps là UTC

