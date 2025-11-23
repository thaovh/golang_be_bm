package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos-layout/internal/biz"
	"github.com/gofrs/uuid/v5"

	"github.com/go-kratos/kratos/v2/log"
)

type userCommandRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserCommandRepo(data *Data, logger log.Logger) biz.UserCommandRepo {
	return &userCommandRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *userCommandRepo) Save(ctx context.Context, u *biz.User) (*biz.User, error) {
	db := r.data.GetWriteDB()
	if err := db.WithContext(ctx).Create(u).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to save user: %v", err)
		return nil, err
	}
	return u, nil
}

func (r *userCommandRepo) Update(ctx context.Context, u *biz.User) (*biz.User, error) {
	db := r.data.GetWriteDB()
	if err := db.WithContext(ctx).Save(u).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to update user: %v", err)
		return nil, err
	}
	return u, nil
}

func (r *userCommandRepo) Delete(ctx context.Context, id uuid.UUID) error {
	db := r.data.GetWriteDB()
	if err := db.WithContext(ctx).Delete(&biz.User{}, "id = ?", id).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to delete user: %v", err)
		return err
	}
	return nil
}

func (r *userCommandRepo) UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error {
	db := r.data.GetWriteDB()
	if err := db.WithContext(ctx).Model(&biz.User{}).
		Where("id = ?", id).
		Update("password_hash", passwordHash).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to update password: %v", err)
		return err
	}
	return nil
}

func (r *userCommandRepo) UpdateLastLogin(ctx context.Context, id uuid.UUID, ip string) error {
	now := time.Now()
	db := r.data.GetWriteDB()
	if err := db.WithContext(ctx).Model(&biz.User{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"last_login_at": &now,
			"last_login_ip": ip,
		}).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to update last login: %v", err)
		return err
	}
	return nil
}

