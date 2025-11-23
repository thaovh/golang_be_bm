# Hướng dẫn Triển khai Backend Service

## Bước 1: Setup Môi trường Development

### Yêu cầu
- Go 1.23+
- Docker & Docker Compose
- Make
- Protobuf compiler

### Cài đặt Tools

```bash
# Install Kratos CLI
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest

# Install Wire (Dependency Injection)
go install github.com/google/wire/cmd/wire@latest

# Install Protobuf tools
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## Bước 2: Tạo Project Structure

### Option 1: Sử dụng Kratos CLI (Khuyến nghị)

```bash
# Tạo project mới
kratos new backend-service
cd backend-service

# Hoặc clone template này
```

### Option 2: Tạo thủ công

Sử dụng cấu trúc trong blueprint này.

## Bước 3: Setup Dependencies

```bash
# Copy go.mod và cập nhật
go mod tidy

# Generate code từ protobuf
make api
make config

# Generate dependency injection
make wire
```

## Bước 4: Cấu hình

### 4.1. Database

Chỉnh sửa `configs/config.yaml`:

```yaml
data:
  database:
    driver: mysql
    source: root:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
```

### 4.2. Service Discovery

```yaml
registry:
  consul:
    address: 127.0.0.1:8500
    scheme: http
    datacenter: dc1
```

### 4.3. Tracing

```yaml
trace:
  endpoint: localhost:4318
```

## Bước 5: Chạy với Docker Compose

```bash
# Start tất cả services
docker-compose up -d

# Check logs
docker-compose logs -f backend-service

# Stop
docker-compose down
```

## Bước 6: Development Workflow

### Local Development

```bash
# Run service locally (không cần Docker)
make run

# Build binary
make build

# Run tests
make test
```

### Hot Reload (với Air)

```bash
# Install Air
go install github.com/cosmtrek/air@latest

# Run với hot reload
air
```

## Bước 7: Kiến trúc Layers

### 1. API Layer (`api/`)
- Protobuf definitions
- Generate gRPC và HTTP handlers

### 2. Service Layer (`internal/service/`)
- Implement gRPC service interface
- Call business logic

### 3. Business Layer (`internal/biz/`)
- Business logic
- Domain models
- Use cases

### 4. Data Layer (`internal/data/`)
- Database access
- Cache
- External API clients

## Bước 8: Middleware Setup

Thêm middleware trong `cmd/service/main.go`:

```go
httpSrv := http.NewServer(
    http.Address(":8000"),
    http.Middleware(
        recovery.Recovery(),
        logging.Server(logger),
        tracing.Server(),
        metrics.Server(),
        validate.Validator(),
    ),
)
```

## Bước 9: Monitoring

### Prometheus Metrics

Truy cập: http://localhost:9090

### Grafana Dashboard

Truy cập: http://localhost:3000
- Username: admin
- Password: admin

### Tracing

Truy cập Jaeger UI (nếu setup) hoặc xem logs từ OTEL Collector.

## Bước 10: Production Deployment

### Build Docker Image

```bash
docker build -t backend-service:v1.0.0 .
```

### Kubernetes Deployment

Xem thư mục `deployments/k8s/` cho Kubernetes manifests.

### Environment Variables

```bash
export CONSUL_ADDR=consul:8500
export DB_HOST=mysql
export REDIS_ADDR=redis:6379
```

## Best Practices

1. **Error Handling**: Sử dụng Kratos errors
2. **Logging**: Structured logging với context
3. **Metrics**: Expose Prometheus metrics
4. **Tracing**: Distributed tracing với OpenTelemetry
5. **Config**: Hot-reload config không restart
6. **Service Discovery**: Auto-register với Consul
7. **Health Check**: Implement health check endpoint
8. **Graceful Shutdown**: Handle signals properly

## Troubleshooting

### Service không register vào Consul
- Check Consul đang chạy: `docker-compose ps`
- Check config registry trong `config.yaml`

### Database connection failed
- Check MySQL đang chạy
- Verify connection string

### Tracing không hoạt động
- Check OTEL Collector đang chạy
- Verify endpoint trong config

