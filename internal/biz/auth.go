package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos-layout/internal/conf"
	"github.com/go-kratos/kratos-layout/internal/pkg/jwt"
	"github.com/go-kratos/kratos-layout/internal/pkg/password"
	"github.com/gofrs/uuid/v5"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	ErrInvalidCredentials = errors.Unauthorized("INVALID_CREDENTIALS", "invalid email or password")
	ErrTokenExpired       = errors.Unauthorized("TOKEN_EXPIRED", "token has expired")
	ErrTokenInvalid       = errors.Unauthorized("TOKEN_INVALID", "invalid token")
	ErrTokenRevoked       = errors.Unauthorized("TOKEN_REVOKED", "token has been revoked")
)

// AuthToken represents authentication token stored in database
type AuthToken struct {
	BaseEntity

	UserID            uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Token             string    `gorm:"type:text;not null" json:"token"` // Access token (for reference)
	RefreshToken      string    `gorm:"type:text;not null;uniqueIndex" json:"refresh_token"`
	ExpiresAt         time.Time `gorm:"not null;index" json:"expires_at"`
	RefreshExpiresAt  time.Time `gorm:"not null;index" json:"refresh_expires_at"`
	IPAddress         string    `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent         string    `gorm:"type:text" json:"user_agent"`
	Revoked           bool      `gorm:"default:false;index" json:"revoked"`
	RevokedAt         *time.Time `gorm:"type:timestamp" json:"revoked_at,omitempty"`
}

// LoginRequest for authentication
type LoginRequest struct {
	Email    string // Email or username
	Password string
	IP       string
	UserAgent string
}

// LoginResponse with tokens
type LoginResponse struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
	TokenType    string // "Bearer"
	User         *User
}

// RegisterRequest for user registration
type RegisterRequest struct {
	Email    string
	Username string
	Password string
	FullName string
	IP       string
	UserAgent string
}

// AuthCommandRepo for write operations
type AuthCommandRepo interface {
	SaveToken(context.Context, *AuthToken) (*AuthToken, error)
	RevokeToken(context.Context, string) error // Revoke by refresh token
	RevokeAllUserTokens(context.Context, uuid.UUID) error
}

// AuthQueryRepo for read operations
type AuthQueryRepo interface {
	FindTokenByRefreshToken(context.Context, string) (*AuthToken, error)
	FindTokenByAccessToken(context.Context, string) (*AuthToken, error)
	ListUserTokens(context.Context, uuid.UUID) ([]*AuthToken, error)
}

// AuthUsecase handles authentication logic
type AuthUsecase struct {
	userQueryRepo   UserQueryRepo
	userCommandRepo UserCommandRepo
	authCommandRepo AuthCommandRepo
	authQueryRepo   AuthQueryRepo
	jwtSecret       []byte
	accessExpiry    time.Duration
	refreshExpiry   time.Duration
	log             *log.Helper
}

// NewAuthUsecase creates a new AuthUsecase
func NewAuthUsecase(
	userQueryRepo UserQueryRepo,
	userCommandRepo UserCommandRepo,
	authCommandRepo AuthCommandRepo,
	authQueryRepo AuthQueryRepo,
	authConfig *AuthConfig,
	logger log.Logger,
) *AuthUsecase {
	return &AuthUsecase{
		userQueryRepo:   userQueryRepo,
		userCommandRepo: userCommandRepo,
		authCommandRepo: authCommandRepo,
		authQueryRepo:   authQueryRepo,
		jwtSecret:       []byte(authConfig.JwtSecret),
		accessExpiry:    time.Duration(authConfig.AccessExpiry) * time.Second,
		refreshExpiry:   time.Duration(authConfig.RefreshExpiry) * time.Second,
		log:             log.NewHelper(logger),
	}
}

// AuthConfig wraps auth configuration
type AuthConfig struct {
	JwtSecret      string
	AccessExpiry   int64
	RefreshExpiry  int64
}

// NewAuthConfigFromConf creates AuthConfig from conf.Auth
func NewAuthConfigFromConf(auth *conf.Auth) *AuthConfig {
	if auth == nil {
		// Default values
		return &AuthConfig{
			JwtSecret:     "default-secret-key-change-in-production",
			AccessExpiry:  3600,   // 1 hour
			RefreshExpiry: 604800, // 7 days
		}
	}
	return &AuthConfig{
		JwtSecret:     auth.JwtSecret,
		AccessExpiry:   auth.AccessTokenExpiry,
		RefreshExpiry: auth.RefreshTokenExpiry,
	}
}

// Login authenticates user and returns tokens
func (uc *AuthUsecase) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	uc.log.WithContext(ctx).Infof("Login attempt: %s", req.Email)

	// Find user by email or username
	user, err := uc.userQueryRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	// If not found by email, try username
	if user == nil {
		user, err = uc.userQueryRepo.FindByUsername(ctx, req.Email)
		if err != nil {
			return nil, err
		}
	}

	if user == nil {
		return nil, ErrInvalidCredentials
	}

	// Verify password
	if !password.Verify(req.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	// Check if user is active
	if !user.IsActive() {
		return nil, errors.Forbidden("USER_INACTIVE", "user account is inactive")
	}

	// Generate tokens
	accessToken, err := jwt.GenerateAccessToken(user.ID, user.Email, user.Role, uc.jwtSecret, uc.accessExpiry)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("Failed to generate access token: %v", err)
		return nil, errors.InternalServer("TOKEN_GENERATION_ERROR", "failed to generate token")
	}

	refreshToken, err := jwt.GenerateRefreshToken()
	if err != nil {
		uc.log.WithContext(ctx).Errorf("Failed to generate refresh token: %v", err)
		return nil, errors.InternalServer("TOKEN_GENERATION_ERROR", "failed to generate refresh token")
	}

	// Save refresh token to database
	now := time.Now()
	authToken := &AuthToken{
		UserID:           user.ID,
		Token:            accessToken, // Store for reference
		RefreshToken:     refreshToken,
		ExpiresAt:        now.Add(uc.accessExpiry),
		RefreshExpiresAt: now.Add(uc.refreshExpiry),
		IPAddress:        req.IP,
		UserAgent:        req.UserAgent,
	}
	
	// Set audit fields - user is creating their own token
	authToken.SetAuditFields(ctx, true)
	// If no user in context (public login), set created_by to the user themselves
	if authToken.CreatedBy == nil {
		authToken.CreatedBy = &user.ID
		authToken.UpdatedBy = &user.ID
	}

	savedToken, err := uc.authCommandRepo.SaveToken(ctx, authToken)
	if err != nil {
		uc.log.WithContext(ctx).Errorf("Failed to save token: %v", err)
		return nil, errors.InternalServer("TOKEN_SAVE_ERROR", "failed to save token")
	}

	// Update last login
	if err := uc.userCommandRepo.UpdateLastLogin(ctx, user.ID, req.IP); err != nil {
		uc.log.WithContext(ctx).Warnf("Failed to update last login: %v", err)
		// Don't fail login if this fails
	}

	uc.log.WithContext(ctx).Infof("User logged in successfully: %s", user.Email)

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: savedToken.RefreshToken,
		ExpiresIn:    int64(uc.accessExpiry.Seconds()),
		TokenType:    "Bearer",
		User:         user,
	}, nil
}

// Register creates new user and returns tokens
func (uc *AuthUsecase) Register(ctx context.Context, req *RegisterRequest) (*LoginResponse, error) {
	uc.log.WithContext(ctx).Infof("Register attempt: %s", req.Email)

	// Hash password
	passwordHash, err := password.Hash(req.Password)
	if err != nil {
		return nil, errors.InternalServer("PASSWORD_HASH_ERROR", "failed to hash password")
	}

	// Create user
	user := &User{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: passwordHash,
		FullName:     req.FullName,
		Role:         "user", // Default role
	}
	
	// Try to set audit fields from context (if authenticated user is creating)
	user.SetAuditFields(ctx, true)

	createdUser, err := uc.userCommandRepo.Save(ctx, user)
	if err != nil {
		return nil, err
	}
	
	// If no created_by was set (public register), set it to the user themselves
	if createdUser.CreatedBy == nil {
		createdUser.CreatedBy = &createdUser.ID
		createdUser.UpdatedBy = &createdUser.ID
		// Update the user to save audit fields
		_, err = uc.userCommandRepo.Update(ctx, createdUser)
		if err != nil {
			uc.log.WithContext(ctx).Warnf("Failed to update user audit fields: %v", err)
			// Don't fail registration if this fails
		}
	}

	// Generate tokens
	accessToken, err := jwt.GenerateAccessToken(createdUser.ID, createdUser.Email, createdUser.Role, uc.jwtSecret, uc.accessExpiry)
	if err != nil {
		return nil, errors.InternalServer("TOKEN_GENERATION_ERROR", "failed to generate token")
	}

	refreshToken, err := jwt.GenerateRefreshToken()
	if err != nil {
		return nil, errors.InternalServer("TOKEN_GENERATION_ERROR", "failed to generate refresh token")
	}

	// Save refresh token
	now := time.Now()
	authToken := &AuthToken{
		UserID:           createdUser.ID,
		Token:            accessToken,
		RefreshToken:     refreshToken,
		ExpiresAt:        now.Add(uc.accessExpiry),
		RefreshExpiresAt: now.Add(uc.refreshExpiry),
		IPAddress:        req.IP,
		UserAgent:        req.UserAgent,
	}
	
	// Set audit fields - user is creating their own token
	authToken.SetAuditFields(ctx, true)
	// If no user in context (public register), set created_by to the user themselves
	if authToken.CreatedBy == nil {
		authToken.CreatedBy = &createdUser.ID
		authToken.UpdatedBy = &createdUser.ID
	}

	savedToken, err := uc.authCommandRepo.SaveToken(ctx, authToken)
	if err != nil {
		return nil, errors.InternalServer("TOKEN_SAVE_ERROR", "failed to save token")
	}

	uc.log.WithContext(ctx).Infof("User registered successfully: %s", createdUser.Email)

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: savedToken.RefreshToken,
		ExpiresIn:    int64(uc.accessExpiry.Seconds()),
		TokenType:    "Bearer",
		User:         createdUser,
	}, nil
}

// RefreshToken generates new access token from refresh token
func (uc *AuthUsecase) RefreshToken(ctx context.Context, refreshToken string) (*LoginResponse, error) {
	uc.log.WithContext(ctx).Info("Refresh token request")

	// Find token in database
	token, err := uc.authQueryRepo.FindTokenByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, err
	}

	if token == nil {
		return nil, ErrTokenInvalid
	}

	// Check if token is revoked
	if token.Revoked {
		return nil, ErrTokenRevoked
	}

	// Check if refresh token is expired
	if time.Now().After(token.RefreshExpiresAt) {
		return nil, ErrTokenExpired
	}

	// Get user
	user, err := uc.userQueryRepo.FindByID(ctx, token.UserID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	// Check if user is active
	if !user.IsActive() {
		return nil, errors.Forbidden("USER_INACTIVE", "user account is inactive")
	}

	// Generate new access token
	accessToken, err := jwt.GenerateAccessToken(user.ID, user.Email, user.Role, uc.jwtSecret, uc.accessExpiry)
	if err != nil {
		return nil, errors.InternalServer("TOKEN_GENERATION_ERROR", "failed to generate token")
	}

	// Update token in database
	token.Token = accessToken
	token.ExpiresAt = time.Now().Add(uc.accessExpiry)
	
	// Set audit fields from context
	token.SetAuditFields(ctx, false)
	// If no user in context, set updated_by to the token's user
	if token.UpdatedBy == nil {
		token.UpdatedBy = &token.UserID
	}
	
	_, err = uc.authCommandRepo.SaveToken(ctx, token)
	if err != nil {
		return nil, errors.InternalServer("TOKEN_SAVE_ERROR", "failed to update token")
	}

	return &LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: token.RefreshToken, // Keep same refresh token
		ExpiresIn:    int64(uc.accessExpiry.Seconds()),
		TokenType:    "Bearer",
		User:         user,
	}, nil
}

// Logout revokes token
func (uc *AuthUsecase) Logout(ctx context.Context, refreshToken string) error {
	uc.log.WithContext(ctx).Info("Logout request")

	if err := uc.authCommandRepo.RevokeToken(ctx, refreshToken); err != nil {
		return err
	}

	return nil
}

// GetUserFromToken extracts user from JWT token
func (uc *AuthUsecase) GetUserFromToken(ctx context.Context, tokenString string) (*User, error) {
	// Validate token
	claims, err := jwt.ValidateToken(tokenString, uc.jwtSecret)
	if err != nil {
		if err == jwt.ErrExpiredToken {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	// Get user from database
	user, err := uc.userQueryRepo.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, ErrUserNotFound
	}

	// Check if user is active
	if !user.IsActive() {
		return nil, errors.Forbidden("USER_INACTIVE", "user account is inactive")
	}

	return user, nil
}

// RevokeAllTokens revokes all tokens for a user
func (uc *AuthUsecase) RevokeAllTokens(ctx context.Context, userID uuid.UUID) error {
	uc.log.WithContext(ctx).Infof("Revoke all tokens for user: %s", userID.String())
	return uc.authCommandRepo.RevokeAllUserTokens(ctx, userID)
}

