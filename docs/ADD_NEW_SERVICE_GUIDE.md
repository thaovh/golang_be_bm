# Hướng Dẫn Thêm Service Mới

## 📋 Tổng Quan

Khi thêm một service mới vào project, cần tuân theo **Clean Architecture** và **CQRS Pattern** với 4 layers:

```
API Layer (api/) 
    ↓
Service Layer (internal/service/)
    ↓
Business Layer (internal/biz/)
    ↓
Data Layer (internal/data/)
```

## 🎯 Quy Trình 7 Bước

### Bước 1: Tạo Domain Model (Business Layer)

**File**: `internal/biz/{service_name}.go`

```go
package biz

import (
    "context"
    "time"
    "github.com/gofrs/uuid/v5"
    "github.com/go-kratos/kratos/v2/errors"
    "github.com/go-kratos/kratos/v2/log"
)

// {Service}Entity - Domain model (embeds BaseEntity)
type {Service}Entity struct {
    BaseEntity
    
    // Your fields here
    Name        string `gorm:"type:varchar(255);not null" json:"name"`
    Description string `gorm:"type:text" json:"description"`
}

// {Service}CommandRepo - Write operations interface
type {Service}CommandRepo interface {
    Save(context.Context, *{Service}Entity) (*{Service}Entity, error)
    Update(context.Context, *{Service}Entity) (*{Service}Entity, error)
    Delete(context.Context, uuid.UUID) error
}

// {Service}QueryRepo - Read operations interface
type {Service}QueryRepo interface {
    FindByID(context.Context, uuid.UUID) (*{Service}Entity, error)
    List(context.Context, *{Service}ListFilter) ([]*{Service}Entity, int64, error)
    Count(context.Context, *{Service}ListFilter) (int64, error)
}

// {Service}ListFilter - Filter for listing
type {Service}ListFilter struct {
    Page     int32
    PageSize int32
    Search   string
    Status   string
}

// {Service}Usecase - Business logic
type {Service}Usecase struct {
    commandRepo {Service}CommandRepo
    queryRepo   {Service}QueryRepo
    log         *log.Helper
}

func New{Service}Usecase(
    commandRepo {Service}CommandRepo,
    queryRepo {Service}QueryRepo,
    logger log.Logger,
) *{Service}Usecase {
    return &{Service}Usecase{
        commandRepo: commandRepo,
        queryRepo:   queryRepo,
        log:         log.NewHelper(logger),
    }
}

// Business methods
func (uc *{Service}Usecase) Create{Service}(ctx context.Context, entity *{Service}Entity) (*{Service}Entity, error) {
    uc.log.WithContext(ctx).Infof("Create{Service}: %s", entity.Name)
    return uc.commandRepo.Save(ctx, entity)
}

func (uc *{Service}Usecase) Get{Service}(ctx context.Context, id uuid.UUID) (*{Service}Entity, error) {
    return uc.queryRepo.FindByID(ctx, id)
}
```

**Checklist:**
- ✅ Entity embeds `BaseEntity`
- ✅ Tách `CommandRepo` và `QueryRepo`
- ✅ Tạo `Usecase` với business logic
- ✅ Sử dụng `uuid.UUID` cho ID
- ✅ Logging với context

---

### Bước 2: Implement Repositories (Data Layer)

**File 1**: `internal/data/{service_name}_command.go` (Write operations)

```go
package data

import (
    "context"
    "github.com/go-kratos/kratos-layout/internal/biz"
    "github.com/gofrs/uuid/v5"
    "github.com/go-kratos/kratos/v2/log"
    "gorm.io/gorm"
)

type {service}CommandRepo struct {
    data *Data
    log  *log.Helper
}

func New{Service}CommandRepo(data *Data, logger log.Logger) biz.{Service}CommandRepo {
    return &{service}CommandRepo{
        data: data,
        log:  log.NewHelper(logger),
    }
}

func (r *{service}CommandRepo) Save(ctx context.Context, entity *biz.{Service}Entity) (*biz.{Service}Entity, error) {
    db := r.data.GetWriteDB() // Use write DB for writes
    if err := db.WithContext(ctx).Save(entity).Error; err != nil {
        r.log.WithContext(ctx).Errorf("Failed to save: %v", err)
        return nil, err
    }
    return entity, nil
}

func (r *{service}CommandRepo) Update(ctx context.Context, entity *biz.{Service}Entity) (*biz.{Service}Entity, error) {
    db := r.data.GetWriteDB()
    if err := db.WithContext(ctx).Model(entity).Updates(entity).Error; err != nil {
        r.log.WithContext(ctx).Errorf("Failed to update: %v", err)
        return nil, err
    }
    return entity, nil
}

func (r *{service}CommandRepo) Delete(ctx context.Context, id uuid.UUID) error {
    db := r.data.GetWriteDB()
    if err := db.WithContext(ctx).Delete(&biz.{Service}Entity{}, id).Error; err != nil {
        r.log.WithContext(ctx).Errorf("Failed to delete: %v", err)
        return nil, err
    }
    return nil
}
```

**File 2**: `internal/data/{service_name}_query.go` (Read operations)

```go
package data

import (
    "context"
    "github.com/go-kratos/kratos-layout/internal/biz"
    "github.com/gofrs/uuid/v5"
    "github.com/go-kratos/kratos/v2/log"
    "gorm.io/gorm"
)

type {service}QueryRepo struct {
    data *Data
    log  *log.Helper
}

func New{Service}QueryRepo(data *Data, logger log.Logger) biz.{Service}QueryRepo {
    return &{service}QueryRepo{
        data: data,
        log:  log.NewHelper(logger),
    }
}

func (r *{service}QueryRepo) FindByID(ctx context.Context, id uuid.UUID) (*biz.{Service}Entity, error) {
    db := r.data.GetReadDB() // Use read DB for reads
    var entity biz.{Service}Entity
    if err := db.WithContext(ctx).Where("id = ?", id).First(&entity).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, nil
        }
        r.log.WithContext(ctx).Errorf("Failed to find by ID: %v", err)
        return nil, err
    }
    return &entity, nil
}

func (r *{service}QueryRepo) List(ctx context.Context, filter *biz.{Service}ListFilter) ([]*biz.{Service}Entity, int64, error) {
    db := r.data.GetReadDB()
    var entities []*biz.{Service}Entity
    var total int64
    
    query := db.WithContext(ctx).Model(&biz.{Service}Entity{})
    
    // Apply filters
    if filter.Search != "" {
        query = query.Where("name LIKE ?", "%"+filter.Search+"%")
    }
    if filter.Status != "" {
        query = query.Where("status = ?", filter.Status)
    }
    
    // Count total
    if err := query.Count(&total).Error; err != nil {
        return nil, 0, err
    }
    
    // Pagination
    offset := (filter.Page - 1) * filter.PageSize
    if err := query.Offset(int(offset)).Limit(int(filter.PageSize)).Find(&entities).Error; err != nil {
        r.log.WithContext(ctx).Errorf("Failed to list: %v", err)
        return nil, 0, err
    }
    
    return entities, total, nil
}

func (r *{service}QueryRepo) Count(ctx context.Context, filter *biz.{Service}ListFilter) (int64, error) {
    db := r.data.GetReadDB()
    var count int64
    
    query := db.WithContext(ctx).Model(&biz.{Service}Entity{})
    
    if filter.Search != "" {
        query = query.Where("name LIKE ?", "%"+filter.Search+"%")
    }
    if filter.Status != "" {
        query = query.Where("status = ?", filter.Status)
    }
    
    if err := query.Count(&count).Error; err != nil {
        return 0, err
    }
    
    return count, nil
}
```

**Checklist:**
- ✅ Tách `_command.go` (write) và `_query.go` (read)
- ✅ Command repo dùng `GetWriteDB()`
- ✅ Query repo dùng `GetReadDB()`
- ✅ Handle `gorm.ErrRecordNotFound`
- ✅ Logging errors

---

### Bước 3: Tạo Protobuf Definitions (API Layer)

**File**: `api/{service_name}/v1/{service_name}.proto`

```protobuf
syntax = "proto3";

package {service_name}.v1;

import "google/api/annotations.proto";

option go_package = "github.com/go-kratos/kratos-layout/api/{service_name}/v1;v1";

service {Service}Service {
  // Commands
  rpc Create{Service} (Create{Service}Request) returns (Create{Service}Response) {
    option (google.api.http) = {
      post: "/api/v1/{service_name}s"
      body: "*"
    };
  }
  
  rpc Update{Service} (Update{Service}Request) returns (Update{Service}Response) {
    option (google.api.http) = {
      put: "/api/v1/{service_name}s/{id}"
      body: "*"
    };
  }
  
  rpc Delete{Service} (Delete{Service}Request) returns (Delete{Service}Response) {
    option (google.api.http) = {
      delete: "/api/v1/{service_name}s/{id}"
    };
  }
  
  // Queries
  rpc Get{Service} (Get{Service}Request) returns (Get{Service}Response) {
    option (google.api.http) = {
      get: "/api/v1/{service_name}s/{id}"
    };
  }
  
  rpc List{Service}s (List{Service}sRequest) returns (List{Service}sResponse) {
    option (google.api.http) = {
      get: "/api/v1/{service_name}s"
    };
  }
}

// Request/Response messages
message Create{Service}Request {
  string name = 1;
  string description = 2;
}

message Create{Service}Response {
  {Service} {service} = 1;
}

message Update{Service}Request {
  string id = 1;
  string name = 2;
  string description = 3;
}

message Update{Service}Response {
  {Service} {service} = 1;
}

message Delete{Service}Request {
  string id = 1;
}

message Delete{Service}Response {
  bool success = 1;
}

message Get{Service}Request {
  string id = 1;
}

message Get{Service}Response {
  {Service} {service} = 1;
}

message List{Service}sRequest {
  int32 page = 1;
  int32 page_size = 2;
  string search = 3;
  string status = 4;
}

message List{Service}sResponse {
  repeated {Service} {service}s = 1;
  int64 total = 2;
}

message {Service} {
  string id = 1;
  string name = 2;
  string description = 3;
  string status = 4;
  string created_at = 5;
  string updated_at = 6;
}
```

**File**: `api/{service_name}/v1/error_reason.proto`

```protobuf
syntax = "proto3";

package {service_name}.v1;

option go_package = "github.com/go-kratos/kratos-layout/api/{service_name}/v1;v1";

enum ErrorReason {
  {SERVICE}_UNSPECIFIED = 0;
  {SERVICE}_NOT_FOUND = 1;
  {SERVICE}_ALREADY_EXISTS = 2;
}
```

**Checklist:**
- ✅ Tạo folder `api/{service_name}/v1/`
- ✅ Define service với HTTP annotations
- ✅ Tách Commands và Queries
- ✅ Define request/response messages
- ✅ Tạo error_reason.proto

---

### Bước 4: Generate Code từ Protobuf

```bash
# Generate protobuf code
make api

# Hoặc manually:
protoc --proto_path=./api \
       --proto_path=./third_party \
       --go_out=paths=source_relative:./api \
       --go-http_out=paths=source_relative:./api \
       --go-grpc_out=paths=source_relative:./api \
       api/{service_name}/v1/*.proto
```

**Checklist:**
- ✅ Run `make api` hoặc protoc command
- ✅ Verify generated files trong `api/{service_name}/v1/`

---

### Bước 5: Implement Service Layer

**File**: `internal/service/{service_name}.go`

```go
package service

import (
    "context"
    v1 "github.com/go-kratos/kratos-layout/api/{service_name}/v1"
    "github.com/go-kratos/kratos-layout/internal/biz"
    "github.com/gofrs/uuid/v5"
    "github.com/go-kratos/kratos/v2/errors"
)

type {Service}Service struct {
    v1.Unimplemented{Service}ServiceServer
    
    uc *biz.{Service}Usecase
}

func New{Service}Service(uc *biz.{Service}Usecase) *{Service}Service {
    return &{Service}Service{uc: uc}
}

func (s *{Service}Service) Create{Service}(ctx context.Context, req *v1.Create{Service}Request) (*v1.Create{Service}Response, error) {
    entity := &biz.{Service}Entity{
        Name:        req.Name,
        Description: req.Description,
    }
    
    created, err := s.uc.Create{Service}(ctx, entity)
    if err != nil {
        return nil, err
    }
    
    return &v1.Create{Service}Response{
        {Service}: toProto{Service}(created),
    }, nil
}

func (s *{Service}Service) Get{Service}(ctx context.Context, req *v1.Get{Service}Request) (*v1.Get{Service}Response, error) {
    id, err := uuid.FromString(req.Id)
    if err != nil {
        return nil, errors.BadRequest("INVALID_ID", "invalid id format")
    }
    
    entity, err := s.uc.Get{Service}(ctx, id)
    if err != nil {
        return nil, err
    }
    
    if entity == nil {
        return nil, errors.NotFound("{SERVICE}_NOT_FOUND", "{service} not found")
    }
    
    return &v1.Get{Service}Response{
        {Service}: toProto{Service}(entity),
    }, nil
}

// Helper function
func toProto{Service}(entity *biz.{Service}Entity) *v1.{Service} {
    return &v1.{Service}{
        Id:          entity.ID.String(),
        Name:        entity.Name,
        Description: entity.Description,
        Status:      entity.Status,
        CreatedAt:   entity.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
        UpdatedAt:   entity.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
    }
}
```

**Checklist:**
- ✅ Implement tất cả methods từ proto
- ✅ Convert request → domain entity
- ✅ Convert domain entity → proto response
- ✅ Handle errors properly
- ✅ Validate input (UUID, etc.)

---

### Bước 6: Update Wire Providers

**File 1**: `internal/data/data.go`

```go
var ProviderSet = wire.NewSet(
    NewData,
    // ... existing repos
    New{Service}CommandRepo,  // Add this
    New{Service}QueryRepo,     // Add this
)
```

**File 2**: `internal/biz/biz.go`

```go
var ProviderSet = wire.NewSet(
    // ... existing usecases
    New{Service}Usecase,  // Add this
)
```

**File 3**: `internal/service/service.go`

```go
var ProviderSet = wire.NewSet(
    // ... existing services
    New{Service}Service,  // Add this
)
```

**File 4**: `internal/server/http.go`

```go
func NewHTTPServer(..., {service} *service.{Service}Service, ...) *http.Server {
    // ...
    {service_name}v1.Register{Service}ServiceHTTPServer(srv, {service})
    return srv
}
```

**File 5**: `internal/server/grpc.go`

```go
func NewGRPCServer(..., {service} *service.{Service}Service, ...) *grpc.Server {
    // ...
    {service_name}v1.Register{Service}ServiceServer(srv, {service})
    return srv
}
```

**File 6**: `cmd/server/wire.go`

```go
// Wire sẽ tự động generate, chỉ cần đảm bảo ProviderSets đã được update
```

**Generate Wire:**
```bash
go generate ./cmd/server/...
```

**Checklist:**
- ✅ Add repos vào `data.ProviderSet`
- ✅ Add usecase vào `biz.ProviderSet`
- ✅ Add service vào `service.ProviderSet`
- ✅ Register service trong HTTP server
- ✅ Register service trong gRPC server
- ✅ Run `go generate` để update wire

---

### Bước 7: Tạo Database Migration

**File**: `migrations/XXX_create_{service_name}s_table.sql`

```sql
-- Migration: Create {service_name}s table
-- Created: YYYY-MM-DD

CREATE TABLE IF NOT EXISTS {service_name}s (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    
    -- Audit fields
    created_by UUID NULL,
    updated_by UUID NULL,
    
    -- Optimistic locking
    version INTEGER NOT NULL DEFAULT 1,
    
    -- Status
    status VARCHAR(20) NOT NULL DEFAULT 'active',
    
    -- Your fields
    name VARCHAR(255) NOT NULL,
    description TEXT,
    
    -- Indexes
    CREATE INDEX IF NOT EXISTS idx_{service_name}s_deleted_at ON {service_name}s(deleted_at);
    CREATE INDEX IF NOT EXISTS idx_{service_name}s_status ON {service_name}s(status);
    CREATE INDEX IF NOT EXISTS idx_{service_name}s_name ON {service_name}s(name);
);

-- Trigger for updated_at
CREATE TRIGGER update_{service_name}s_updated_at BEFORE UPDATE ON {service_name}s
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Comments
COMMENT ON TABLE {service_name}s IS 'Stores {service_name} information';
```

**Run Migration:**
```bash
DB_PASSWORD=your_password ./scripts/migrate.sh migrations/XXX_create_{service_name}s_table.sql
```

**Checklist:**
- ✅ Tạo migration file
- ✅ Include BaseEntity fields
- ✅ Add indexes cho common queries
- ✅ Add trigger cho updated_at
- ✅ Run migration

---

## 📝 Checklist Tổng Hợp

### Domain Layer (biz/)
- [ ] Tạo `{service_name}.go` với Entity
- [ ] Entity embeds `BaseEntity`
- [ ] Define `CommandRepo` interface
- [ ] Define `QueryRepo` interface
- [ ] Create `Usecase` với business logic
- [ ] Add vào `biz.ProviderSet`

### Data Layer (data/)
- [ ] Tạo `{service_name}_command.go` (write)
- [ ] Tạo `{service_name}_query.go` (read)
- [ ] Use `GetWriteDB()` cho commands
- [ ] Use `GetReadDB()` cho queries
- [ ] Add vào `data.ProviderSet`

### API Layer (api/)
- [ ] Tạo folder `api/{service_name}/v1/`
- [ ] Tạo `{service_name}.proto`
- [ ] Tạo `error_reason.proto`
- [ ] Run `make api` để generate code

### Service Layer (service/)
- [ ] Tạo `{service_name}.go`
- [ ] Implement tất cả proto methods
- [ ] Add conversion functions
- [ ] Add vào `service.ProviderSet`

### Server Registration
- [ ] Register trong `http.go`
- [ ] Register trong `grpc.go`
- [ ] Update imports

### Database
- [ ] Tạo migration file
- [ ] Run migration

### Wire
- [ ] Update all ProviderSets
- [ ] Run `go generate ./cmd/server/...`
- [ ] Verify `wire_gen.go` updated

### Testing
- [ ] Test với curl/Postman
- [ ] Verify CRUD operations
- [ ] Test error cases

---

## 🎯 Naming Conventions

### Files
- Domain: `internal/biz/{service_name}.go`
- Command Repo: `internal/data/{service_name}_command.go`
- Query Repo: `internal/data/{service_name}_query.go`
- Service: `internal/service/{service_name}.go`
- Proto: `api/{service_name}/v1/{service_name}.proto`

### Types
- Entity: `{Service}Entity` (PascalCase)
- Repo Interface: `{Service}CommandRepo`, `{Service}QueryRepo`
- Usecase: `{Service}Usecase`
- Service: `{Service}Service`
- Proto Service: `{Service}Service`

### Functions
- Constructor: `New{Service}Usecase`, `New{Service}Service`
- Repo Constructor: `New{Service}CommandRepo`, `New{Service}QueryRepo`

---

## 🔍 Ví Dụ: Thêm Product Service

### 1. Domain Model
```go
// internal/biz/product.go
type Product struct {
    BaseEntity
    Name        string  `gorm:"type:varchar(255);not null"`
    Price       float64 `gorm:"type:decimal(10,2)"`
    Description string  `gorm:"type:text"`
}
```

### 2. Repositories
```go
// internal/data/product_command.go
func NewProductCommandRepo(...) biz.ProductCommandRepo { ... }

// internal/data/product_query.go
func NewProductQueryRepo(...) biz.ProductQueryRepo { ... }
```

### 3. Proto
```protobuf
// api/product/v1/product.proto
service ProductService {
  rpc CreateProduct (CreateProductRequest) returns (CreateProductResponse) { ... }
}
```

### 4. Service
```go
// internal/service/product.go
func NewProductService(uc *biz.ProductUsecase) *ProductService { ... }
```

### 5. Wire
```go
// Add to ProviderSets
data.ProviderSet: NewProductCommandRepo, NewProductQueryRepo
biz.ProviderSet: NewProductUsecase
service.ProviderSet: NewProductService
```

---

## ⚠️ Common Mistakes

1. **Quên tách Command/Query Repo**
   - ❌ Dùng 1 repo cho cả read và write
   - ✅ Tách `_command.go` và `_query.go`

2. **Sai database connection**
   - ❌ Dùng `readDB` cho write operations
   - ✅ Command → `GetWriteDB()`, Query → `GetReadDB()`

3. **Quên update Wire**
   - ❌ Chỉ tạo code, quên add vào ProviderSet
   - ✅ Update tất cả ProviderSets và run `go generate`

4. **Quên register service**
   - ❌ Service không accessible
   - ✅ Register trong `http.go` và `grpc.go`

5. **Quên migration**
   - ❌ Entity không có table
   - ✅ Tạo và run migration

---

## 🚀 Quick Start Template

```bash
# 1. Create domain model
touch internal/biz/product.go

# 2. Create repositories
touch internal/data/product_command.go
touch internal/data/product_query.go

# 3. Create proto
mkdir -p api/product/v1
touch api/product/v1/product.proto
touch api/product/v1/error_reason.proto

# 4. Generate proto code
make api

# 5. Create service
touch internal/service/product.go

# 6. Update ProviderSets
# Edit: data/data.go, biz/biz.go, service/service.go

# 7. Register in servers
# Edit: server/http.go, server/grpc.go

# 8. Generate wire
go generate ./cmd/server/...

# 9. Create migration
touch migrations/XXX_create_products_table.sql

# 10. Run migration
DB_PASSWORD=xxx ./scripts/migrate.sh migrations/XXX_create_products_table.sql
```

---

## 📚 References

- Clean Architecture: [ARCHITECTURE.md](ARCHITECTURE.md)
- CQRS Pattern: See User/Auth service implementations
- BaseEntity: `../internal/biz/base.go`
- Wire DI: `../cmd/server/wire.go`

---

**Lưu ý**: Thay `{service_name}` và `{Service}` bằng tên service thực tế của bạn!

