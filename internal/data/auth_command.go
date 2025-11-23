package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos-layout/internal/biz"
	"github.com/gofrs/uuid/v5"

	"github.com/go-kratos/kratos/v2/log"
)

type authCommandRepo struct {
	data *Data
	log  *log.Helper
}

func NewAuthCommandRepo(data *Data, logger log.Logger) biz.AuthCommandRepo {
	return &authCommandRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *authCommandRepo) SaveToken(ctx context.Context, token *biz.AuthToken) (*biz.AuthToken, error) {
	db := r.data.GetWriteDB()
	
	// If token has ID, update it; otherwise create new
	if token.ID != uuid.Nil {
		if err := db.WithContext(ctx).Save(token).Error; err != nil {
			r.log.WithContext(ctx).Errorf("Failed to update token: %v", err)
			return nil, err
		}
	} else {
		if err := db.WithContext(ctx).Create(token).Error; err != nil {
			r.log.WithContext(ctx).Errorf("Failed to save token: %v", err)
			return nil, err
		}
	}
	
	return token, nil
}

func (r *authCommandRepo) RevokeToken(ctx context.Context, refreshToken string) error {
	db := r.data.GetWriteDB()
	now := time.Now()
	
	result := db.WithContext(ctx).Model(&biz.AuthToken{}).
		Where("refresh_token = ? AND revoked = ?", refreshToken, false).
		Updates(map[string]interface{}{
			"revoked":    true,
			"revoked_at": &now,
		})
	
	if result.Error != nil {
		r.log.WithContext(ctx).Errorf("Failed to revoke token: %v", result.Error)
		return result.Error
	}
	
	if result.RowsAffected == 0 {
		r.log.WithContext(ctx).Warnf("Token not found or already revoked: %s", refreshToken)
	}
	
	return nil
}

func (r *authCommandRepo) RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error {
	db := r.data.GetWriteDB()
	now := time.Now()
	
	result := db.WithContext(ctx).Model(&biz.AuthToken{}).
		Where("user_id = ? AND revoked = ?", userID, false).
		Updates(map[string]interface{}{
			"revoked":    true,
			"revoked_at": &now,
		})
	
	if result.Error != nil {
		r.log.WithContext(ctx).Errorf("Failed to revoke all tokens: %v", result.Error)
		return result.Error
	}
	
	r.log.WithContext(ctx).Infof("Revoked %d tokens for user %s", result.RowsAffected, userID.String())
	
	return nil
}

