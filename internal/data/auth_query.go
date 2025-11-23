package data

import (
	"context"

	"github.com/go-kratos/kratos-layout/internal/biz"
	"github.com/gofrs/uuid/v5"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

type authQueryRepo struct {
	data *Data
	log  *log.Helper
}

func NewAuthQueryRepo(data *Data, logger log.Logger) biz.AuthQueryRepo {
	return &authQueryRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *authQueryRepo) FindTokenByRefreshToken(ctx context.Context, refreshToken string) (*biz.AuthToken, error) {
	db := r.data.GetReadDB()
	var token biz.AuthToken
	
	if err := db.WithContext(ctx).
		Where("refresh_token = ?", refreshToken).
		First(&token).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.log.WithContext(ctx).Errorf("Failed to find token by refresh token: %v", err)
		return nil, err
	}
	
	return &token, nil
}

func (r *authQueryRepo) FindTokenByAccessToken(ctx context.Context, accessToken string) (*biz.AuthToken, error) {
	db := r.data.GetReadDB()
	var token biz.AuthToken
	
	if err := db.WithContext(ctx).
		Where("token = ? AND revoked = ?", accessToken, false).
		First(&token).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.log.WithContext(ctx).Errorf("Failed to find token by access token: %v", err)
		return nil, err
	}
	
	return &token, nil
}

func (r *authQueryRepo) ListUserTokens(ctx context.Context, userID uuid.UUID) ([]*biz.AuthToken, error) {
	db := r.data.GetReadDB()
	var tokens []*biz.AuthToken
	
	if err := db.WithContext(ctx).
		Where("user_id = ? AND revoked = ?", userID, false).
		Order("created_at DESC").
		Find(&tokens).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to list user tokens: %v", err)
		return nil, err
	}
	
	return tokens, nil
}

