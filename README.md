# Backend Service Blueprint

Kiến trúc backend service dựa trên các công nghệ từ gateway project.

## Công nghệ Stack

- **Framework**: Kratos v2 (Go microservices framework)
- **Protocols**: HTTP, gRPC
- **Service Discovery**: Consul
- **Config**: YAML với hot-reload
- **Observability**: 
  - Metrics: Prometheus
  - Tracing: OpenTelemetry
  - Logging: Structured logging
- **Resilience**: Circuit Breaker, Retry, Timeout
- **Container**: Docker + Docker Compose

## Kiến trúc

```
┌─────────────┐
│   Gateway   │
└──────┬──────┘
       │
       ▼
┌─────────────────────────────────────┐
│         Backend Services            │
├─────────────────────────────────────┤
│  ┌──────────┐  ┌──────────┐        │
│  │ Service1 │  │ Service2 │  ...   │
│  │ (User)   │  │ (Order)  │        │
│  └────┬─────┘  └────┬─────┘        │
│       │             │               │
│  ┌────▼─────────────▼─────┐         │
│  │   Shared Components    │         │
│  │  - Database            │         │
│  │  - Cache (Redis)       │         │
│  │  - Message Queue       │         │
│  └────────────────────────┘         │
└─────────────────────────────────────┘
```

## Cấu trúc Project

```
backend-service/
├── api/                    # Protobuf definitions
│   └── service/
│       └── v1/
│           ├── service.proto
│           └── service.pb.go
├── cmd/                    # Application entrypoints
│   └── service/
│       ├── main.go
│       └── wire.go         # Dependency injection
├── internal/              # Private application code
│   ├── biz/               # Business logic
│   ├── data/              # Data access layer
│   ├── service/           # Service implementation
│   └── conf/              # Configuration
├── configs/               # Configuration files
│   └── config.yaml
├── deployments/          # Deployment configs
│   ├── docker-compose.yaml
│   └── k8s/
├── Dockerfile
├── Makefile
└── go.mod
```

## Quick Start

### 1. Tạo project với Kratos CLI

```bash
# Install Kratos CLI
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest

# Tạo project
kratos new backend-service
cd backend-service
```

### 2. Setup dependencies

```bash
go mod tidy
```

### 3. Chạy với Docker Compose

```bash
docker-compose up -d
```

## Development

```bash
# Generate code từ protobuf
make api

# Run service
make run

# Build
make build
```

## 📚 Documentation

Tất cả documentation được lưu trong thư mục [`docs/`](docs/).

### Getting Started
- **[Quick Start Guide](docs/QUICK_START.md)** - Hướng dẫn setup project từ đầu
- **[Architecture](docs/ARCHITECTURE.md)** - Kiến trúc tổng quan của hệ thống
- **[Deployment Guide](docs/DEPLOYMENT_GUIDE.md)** - Hướng dẫn triển khai

### Adding New Services
- **[Add New Service Guide](docs/ADD_NEW_SERVICE_GUIDE.md)** - Hướng dẫn chi tiết thêm service mới (812 dòng)
- **[Quick Reference](docs/QUICK_REFERENCE.md)** - Quick reference card cho việc thêm service

### Services Documentation
- **[Auth Service](docs/AUTH_SERVICE.md)** - Authentication & Authorization
- **[API Testing](docs/API_TESTING.md)** - Hướng dẫn test API
- **[Middleware Implementation](docs/MIDDLEWARE_IMPLEMENTATION.md)** - Rate limiting & Auth middleware

### Database
- **[Migrations](migrations/README.md)** - Database migrations guide

