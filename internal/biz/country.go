package biz

import (
	"context"

	"github.com/gofrs/uuid/v5"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	ErrCountryNotFound      = errors.NotFound("COUNTRY_NOT_FOUND", "country not found")
	ErrCountryAlreadyExists = errors.Conflict("COUNTRY_ALREADY_EXISTS", "country already exists")
	ErrInvalidCountryCode   = errors.BadRequest("INVALID_COUNTRY_CODE", "invalid country code format")
)

// Country là domain model cho Quốc gia
type Country struct {
	BaseEntity

	// Mã quốc gia (ISO 3166-1 alpha-2)
	Code string `gorm:"type:varchar(2);uniqueIndex;not null" json:"code"` // VN, US, JP...

	// Tên quốc gia
	Name   string `gorm:"type:varchar(255);not null;index" json:"name"`   // Việt Nam
	NameEn string `gorm:"type:varchar(255);not null;index" json:"name_en"` // Vietnam

	// Thông tin địa lý
	Region    string `gorm:"type:varchar(100);index" json:"region,omitempty"`     // Asia, Europe...
	SubRegion string `gorm:"type:varchar(100)" json:"sub_region,omitempty"`       // Southeast Asia

	// Thông tin kinh tế
	CurrencyCode   string `gorm:"type:varchar(3)" json:"currency_code,omitempty"`   // VND, USD
	CurrencySymbol string `gorm:"type:varchar(10)" json:"currency_symbol,omitempty"` // ₫, $

	// Thông tin liên lạc
	PhoneCode string `gorm:"type:varchar(10)" json:"phone_code,omitempty"`     // +84
	TimeZone  string `gorm:"type:varchar(50)" json:"time_zone,omitempty"`      // Asia/Ho_Chi_Minh

	// Thông tin bổ sung
	Flag       string `gorm:"type:varchar(10)" json:"flag,omitempty"`         // 🇻🇳 (emoji hoặc URL)
	Capital    string `gorm:"type:varchar(100)" json:"capital,omitempty"`     // Hà Nội
	Population int64  `gorm:"type:bigint" json:"population,omitempty"`

	// Metadata
	ISO3166Alpha3  string `gorm:"type:varchar(3);uniqueIndex" json:"iso3166_alpha3,omitempty"`  // VNM
	ISO3166Numeric string `gorm:"type:varchar(3)" json:"iso3166_numeric,omitempty"`            // 704
}

// CountryCommandRepo là repository interface cho write operations
type CountryCommandRepo interface {
	Save(context.Context, *Country) (*Country, error)
	Update(context.Context, *Country) (*Country, error)
	Delete(context.Context, uuid.UUID) error
}

// CountryQueryRepo là repository interface cho read operations
type CountryQueryRepo interface {
	FindByID(context.Context, uuid.UUID) (*Country, error)
	FindByCode(context.Context, string) (*Country, error)
	List(context.Context, *CountryListFilter) ([]*Country, int64, error)
	Count(context.Context, *CountryListFilter) (int64, error)
	Search(context.Context, string) ([]*Country, error)
}

// CountryListFilter cho pagination và filtering
type CountryListFilter struct {
	Page     int32
	PageSize int32
	Search   string // Search by name or name_en
	Region   string // Filter by region
	Status   string // Filter by status
	Code     string // Filter by code
}

// CountryUsecase là usecase cho Country với CQRS pattern
type CountryUsecase struct {
	commandRepo CountryCommandRepo
	queryRepo   CountryQueryRepo
	log         *log.Helper
}

// NewCountryUsecase tạo CountryUsecase mới
func NewCountryUsecase(
	commandRepo CountryCommandRepo,
	queryRepo CountryQueryRepo,
	logger log.Logger,
) *CountryUsecase {
	return &CountryUsecase{
		commandRepo: commandRepo,
		queryRepo:   queryRepo,
		log:         log.NewHelper(logger),
	}
}

// CreateCountry creates a new country (Command)
func (uc *CountryUsecase) CreateCountry(ctx context.Context, country *Country) (*Country, error) {
	uc.log.WithContext(ctx).Infof("CreateCountry: %s (%s)", country.Name, country.Code)

	// Validate country code format (should be 2 uppercase letters)
	if len(country.Code) != 2 {
		return nil, ErrInvalidCountryCode
	}

	// Check if country code already exists
	existing, _ := uc.queryRepo.FindByCode(ctx, country.Code)
	if existing != nil {
		return nil, ErrCountryAlreadyExists
	}

	// Set audit fields from context
	country.SetAuditFields(ctx, true)

	return uc.commandRepo.Save(ctx, country)
}

// UpdateCountry updates a country (Command)
func (uc *CountryUsecase) UpdateCountry(ctx context.Context, country *Country) (*Country, error) {
	uc.log.WithContext(ctx).Infof("UpdateCountry: %s", country.ID.String())

	// Validate country code format if provided
	if country.Code != "" && len(country.Code) != 2 {
		return nil, ErrInvalidCountryCode
	}

	// Set audit fields from context
	country.SetAuditFields(ctx, false)

	return uc.commandRepo.Update(ctx, country)
}

// DeleteCountry deletes a country (Command)
func (uc *CountryUsecase) DeleteCountry(ctx context.Context, id uuid.UUID) error {
	uc.log.WithContext(ctx).Infof("DeleteCountry: %s", id.String())
	return uc.commandRepo.Delete(ctx, id)
}

// GetCountry gets a country by ID (Query)
func (uc *CountryUsecase) GetCountry(ctx context.Context, id uuid.UUID) (*Country, error) {
	country, err := uc.queryRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if country == nil {
		return nil, ErrCountryNotFound
	}
	return country, nil
}

// GetCountryByCode gets a country by code (Query)
func (uc *CountryUsecase) GetCountryByCode(ctx context.Context, code string) (*Country, error) {
	country, err := uc.queryRepo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	if country == nil {
		return nil, ErrCountryNotFound
	}
	return country, nil
}

// ListCountries lists countries with pagination and filters (Query)
func (uc *CountryUsecase) ListCountries(ctx context.Context, filter *CountryListFilter) ([]*Country, int64, error) {
	return uc.queryRepo.List(ctx, filter)
}

// SearchCountries searches countries by name (Query)
func (uc *CountryUsecase) SearchCountries(ctx context.Context, query string) ([]*Country, error) {
	return uc.queryRepo.Search(ctx, query)
}

