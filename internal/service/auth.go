package service

import (
	"context"

	v1 "github.com/go-kratos/kratos-layout/api/auth/v1"
	"github.com/go-kratos/kratos-layout/internal/biz"
	"github.com/go-kratos/kratos-layout/internal/pkg/jwt"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/transport/http"
)

type AuthService struct {
	v1.UnimplementedAuthServiceServer

	uc *biz.AuthUsecase
}

func NewAuthService(uc *biz.AuthUsecase) *AuthService {
	return &AuthService{uc: uc}
}

// Login authenticates user and returns tokens
func (s *AuthService) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error) {
	// Extract IP and User-Agent from context
	ip := extractIPFromContext(ctx)
	userAgent := extractUserAgentFromContext(ctx)

	loginReq := &biz.LoginRequest{
		Email:    req.Identifier, // Identifier can be email or username
		Password: req.Password,
		IP:       ip,
		UserAgent: userAgent,
	}

	result, err := s.uc.Login(ctx, loginReq)
	if err != nil {
		return nil, err
	}

	return &v1.LoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
		TokenType:    result.TokenType,
		User:         toProtoAuthUser(result.User),
	}, nil
}

// Register creates new user and returns tokens
func (s *AuthService) Register(ctx context.Context, req *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	// Extract IP and User-Agent from context
	ip := extractIPFromContext(ctx)
	userAgent := extractUserAgentFromContext(ctx)

	registerReq := &biz.RegisterRequest{
		Email:     req.Email,
		Username:  req.Username,
		Password:  req.Password,
		FullName:  req.FullName,
		IP:        ip,
		UserAgent: userAgent,
	}

	result, err := s.uc.Register(ctx, registerReq)
	if err != nil {
		return nil, err
	}

	return &v1.RegisterResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
		TokenType:    result.TokenType,
		User:         toProtoAuthUser(result.User),
	}, nil
}

// RefreshToken generates new access token from refresh token
func (s *AuthService) RefreshToken(ctx context.Context, req *v1.RefreshTokenRequest) (*v1.RefreshTokenResponse, error) {
	result, err := s.uc.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &v1.RefreshTokenResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    result.ExpiresIn,
		TokenType:    result.TokenType,
	}, nil
}

// Logout revokes token
func (s *AuthService) Logout(ctx context.Context, req *v1.LogoutRequest) (*v1.LogoutResponse, error) {
	if err := s.uc.Logout(ctx, req.RefreshToken); err != nil {
		return nil, err
	}

	return &v1.LogoutResponse{Success: true}, nil
}

// GetCurrentUser gets current user from token
func (s *AuthService) GetCurrentUser(ctx context.Context, req *v1.GetCurrentUserRequest) (*v1.GetCurrentUserResponse, error) {
	// Extract token from Authorization header
	token, err := extractTokenFromContext(ctx)
	if err != nil {
		return nil, errors.Unauthorized("UNAUTHORIZED", "missing or invalid authorization header")
	}

	user, err := s.uc.GetUserFromToken(ctx, token)
	if err != nil {
		return nil, err
	}

	return &v1.GetCurrentUserResponse{
		User: toProtoAuthUser(user),
	}, nil
}

// RevokeAllTokens revokes all tokens for current user
func (s *AuthService) RevokeAllTokens(ctx context.Context, req *v1.RevokeAllTokensRequest) (*v1.RevokeAllTokensResponse, error) {
	// Extract token from Authorization header
	token, err := extractTokenFromContext(ctx)
	if err != nil {
		return nil, errors.Unauthorized("UNAUTHORIZED", "missing or invalid authorization header")
	}

	// Get user from token
	user, err := s.uc.GetUserFromToken(ctx, token)
	if err != nil {
		return nil, err
	}

	if err := s.uc.RevokeAllTokens(ctx, user.ID); err != nil {
		return nil, err
	}

	return &v1.RevokeAllTokensResponse{Success: true}, nil
}

// Helper functions

func toProtoAuthUser(user *biz.User) *v1.User {
	return &v1.User{
		Id:       user.ID.String(),
		Email:    user.Email,
		Username: user.Username,
		FullName: user.FullName,
		Role:     user.Role,
		Status:   user.Status,
	}
}

func extractIPFromContext(ctx context.Context) string {
	// Try to get from HTTP request context
	if req, ok := http.RequestFromServerContext(ctx); ok {
		return req.RemoteAddr
	}
	return ""
}

func extractUserAgentFromContext(ctx context.Context) string {
	// Try to get from HTTP request context
	if req, ok := http.RequestFromServerContext(ctx); ok {
		return req.UserAgent()
	}
	return ""
}

func extractTokenFromContext(ctx context.Context) (string, error) {
	// Try to get from HTTP request context
	if req, ok := http.RequestFromServerContext(ctx); ok {
		authHeader := req.Header.Get("Authorization")
		return jwt.ExtractTokenFromHeader(authHeader)
	}
	return "", errors.Unauthorized("UNAUTHORIZED", "missing authorization header")
}

