package data

import (
	"context"
	"strings"

	"github.com/go-kratos/kratos-layout/internal/biz"
	"github.com/gofrs/uuid/v5"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type countryQueryRepo struct {
	data *Data
	log  *log.Helper
}

func NewCountryQueryRepo(data *Data, logger log.Logger) biz.CountryQueryRepo {
	return &countryQueryRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *countryQueryRepo) FindByID(ctx context.Context, id uuid.UUID) (*biz.Country, error) {
	db := r.data.GetReadDB()
	var country biz.Country
	if err := db.WithContext(ctx).Where("id = ?", id).First(&country).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.log.WithContext(ctx).Errorf("Failed to find country by ID: %v", err)
		return nil, err
	}
	return &country, nil
}

func (r *countryQueryRepo) FindByCode(ctx context.Context, code string) (*biz.Country, error) {
	db := r.data.GetReadDB()
	var country biz.Country
	if err := db.WithContext(ctx).Where("code = ?", strings.ToUpper(code)).First(&country).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.log.WithContext(ctx).Errorf("Failed to find country by code: %v", err)
		return nil, err
	}
	return &country, nil
}

func (r *countryQueryRepo) List(ctx context.Context, filter *biz.CountryListFilter) ([]*biz.Country, int64, error) {
	db := r.data.GetReadDB()
	var countries []*biz.Country
	var total int64

	query := db.WithContext(ctx).Model(&biz.Country{})

	// Apply filters
	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		query = query.Where("name LIKE ? OR name_en LIKE ?", searchPattern, searchPattern)
	}
	if filter.Region != "" {
		query = query.Where("region = ?", filter.Region)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Code != "" {
		query = query.Where("code = ?", strings.ToUpper(filter.Code))
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to count countries: %v", err)
		return nil, 0, err
	}

	// Pagination
	if filter.PageSize > 0 {
		offset := (filter.Page - 1) * filter.PageSize
		query = query.Offset(int(offset)).Limit(int(filter.PageSize))
	}

	// Order by name
	query = query.Order("name ASC")

	if err := query.Find(&countries).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to list countries: %v", err)
		return nil, 0, err
	}

	return countries, total, nil
}

func (r *countryQueryRepo) Count(ctx context.Context, filter *biz.CountryListFilter) (int64, error) {
	db := r.data.GetReadDB()
	var count int64

	query := db.WithContext(ctx).Model(&biz.Country{})

	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		query = query.Where("name LIKE ? OR name_en LIKE ?", searchPattern, searchPattern)
	}
	if filter.Region != "" {
		query = query.Where("region = ?", filter.Region)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	if filter.Code != "" {
		query = query.Where("code = ?", strings.ToUpper(filter.Code))
	}

	if err := query.Count(&count).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to count countries: %v", err)
		return 0, err
	}

	return count, nil
}

func (r *countryQueryRepo) Search(ctx context.Context, query string) ([]*biz.Country, error) {
	db := r.data.GetReadDB()
	var countries []*biz.Country

	searchPattern := "%" + query + "%"
	if err := db.WithContext(ctx).
		Where("name LIKE ? OR name_en LIKE ? OR code LIKE ?", searchPattern, searchPattern, strings.ToUpper(query)).
		Order("name ASC").
		Limit(50). // Limit search results
		Find(&countries).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to search countries: %v", err)
		return nil, err
	}

	return countries, nil
}

