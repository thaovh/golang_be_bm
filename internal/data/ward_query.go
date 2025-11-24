package data

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos-layout/internal/biz"
	"github.com/gofrs/uuid/v5"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type wardQueryRepo struct {
	data *Data
	log  *log.Helper
}

// NewWardQueryRepo creates a new WardQueryRepo
func NewWardQueryRepo(data *Data, logger log.Logger) biz.WardQueryRepo {
	return &wardQueryRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// FindByID finds a ward by ID from the read database
func (r *wardQueryRepo) FindByID(ctx context.Context, id uuid.UUID) (*biz.Ward, error) {
	db := r.data.GetReadDB()
	var ward biz.Ward
	if err := db.WithContext(ctx).Where("id = ?", id).First(&ward).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Return nil, nil if not found
		}
		r.log.WithContext(ctx).Errorf("Failed to find ward by ID: %v", err)
		return nil, err
	}
	return &ward, nil
}

// FindByCode finds a ward by its code and province ID from the read database
func (r *wardQueryRepo) FindByCode(ctx context.Context, code string, provinceID uuid.UUID) (*biz.Ward, error) {
	db := r.data.GetReadDB()
	var ward biz.Ward
	if err := db.WithContext(ctx).Where("code = ? AND province_id = ?", code, provinceID).First(&ward).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Return nil, nil if not found
		}
		r.log.WithContext(ctx).Errorf("Failed to find ward by code: %v", err)
		return nil, err
	}
	return &ward, nil
}

// List lists wards with pagination and filters from the read database
func (r *wardQueryRepo) List(ctx context.Context, filter *biz.WardListFilter) ([]*biz.Ward, int64, error) {
	db := r.data.GetReadDB()
	var wards []*biz.Ward
	var total int64

	query := db.WithContext(ctx).Model(&biz.Ward{})

	if filter.Search != "" {
		search := fmt.Sprintf("%%%s%%", filter.Search)
		query = query.Where("name ILIKE ? OR name_en ILIKE ? OR code ILIKE ?", search, search, search)
	}
	if filter.ProvinceID != uuid.Nil {
		query = query.Where("province_id = ?", filter.ProvinceID)
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
		r.log.WithContext(ctx).Errorf("Failed to count wards: %v", err)
		return nil, 0, err
	}

	if filter.Page > 0 && filter.PageSize > 0 {
		offset := (filter.Page - 1) * filter.PageSize
		query = query.Offset(int(offset)).Limit(int(filter.PageSize))
	}

	// Order by sort_order, then by name
	query = query.Order("sort_order ASC, name ASC")

	if err := query.Find(&wards).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to list wards: %v", err)
		return nil, 0, err
	}

	return wards, total, nil
}

// ListByProvince lists all wards of a specific province from the read database
func (r *wardQueryRepo) ListByProvince(ctx context.Context, provinceID uuid.UUID) ([]*biz.Ward, error) {
	db := r.data.GetReadDB()
	var wards []*biz.Ward
	if err := db.WithContext(ctx).Where("province_id = ?", provinceID).Order("sort_order ASC, name ASC").Find(&wards).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to list wards by province: %v", err)
		return nil, err
	}
	return wards, nil
}

// Search searches wards by name or code from the read database
func (r *wardQueryRepo) Search(ctx context.Context, search string) ([]*biz.Ward, error) {
	db := r.data.GetReadDB()
	var wards []*biz.Ward
	searchPattern := fmt.Sprintf("%%%s%%", search)
	if err := db.WithContext(ctx).Where("name ILIKE ? OR name_en ILIKE ? OR code ILIKE ?", searchPattern, searchPattern, searchPattern).Order("sort_order ASC, name ASC").Find(&wards).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to search wards: %v", err)
		return nil, err
	}
	return wards, nil
}

// Count counts wards based on filters from the read database
func (r *wardQueryRepo) Count(ctx context.Context, filter *biz.WardListFilter) (int64, error) {
	db := r.data.GetReadDB()
	var count int64

	query := db.WithContext(ctx).Model(&biz.Ward{})

	if filter.Search != "" {
		search := fmt.Sprintf("%%%s%%", filter.Search)
		query = query.Where("name ILIKE ? OR name_en ILIKE ? OR code ILIKE ?", search, search, search)
	}
	if filter.ProvinceID != uuid.Nil {
		query = query.Where("province_id = ?", filter.ProvinceID)
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
		r.log.WithContext(ctx).Errorf("Failed to count wards: %v", err)
		return 0, err
	}
	return count, nil
}

