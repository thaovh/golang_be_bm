package data

import (
	"context"

	"github.com/go-kratos/kratos-layout/internal/biz"
	"github.com/gofrs/uuid/v5"

	"github.com/go-kratos/kratos/v2/log"
)

type countryCommandRepo struct {
	data *Data
	log  *log.Helper
}

func NewCountryCommandRepo(data *Data, logger log.Logger) biz.CountryCommandRepo {
	return &countryCommandRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *countryCommandRepo) Save(ctx context.Context, c *biz.Country) (*biz.Country, error) {
	db := r.data.GetWriteDB()
	if err := db.WithContext(ctx).Create(c).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to save country: %v", err)
		return nil, err
	}
	return c, nil
}

func (r *countryCommandRepo) Update(ctx context.Context, c *biz.Country) (*biz.Country, error) {
	db := r.data.GetWriteDB()
	if err := db.WithContext(ctx).Model(c).Updates(c).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to update country: %v", err)
		return nil, err
	}
	return c, nil
}

func (r *countryCommandRepo) Delete(ctx context.Context, id uuid.UUID) error {
	db := r.data.GetWriteDB()
	if err := db.WithContext(ctx).Delete(&biz.Country{}, "id = ?", id).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to delete country: %v", err)
		return err
	}
	return nil
}

