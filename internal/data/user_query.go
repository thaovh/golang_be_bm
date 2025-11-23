package data

import (
	"context"
	"fmt"

	"github.com/go-kratos/kratos-layout/internal/biz"
	"github.com/gofrs/uuid/v5"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type userQueryRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserQueryRepo(data *Data, logger log.Logger) biz.UserQueryRepo {
	return &userQueryRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *userQueryRepo) FindByID(ctx context.Context, id uuid.UUID) (*biz.User, error) {
	db := r.data.GetReadDB()
	var user biz.User
	if err := db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, biz.ErrUserNotFound
		}
		r.log.WithContext(ctx).Errorf("Failed to find user by ID: %v", err)
		return nil, err
	}
	return &user, nil
}

func (r *userQueryRepo) FindByEmail(ctx context.Context, email string) (*biz.User, error) {
	db := r.data.GetReadDB()
	var user biz.User
	if err := db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Return nil, nil if not found (not an error)
		}
		r.log.WithContext(ctx).Errorf("Failed to find user by email: %v", err)
		return nil, err
	}
	return &user, nil
}

func (r *userQueryRepo) FindByUsername(ctx context.Context, username string) (*biz.User, error) {
	db := r.data.GetReadDB()
	var user biz.User
	if err := db.WithContext(ctx).Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.log.WithContext(ctx).Errorf("Failed to find user by username: %v", err)
		return nil, err
	}
	return &user, nil
}

func (r *userQueryRepo) List(ctx context.Context, filter *biz.UserListFilter) ([]*biz.User, int64, error) {
	db := r.data.GetReadDB()
	var users []*biz.User
	var total int64

	query := db.WithContext(ctx).Model(&biz.User{})

	// Apply filters
	if filter.Search != "" {
		searchPattern := fmt.Sprintf("%%%s%%", filter.Search)
		query = query.Where("email ILIKE ? OR username ILIKE ? OR full_name ILIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	if filter.Role != "" {
		query = query.Where("role = ?", filter.Role)
	}

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to count users: %v", err)
		return nil, 0, err
	}

	// Apply pagination
	page := int(filter.Page)
	if page < 1 {
		page = 1
	}
	pageSize := int(filter.PageSize)
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	// Fetch users
	if err := query.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&users).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to list users: %v", err)
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userQueryRepo) Count(ctx context.Context, filter *biz.UserListFilter) (int64, error) {
	db := r.data.GetReadDB()
	var total int64

	query := db.WithContext(ctx).Model(&biz.User{})

	if filter.Search != "" {
		searchPattern := fmt.Sprintf("%%%s%%", filter.Search)
		query = query.Where("email ILIKE ? OR username ILIKE ? OR full_name ILIKE ?",
			searchPattern, searchPattern, searchPattern)
	}

	if filter.Role != "" {
		query = query.Where("role = ?", filter.Role)
	}

	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to count users: %v", err)
		return 0, err
	}

	return total, nil
}

