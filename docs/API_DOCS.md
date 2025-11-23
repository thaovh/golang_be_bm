# API Documentation Guide

## Swagger UI

Project đã được tích hợp Swagger UI để xem và test API documentation.

### Truy cập Swagger UI

Sau khi start service, truy cập:

```
http://localhost:8000/api/docs
```

### Endpoints

- **Swagger UI**: `GET /api/docs` - Interactive API documentation
- **OpenAPI Spec**: `GET /api/openapi.yaml` - OpenAPI 3.0 specification file

### Features

1. **Interactive API Testing**
   - Xem tất cả endpoints
   - Test API trực tiếp từ browser
   - Xem request/response schemas

2. **Auto-generated từ Protobuf**
   - OpenAPI spec được generate tự động từ `.proto` files
   - Chạy `make api` để regenerate khi có thay đổi

3. **Authentication Support**
   - Có thể test protected endpoints
   - Sử dụng "Authorize" button để thêm JWT token

### Generate OpenAPI Spec

```bash
# Generate OpenAPI spec từ protobuf
make api

# File openapi.yaml sẽ được tạo/update ở root directory
```

### Customization

File Swagger UI được embed trong `internal/server/swagger-ui/` và được serve qua `internal/server/swagger.go`.

Để customize:
1. Edit `internal/server/swagger.go` - function `createSwaggerUIHTML()`
2. Update Swagger UI files trong `internal/server/swagger-ui/` nếu cần

### Troubleshooting

**Swagger UI không load:**
- Kiểm tra service đang chạy: `curl http://localhost:8000/api/openapi.yaml`
- Kiểm tra file `openapi.yaml` tồn tại ở root directory
- Xem logs: `tail -f /tmp/server.log`

**OpenAPI spec không update:**
- Chạy `make api` để regenerate
- Restart service

### Example Usage

1. Start service:
```bash
go run cmd/server/main.go cmd/server/wire_gen.go
```

2. Open browser:
```
http://localhost:8000/api/docs
```

3. Test API:
   - Click vào endpoint để expand
   - Click "Try it out"
   - Fill request body
   - Click "Execute"
   - Xem response

### Protected Endpoints

Để test protected endpoints:
1. Login trước để lấy token: `POST /api/v1/auth/login`
2. Copy `access_token` từ response
3. Click "Authorize" button ở top right
4. Enter: `Bearer <your-token>`
5. Click "Authorize"
6. Test protected endpoints

