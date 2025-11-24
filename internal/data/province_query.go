package data

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos-layout/internal/biz"
	"github.com/gofrs/uuid/v5"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type provinceQueryRepo struct {
	data *Data
	log  *log.Helper
}

// NewProvinceQueryRepo creates a new ProvinceQueryRepo
func NewProvinceQueryRepo(data *Data, logger log.Logger) biz.ProvinceQueryRepo {
	return &provinceQueryRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// FindByID finds a province by ID from the read database
func (r *provinceQueryRepo) FindByID(ctx context.Context, id uuid.UUID) (*biz.Province, error) {
	db := r.data.GetReadDB()
	var province biz.Province
	if err := db.WithContext(ctx).Where("id = ?", id).First(&province).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Return nil, nil if not found
		}
		r.log.WithContext(ctx).Errorf("Failed to find province by ID: %v", err)
		return nil, err
	}
	return &province, nil
}

// FindByCode finds a province by its code and country ID from the read database
func (r *provinceQueryRepo) FindByCode(ctx context.Context, code string, countryID uuid.UUID) (*biz.Province, error) {
	db := r.data.GetReadDB()
	var province biz.Province
	if err := db.WithContext(ctx).Where("code = ? AND country_id = ?", code, countryID).First(&province).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Return nil, nil if not found
		}
		r.log.WithContext(ctx).Errorf("Failed to find province by code: %v", err)
		return nil, err
	}
	return &province, nil
}

// List lists provinces with pagination and filters from the read database
func (r *provinceQueryRepo) List(ctx context.Context, filter *biz.ProvinceListFilter) ([]*biz.Province, int64, error) {
	db := r.data.GetReadDB()
	var provinces []*biz.Province
	var total int64

	query := db.WithContext(ctx).Model(&biz.Province{})

	if filter.Search != "" {
		search := fmt.Sprintf("%%%s%%", filter.Search)
		query = query.Where("name ILIKE ? OR name_en ILIKE ? OR code ILIKE ?", search, search, search)
	}
	if filter.CountryID != uuid.Nil {
		query = query.Where("country_id = ?", filter.CountryID)
	}
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Code != "" {
		query = query.Where("code = ?", filter.Code)
	}

	if err := query.Count(&total).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to count provinces: %v", err)
		return nil, 0, err
	}

	if filter.Page > 0 && filter.PageSize > 0 {
		offset := (filter.Page - 1) * filter.PageSize
		query = query.Offset(int(offset)).Limit(int(filter.PageSize))
	}

	// Order by sort_order, then by name
	query = query.Order("sort_order ASC, name ASC")

	if err := query.Find(&provinces).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to list provinces: %v", err)
		return nil, 0, err
	}

	return provinces, total, nil
}

// ListByCountry lists all provinces of a specific country from the read database
func (r *provinceQueryRepo) ListByCountry(ctx context.Context, countryID uuid.UUID) ([]*biz.Province, error) {
	db := r.data.GetReadDB()
	var provinces []*biz.Province
	if err := db.WithContext(ctx).Where("country_id = ?", countryID).Order("sort_order ASC, name ASC").Find(&provinces).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to list provinces by country: %v", err)
		return nil, err
	}
	return provinces, nil
}

// Search searches provinces by name or code from the read database
func (r *provinceQueryRepo) Search(ctx context.Context, search string) ([]*biz.Province, error) {
	db := r.data.GetReadDB()
	var provinces []*biz.Province
	searchPattern := fmt.Sprintf("%%%s%%", search)
	if err := db.WithContext(ctx).Where("name ILIKE ? OR name_en ILIKE ? OR code ILIKE ?", searchPattern, searchPattern, searchPattern).Order("sort_order ASC, name ASC").Find(&provinces).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to search provinces: %v", err)
		return nil, err
	}
	return provinces, nil
}

// Count counts provinces based on filters from the read database
func (r *provinceQueryRepo) Count(ctx context.Context, filter *biz.ProvinceListFilter) (int64, error) {
	db := r.data.GetReadDB()
	var count int64

	query := db.WithContext(ctx).Model(&biz.Province{})

	if filter.Search != "" {
		search := fmt.Sprintf("%%%s%%", filter.Search)
		query = query.Where("name ILIKE ? OR name_en ILIKE ? OR code ILIKE ?", search, search, search)
	}
	if filter.CountryID != uuid.Nil {
		query = query.Where("country_id = ?", filter.CountryID)
	}
	if filter.Type != "" {
		query = query.Where("type = ?", filter.Type)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Code != "" {
		query = query.Where("code = ?", filter.Code)
	}

	if err := query.Count(&count).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to count provinces: %v", err)
		return 0, err
	}
	return count, nil
}

