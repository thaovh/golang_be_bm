package biz

import (
	"context"

	"github.com/gofrs/uuid/v5"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	ErrProvinceNotFound      = errors.NotFound("PROVINCE_NOT_FOUND", "province not found")
	ErrProvinceAlreadyExists = errors.Conflict("PROVINCE_ALREADY_EXISTS", "province code already exists in this country")
	ErrInvalidProvinceCode   = errors.BadRequest("INVALID_PROVINCE_CODE", "invalid province code format")
	ErrCountryRequired       = errors.BadRequest("COUNTRY_REQUIRED", "country is required")
)

// Province là domain model cho Tỉnh/Thành phố
type Province struct {
	BaseEntity

	// Foreign key to Country
	CountryID uuid.UUID `gorm:"type:uuid;not null;index" json:"country_id"`
	Country   *Country  `gorm:"foreignKey:CountryID" json:"country,omitempty"` // Optional: eager load

	// Mã tỉnh/thành phố (unique trong country)
	Code string `gorm:"type:varchar(20);not null;index" json:"code"` // 01, 02, HCM, HN...

	// Tên tỉnh/thành phố
	Name   string `gorm:"type:varchar(255);not null;index" json:"name"`   // Hà Nội
	NameEn string `gorm:"type:varchar(255);not null;index" json:"name_en"` // Hanoi

	// Loại đơn vị hành chính
	Type string `gorm:"type:varchar(50);index" json:"type,omitempty"` // province, city, municipality

	// Thông tin địa lý
	Area        float64 `gorm:"type:decimal(15,2)" json:"area,omitempty"` // km²
	Population  int64   `gorm:"type:bigint" json:"population,omitempty"`
	Coordinates string  `gorm:"type:varchar(100)" json:"coordinates,omitempty"` // lat,lng

	// Thông tin hành chính
	Capital     string `gorm:"type:varchar(100)" json:"capital,omitempty"` // Thủ phủ/Trung tâm
	PostalCode  string `gorm:"type:varchar(20)" json:"postal_code,omitempty"`
	PhonePrefix string `gorm:"type:varchar(10)" json:"phone_prefix,omitempty"` // 024, 028...

	// Thứ tự sắp xếp
	SortOrder int `gorm:"type:integer;default:0;index" json:"sort_order,omitempty"`
}

// ProvinceCommandRepo là repository interface cho write operations
type ProvinceCommandRepo interface {
	Save(context.Context, *Province) (*Province, error)
	Update(context.Context, *Province) (*Province, error)
	Delete(context.Context, uuid.UUID) error
}

// ProvinceQueryRepo là repository interface cho read operations
type ProvinceQueryRepo interface {
	FindByID(context.Context, uuid.UUID) (*Province, error)
	FindByCode(context.Context, string, uuid.UUID) (*Province, error) // code + country_id
	List(context.Context, *ProvinceListFilter) ([]*Province, int64, error)
	ListByCountry(context.Context, uuid.UUID) ([]*Province, error)
	Search(context.Context, string) ([]*Province, error)
	Count(context.Context, *ProvinceListFilter) (int64, error)
}

// ProvinceListFilter cho pagination và filtering
type ProvinceListFilter struct {
	Page      int32
	PageSize  int32
	Search    string // Search by name, name_en, code
	CountryID uuid.UUID // Filter by country
	Type      string    // Filter by type (province, city, municipality)
	Status    string    // Filter by status
	Code      string    // Filter by exact code
}

// ProvinceUsecase là usecase cho Province với CQRS pattern
type ProvinceUsecase struct {
	commandRepo ProvinceCommandRepo
	queryRepo   ProvinceQueryRepo
	countryRepo CountryQueryRepo // To validate country exists
	log         *log.Helper
}

// NewProvinceUsecase tạo ProvinceUsecase mới
func NewProvinceUsecase(
	commandRepo ProvinceCommandRepo,
	queryRepo ProvinceQueryRepo,
	countryRepo CountryQueryRepo,
	logger log.Logger,
) *ProvinceUsecase {
	return &ProvinceUsecase{
		commandRepo: commandRepo,
		queryRepo:   queryRepo,
		countryRepo: countryRepo,
		log:         log.NewHelper(logger),
	}
}

// CreateProvince creates a new province (Command)
func (uc *ProvinceUsecase) CreateProvince(ctx context.Context, province *Province) (*Province, error) {
	uc.log.WithContext(ctx).Infof("CreateProvince: %s (%s)", province.Name, province.Code)

	// Validate country ID is provided
	if province.CountryID == uuid.Nil {
		return nil, ErrCountryRequired
	}

	// Validate country exists
	country, err := uc.countryRepo.FindByID(ctx, province.CountryID)
	if err != nil {
		return nil, err
	}
	if country == nil {
		return nil, ErrCountryNotFound
	}

	// Validate province code is not empty
	if province.Code == "" {
		return nil, ErrInvalidProvinceCode
	}

	// Check if province code already exists in this country
	existing, _ := uc.queryRepo.FindByCode(ctx, province.Code, province.CountryID)
	if existing != nil {
		return nil, ErrProvinceAlreadyExists
	}

	// Set audit fields from context
	province.SetAuditFields(ctx, true)

	return uc.commandRepo.Save(ctx, province)
}

// UpdateProvince updates an existing province (Command)
func (uc *ProvinceUsecase) UpdateProvince(ctx context.Context, province *Province) (*Province, error) {
	uc.log.WithContext(ctx).Infof("UpdateProvince: %s", province.ID.String())

	// Validate province code if provided
	if province.Code != "" && len(province.Code) == 0 {
		return nil, ErrInvalidProvinceCode
	}

	// Check if province exists
	existing, err := uc.queryRepo.FindByID(ctx, province.ID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrProvinceNotFound
	}

	// If country changed, validate new country exists
	if province.CountryID != uuid.Nil && province.CountryID != existing.CountryID {
		country, err := uc.countryRepo.FindByID(ctx, province.CountryID)
		if err != nil {
			return nil, err
		}
		if country == nil {
			return nil, ErrCountryNotFound
		}
	}

	// If code changed, check for duplicate in the country
	if province.Code != "" && province.Code != existing.Code {
		// Use existing country ID if not changed
		countryID := province.CountryID
		if countryID == uuid.Nil {
			countryID = existing.CountryID
		}
		duplicate, _ := uc.queryRepo.FindByCode(ctx, province.Code, countryID)
		if duplicate != nil && duplicate.ID != province.ID {
			return nil, ErrProvinceAlreadyExists
		}
	}

	// Set audit fields from context
	province.SetAuditFields(ctx, false)

	return uc.commandRepo.Update(ctx, province)
}

// DeleteProvince deletes a province (Command)
func (uc *ProvinceUsecase) DeleteProvince(ctx context.Context, id uuid.UUID) error {
	uc.log.WithContext(ctx).Infof("DeleteProvince: %s", id.String())
	return uc.commandRepo.Delete(ctx, id)
}

// GetProvince retrieves a province by ID (Query)
func (uc *ProvinceUsecase) GetProvince(ctx context.Context, id uuid.UUID) (*Province, error) {
	uc.log.WithContext(ctx).Infof("GetProvince: %s", id.String())
	province, err := uc.queryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if province == nil {
		return nil, ErrProvinceNotFound
	}
	return province, nil
}

// GetProvinceByCode retrieves a province by its code and country ID (Query)
func (uc *ProvinceUsecase) GetProvinceByCode(ctx context.Context, code string, countryID uuid.UUID) (*Province, error) {
	uc.log.WithContext(ctx).Infof("GetProvinceByCode: %s (country: %s)", code, countryID.String())
	province, err := uc.queryRepo.FindByCode(ctx, code, countryID)
	if err != nil {
		return nil, err
	}
	if province == nil {
		return nil, ErrProvinceNotFound
	}
	return province, nil
}

// ListProvinces lists provinces with pagination and filters (Query)
func (uc *ProvinceUsecase) ListProvinces(ctx context.Context, filter *ProvinceListFilter) ([]*Province, int64, error) {
	uc.log.WithContext(ctx).Infof("ListProvinces: Page %d, PageSize %d, CountryID %s", filter.Page, filter.PageSize, filter.CountryID.String())
	return uc.queryRepo.List(ctx, filter)
}

// ListProvincesByCountry lists all provinces of a specific country (Query)
func (uc *ProvinceUsecase) ListProvincesByCountry(ctx context.Context, countryID uuid.UUID) ([]*Province, error) {
	uc.log.WithContext(ctx).Infof("ListProvincesByCountry: %s", countryID.String())
	return uc.queryRepo.ListByCountry(ctx, countryID)
}

// SearchProvinces searches provinces by name or code (Query)
func (uc *ProvinceUsecase) SearchProvinces(ctx context.Context, query string) ([]*Province, error) {
	uc.log.WithContext(ctx).Infof("SearchProvinces: %s", query)
	return uc.queryRepo.Search(ctx, query)
}

