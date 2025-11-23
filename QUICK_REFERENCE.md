# Quick Reference: Thêm Service Mới

## 🎯 7 Bước Nhanh

```
1. Domain Model (biz/{name}.go)
   └─ Entity + CommandRepo + QueryRepo + Usecase

2. Repositories (data/{name}_command.go + {name}_query.go)
   └─ Command: GetWriteDB() | Query: GetReadDB()

3. Protobuf (api/{name}/v1/{name}.proto)
   └─ Service + Messages + ErrorReason

4. Generate Code
   └─ make api

5. Service Layer (service/{name}.go)
   └─ Implement proto methods + conversions

6. Wire Setup
   └─ Update ProviderSets + go generate

7. Database Migration
   └─ migrations/XXX_create_{name}s_table.sql
```

## 📁 File Structure

```
internal/
├── biz/
│   └── {name}.go              # Domain model + usecase
├── data/
│   ├── {name}_command.go      # Write operations
│   └── {name}_query.go        # Read operations
└── service/
    └── {name}.go              # API service

api/
└── {name}/v1/
    ├── {name}.proto           # API definition
    └── error_reason.proto     # Error codes

migrations/
└── XXX_create_{name}s_table.sql
```

## 🔧 ProviderSets Update

```go
// data/data.go
var ProviderSet = wire.NewSet(
    New{Service}CommandRepo,
    New{Service}QueryRepo,
)

// biz/biz.go
var ProviderSet = wire.NewSet(
    New{Service}Usecase,
)

// service/service.go
var ProviderSet = wire.NewSet(
    New{Service}Service,
)
```

## 🚀 Commands

```bash
# Generate protobuf
make api

# Generate wire
go generate ./cmd/server/...

# Run migration
DB_PASSWORD=xxx ./scripts/migrate.sh migrations/XXX_create_{name}s_table.sql
```

## ✅ Checklist

- [ ] Entity embeds `BaseEntity`
- [ ] Tách Command/Query repos
- [ ] Command → `GetWriteDB()`
- [ ] Query → `GetReadDB()`
- [ ] Update all ProviderSets
- [ ] Register in http.go & grpc.go
- [ ] Create & run migration
- [ ] Test endpoints

## 🎨 Naming

- Entity: `ProductEntity`
- Repo: `ProductCommandRepo`, `ProductQueryRepo`
- Usecase: `ProductUsecase`
- Service: `ProductService`
- Files: `product.go`, `product_command.go`

## ⚠️ Common Errors

1. Quên tách Command/Query → ❌
2. Dùng sai DB (read cho write) → ❌
3. Quên update Wire → ❌
4. Quên register service → ❌
5. Quên migration → ❌

---

**Xem chi tiết**: `ADD_NEW_SERVICE_GUIDE.md`

