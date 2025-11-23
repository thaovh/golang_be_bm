# Kiến trúc Backend Service

## Tổng quan

Backend service được xây dựng dựa trên **Clean Architecture** và **Domain-Driven Design (DDD)** principles, sử dụng Kratos framework.

## Kiến trúc Layers

```
┌─────────────────────────────────────────┐
│         API Layer (api/)                │
│  - Protobuf definitions                 │
│  - gRPC/HTTP interfaces                 │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│      Service Layer (internal/service/)  │
│  - gRPC service implementation         │
│  - HTTP handlers                        │
│  - Request/Response transformation      │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│    Business Layer (internal/biz/)       │
│  - Business logic                       │
│  - Domain models                        │
│  - Use cases                            │
│  - Validation rules                     │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│      Data Layer (internal/data/)        │
│  - Database access (GORM)              │
│  - Cache (Redis)                        │
│  - External API clients                 │
│  - Repository implementations           │
└─────────────────────────────────────────┘
```

## Components

### 1. API Layer

**Mục đích**: Định nghĩa contracts giữa services

**File**: `api/service/v1/service.proto`

```protobuf
service Service {
  rpc GetService (GetServiceRequest) returns (GetServiceResponse) {
    option (google.api.http) = {
      get: "/api/v1/service/{id}",
    };
  }
}
```

**Generate code**:
```bash
make api
```

### 2. Service Layer

**Mục đích**: Implement gRPC/HTTP handlers, transform requests/responses

**Pattern**: 
- Implement generated interface từ protobuf
- Call business logic từ biz layer
- Handle errors và transform

### 3. Business Layer

**Mục đích**: Core business logic, domain rules

**Pattern**:
- Use cases (ServiceUsecase)
- Domain models
- Business validation

**Không phụ thuộc vào**:
- Database
- HTTP/gRPC
- External services

### 4. Data Layer

**Mục đích**: Data access, external integrations

**Pattern**:
- Repository pattern
- Implement interfaces từ biz layer
- Database, Cache, External APIs

## Dependency Flow

```
API → Service → Business → Data
  ↓      ↓         ↓         ↓
  └──────┴─────────┴─────────┘
         (Dependency Injection)
```

**Rule**: Dependencies chỉ đi vào trong, không đi ra ngoài.

## Technology Stack

### Core
- **Framework**: Kratos v2
- **Language**: Go 1.23+
- **Protocols**: HTTP, gRPC

### Data
- **Database**: MySQL (GORM)
- **Cache**: Redis
- **ORM**: GORM

### Infrastructure
- **Service Discovery**: Consul
- **Config**: YAML với hot-reload
- **Container**: Docker

### Observability
- **Metrics**: Prometheus
- **Tracing**: OpenTelemetry
- **Logging**: Structured logging

### Resilience
- **Circuit Breaker**: go-kratos/aegis
- **Retry**: Built-in với backoff
- **Rate Limiting**: BBR algorithm

## Design Patterns

### 1. Dependency Injection (Wire)

```go
// wire.go
func wireApp(...) (*kratos.App, func(), error) {
    panic(wire.Build(
        server.ProviderSet,
        data.ProviderSet,
        biz.ProviderSet,
        service.ProviderSet,
        newApp,
    ))
}
```

### 2. Repository Pattern

```go
// biz/service.go
type ServiceRepo interface {
    Save(context.Context, *Service) (*Service, error)
    FindByID(context.Context, int64) (*Service, error)
}

// data/service.go
type serviceRepo struct {
    data *Data
}
```

### 3. Middleware Chain

```go
httpSrv := http.NewServer(
    http.Middleware(
        recovery.Recovery(),
        logging.Server(logger),
        tracing.Server(),
        metrics.Server(),
    ),
)
```

## Service Discovery Flow

```
1. Service starts
2. Register với Consul
3. Health check endpoint
4. Consul watches service
5. Gateway discovers service
6. Load balance requests
```

## Request Flow

```
Client Request
    ↓
Gateway (Load Balance)
    ↓
Service (HTTP/gRPC)
    ↓
Middleware Chain
    ↓
Service Handler
    ↓
Business Logic
    ↓
Data Layer
    ↓
Database/Cache
    ↓
Response
```

## Configuration Management

### Hot Reload

```go
c := config.New(
    config.WithSource(
        file.NewSource(flagconf),
    ),
)
c.Watch("server", func(key string, value config.Value) {
    // Reload config
})
```

### Environment-based Config

```yaml
# configs/config.yaml (default)
# configs/config.prod.yaml (production)
# configs/config.dev.yaml (development)
```

## Error Handling

### Kratos Errors

```go
import "github.com/go-kratos/kratos/v2/errors"

// Define errors
var ErrServiceNotFound = errors.NotFound("SERVICE_NOT_FOUND", "service not found")

// Return errors
return nil, ErrServiceNotFound
```

### HTTP Status Mapping

- `NotFound` → 404
- `BadRequest` → 400
- `Unauthorized` → 401
- `Internal` → 500

## Testing Strategy

### Unit Tests
- Business logic
- Repository mocks

### Integration Tests
- Database operations
- External API calls

### E2E Tests
- Full request flow
- Service discovery

## Scalability

### Horizontal Scaling
- Stateless services
- Service discovery
- Load balancing

### Vertical Scaling
- Resource limits
- Connection pooling
- Cache optimization

## Security

- TLS/SSL
- Authentication (JWT)
- Authorization (RBAC)
- Input validation
- SQL injection prevention (GORM)

## Performance

- Connection pooling
- Cache strategy
- Database indexing
- Async processing
- Circuit breaker

