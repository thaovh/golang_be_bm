# Quick Start Guide

## Bước 1: Tạo Project

### Option A: Sử dụng Kratos CLI (Khuyến nghị)

```bash
# Install Kratos CLI
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest

# Tạo project
kratos new backend-service
cd backend-service
```

### Option B: Copy từ Blueprint

```bash
# Copy thư mục backend-blueprint
cp -r backend-blueprint backend-service
cd backend-service
```

## Bước 2: Setup Dependencies

```bash
# Install tools
make init

# Download dependencies
go mod tidy
```

## Bước 3: Generate Code

```bash
# Generate từ protobuf
make api
make config

# Generate dependency injection (Wire)
make wire
```

## Bước 4: Cấu hình

Chỉnh sửa `configs/config.yaml`:

```yaml
server:
  http:
    addr: 0.0.0.0:8000
  grpc:
    addr: 0.0.0.0:9000

data:
  database:
    source: root:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
  redis:
    addr: 127.0.0.1:6379
```

## Bước 5: Chạy Services

### Với Docker Compose (Khuyến nghị)

```bash
# Start tất cả services (MySQL, Redis, Consul, etc.)
docker-compose up -d

# Check logs
docker-compose logs -f backend-service
```

### Chạy Local (Development)

```bash
# Cần MySQL và Redis đang chạy
make run
```

## Bước 6: Test Service

### Test HTTP

```bash
curl http://localhost:8000/api/v1/service/1
```

### Test gRPC

```bash
# Install grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# List services
grpcurl -plaintext localhost:9000 list

# Call service
grpcurl -plaintext -d '{"id": 1}' localhost:9000 api.service.v1.Service/GetService
```

## Bước 7: Monitoring

- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000 (admin/admin)
- **Consul UI**: http://localhost:8500

## Next Steps

1. Đọc [ARCHITECTURE.md](./ARCHITECTURE.md) để hiểu kiến trúc
2. Đọc [DEPLOYMENT_GUIDE.md](./DEPLOYMENT_GUIDE.md) để triển khai
3. Xem examples trong gateway project để học patterns

## Troubleshooting

### Lỗi: "cannot find package"
```bash
go mod tidy
make wire
```

### Lỗi: "protoc not found"
```bash
# Install protoc
# macOS: brew install protobuf
# Linux: apt-get install protobuf-compiler
```

### Service không register vào Consul
- Check Consul đang chạy: `docker-compose ps`
- Check config trong `configs/config.yaml`

