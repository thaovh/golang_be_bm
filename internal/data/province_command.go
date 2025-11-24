package data

import (
	"context"

	"github.com/go-kratos/kratos-layout/internal/biz"
	"github.com/gofrs/uuid/v5"

	"github.com/go-kratos/kratos/v2/log"
)

type provinceCommandRepo struct {
	data *Data
	log  *log.Helper
}

// NewProvinceCommandRepo creates a new ProvinceCommandRepo
func NewProvinceCommandRepo(data *Data, logger log.Logger) biz.ProvinceCommandRepo {
	return &provinceCommandRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// Save saves a Province using write database
func (r *provinceCommandRepo) Save(ctx context.Context, p *biz.Province) (*biz.Province, error) {
	db := r.data.GetWriteDB()
	if err := db.WithContext(ctx).Create(p).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to save province: %v", err)
		return nil, err
	}
	return p, nil
}

// Update updates a Province using write database
func (r *provinceCommandRepo) Update(ctx context.Context, p *biz.Province) (*biz.Province, error) {
	db := r.data.GetWriteDB()
	if err := db.WithContext(ctx).Save(p).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to update province: %v", err)
		return nil, err
	}
	return p, nil
}

// Delete deletes a Province using write database
func (r *provinceCommandRepo) Delete(ctx context.Context, id uuid.UUID) error {
	db := r.data.GetWriteDB()
	if err := db.WithContext(ctx).Delete(&biz.Province{}, "id = ?", id).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to delete province: %v", err)
		return err
	}
	return nil
}

