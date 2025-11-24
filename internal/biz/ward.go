package biz

import (
	"context"

	"github.com/gofrs/uuid/v5"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	ErrWardNotFound      = errors.NotFound("WARD_NOT_FOUND", "ward not found")
	ErrWardAlreadyExists = errors.Conflict("WARD_ALREADY_EXISTS", "ward code already exists in this province")
	ErrInvalidWardCode   = errors.BadRequest("INVALID_WARD_CODE", "invalid ward code format")
	ErrProvinceRequired  = errors.BadRequest("PROVINCE_REQUIRED", "province is required")
)

// Ward là domain model cho Xã/Phường
type Ward struct {
	BaseEntity

	// Foreign key to Province
	ProvinceID uuid.UUID `gorm:"type:uuid;not null;index" json:"province_id"`
	Province   *Province `gorm:"foreignKey:ProvinceID" json:"province,omitempty"` // Optional: eager load

	// Mã xã/phường (unique trong province)
	Code string `gorm:"type:varchar(20);not null;index" json:"code"` // 00001, 00002...

	// Tên xã/phường
	Name   string `gorm:"type:varchar(255);not null;index" json:"name"`   // Phường Cửa Đông
	NameEn string `gorm:"type:varchar(255);not null;index" json:"name_en"` // Cua Dong Ward

	// Loại đơn vị hành chính
	Type string `gorm:"type:varchar(50);index" json:"type,omitempty"` // ward, commune, town

	// Thông tin địa lý
	Area        float64 `gorm:"type:decimal(15,2)" json:"area,omitempty"` // km²
	Population  int64   `gorm:"type:bigint" json:"population,omitempty"`
	Coordinates string  `gorm:"type:varchar(100)" json:"coordinates,omitempty"` // lat,lng

	// Thông tin hành chính
	PostalCode string `gorm:"type:varchar(20)" json:"postal_code,omitempty"`
	Address    string `gorm:"type:varchar(500)" json:"address,omitempty"` // Địa chỉ trụ sở

	// Thứ tự sắp xếp
	SortOrder int `gorm:"type:integer;default:0;index" json:"sort_order,omitempty"`
}

// WardCommandRepo là repository interface cho write operations
type WardCommandRepo interface {
	Save(context.Context, *Ward) (*Ward, error)
	Update(context.Context, *Ward) (*Ward, error)
	Delete(context.Context, uuid.UUID) error
}

// WardQueryRepo là repository interface cho read operations
type WardQueryRepo interface {
	FindByID(context.Context, uuid.UUID) (*Ward, error)
	FindByCode(context.Context, string, uuid.UUID) (*Ward, error) // code + province_id
	List(context.Context, *WardListFilter) ([]*Ward, int64, error)
	ListByProvince(context.Context, uuid.UUID) ([]*Ward, error)
	Search(context.Context, string) ([]*Ward, error)
	Count(context.Context, *WardListFilter) (int64, error)
}

// WardListFilter cho pagination và filtering
type WardListFilter struct {
	Page       int32
	PageSize   int32
	Search     string // Search by name, name_en, code
	ProvinceID uuid.UUID // Filter by province
	Type       string    // Filter by type (ward, commune, town)
	Status     string    // Filter by status
	Code       string    // Filter by exact code
}

// WardUsecase là usecase cho Ward với CQRS pattern
type WardUsecase struct {
	commandRepo  WardCommandRepo
	queryRepo    WardQueryRepo
	provinceRepo ProvinceQueryRepo // To validate province exists
	log          *log.Helper
}

// NewWardUsecase tạo WardUsecase mới
func NewWardUsecase(
	commandRepo WardCommandRepo,
	queryRepo WardQueryRepo,
	provinceRepo ProvinceQueryRepo,
	logger log.Logger,
) *WardUsecase {
	return &WardUsecase{
		commandRepo:  commandRepo,
		queryRepo:    queryRepo,
		provinceRepo: provinceRepo,
		log:          log.NewHelper(logger),
	}
}

// CreateWard creates a new ward (Command)
func (uc *WardUsecase) CreateWard(ctx context.Context, ward *Ward) (*Ward, error) {
	uc.log.WithContext(ctx).Infof("CreateWard: %s (%s)", ward.Name, ward.Code)

	// Validate province ID is provided
	if ward.ProvinceID == uuid.Nil {
		return nil, ErrProvinceRequired
	}

	// Validate province exists
	province, err := uc.provinceRepo.FindByID(ctx, ward.ProvinceID)
	if err != nil {
		return nil, err
	}
	if province == nil {
		return nil, ErrProvinceNotFound
	}

	// Validate ward code is not empty
	if ward.Code == "" {
		return nil, ErrInvalidWardCode
	}

	// Check if ward code already exists in this province
	existing, _ := uc.queryRepo.FindByCode(ctx, ward.Code, ward.ProvinceID)
	if existing != nil {
		return nil, ErrWardAlreadyExists
	}

	// Set audit fields from context
	ward.SetAuditFields(ctx, true)

	return uc.commandRepo.Save(ctx, ward)
}

// UpdateWard updates an existing ward (Command)
func (uc *WardUsecase) UpdateWard(ctx context.Context, ward *Ward) (*Ward, error) {
	uc.log.WithContext(ctx).Infof("UpdateWard: %s", ward.ID.String())

	// Validate ward code if provided
	if ward.Code != "" && len(ward.Code) == 0 {
		return nil, ErrInvalidWardCode
	}

	// Check if ward exists
	existing, err := uc.queryRepo.FindByID(ctx, ward.ID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrWardNotFound
	}

	// If province changed, validate new province exists
	if ward.ProvinceID != uuid.Nil && ward.ProvinceID != existing.ProvinceID {
		province, err := uc.provinceRepo.FindByID(ctx, ward.ProvinceID)
		if err != nil {
			return nil, err
		}
		if province == nil {
			return nil, ErrProvinceNotFound
		}
	}

	// If code changed, check for duplicate in the province
	if ward.Code != "" && ward.Code != existing.Code {
		// Use existing province ID if not changed
		provinceID := ward.ProvinceID
		if provinceID == uuid.Nil {
			provinceID = existing.ProvinceID
		}
		duplicate, _ := uc.queryRepo.FindByCode(ctx, ward.Code, provinceID)
		if duplicate != nil && duplicate.ID != ward.ID {
			return nil, ErrWardAlreadyExists
		}
	}

	// Set audit fields from context
	ward.SetAuditFields(ctx, false)

	return uc.commandRepo.Update(ctx, ward)
}

// DeleteWard deletes a ward (Command)
func (uc *WardUsecase) DeleteWard(ctx context.Context, id uuid.UUID) error {
	uc.log.WithContext(ctx).Infof("DeleteWard: %s", id.String())
	return uc.commandRepo.Delete(ctx, id)
}

// GetWard retrieves a ward by ID (Query)
func (uc *WardUsecase) GetWard(ctx context.Context, id uuid.UUID) (*Ward, error) {
	uc.log.WithContext(ctx).Infof("GetWard: %s", id.String())
	ward, err := uc.queryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if ward == nil {
		return nil, ErrWardNotFound
	}
	return ward, nil
}

// GetWardByCode retrieves a ward by its code and province ID (Query)
func (uc *WardUsecase) GetWardByCode(ctx context.Context, code string, provinceID uuid.UUID) (*Ward, error) {
	uc.log.WithContext(ctx).Infof("GetWardByCode: %s (province: %s)", code, provinceID.String())
	ward, err := uc.queryRepo.FindByCode(ctx, code, provinceID)
	if err != nil {
		return nil, err
	}
	if ward == nil {
		return nil, ErrWardNotFound
	}
	return ward, nil
}

// ListWards lists wards with pagination and filters (Query)
func (uc *WardUsecase) ListWards(ctx context.Context, filter *WardListFilter) ([]*Ward, int64, error) {
	uc.log.WithContext(ctx).Infof("ListWards: Page %d, PageSize %d, ProvinceID %s", filter.Page, filter.PageSize, filter.ProvinceID.String())
	return uc.queryRepo.List(ctx, filter)
}

// ListWardsByProvince lists all wards of a specific province (Query)
func (uc *WardUsecase) ListWardsByProvince(ctx context.Context, provinceID uuid.UUID) ([]*Ward, error) {
	uc.log.WithContext(ctx).Infof("ListWardsByProvince: %s", provinceID.String())
	return uc.queryRepo.ListByProvince(ctx, provinceID)
}

// SearchWards searches wards by name or code (Query)
func (uc *WardUsecase) SearchWards(ctx context.Context, query string) ([]*Ward, error) {
	uc.log.WithContext(ctx).Infof("SearchWards: %s", query)
	return uc.queryRepo.Search(ctx, query)
}

