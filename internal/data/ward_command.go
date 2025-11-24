package data

import (
	"context"

	"github.com/go-kratos/kratos-layout/internal/biz"
	"github.com/gofrs/uuid/v5"

	"github.com/go-kratos/kratos/v2/log"
)

type wardCommandRepo struct {
	data *Data
	log  *log.Helper
}

// NewWardCommandRepo creates a new WardCommandRepo
func NewWardCommandRepo(data *Data, logger log.Logger) biz.WardCommandRepo {
	return &wardCommandRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// Save saves a Ward using write database
func (r *wardCommandRepo) Save(ctx context.Context, w *biz.Ward) (*biz.Ward, error) {
	db := r.data.GetWriteDB()
	if err := db.WithContext(ctx).Create(w).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to save ward: %v", err)
		return nil, err
	}
	return w, nil
}

// Update updates a Ward using write database
func (r *wardCommandRepo) Update(ctx context.Context, w *biz.Ward) (*biz.Ward, error) {
	db := r.data.GetWriteDB()
	if err := db.WithContext(ctx).Save(w).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to update ward: %v", err)
		return nil, err
	}
	return w, nil
}

// Delete deletes a Ward using write database
func (r *wardCommandRepo) Delete(ctx context.Context, id uuid.UUID) error {
	db := r.data.GetWriteDB()
	if err := db.WithContext(ctx).Delete(&biz.Ward{}, "id = ?", id).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to delete ward: %v", err)
		return nil, err
	}
	return nil
}

