package service

import (
	"context"
	"time"

	v1 "github.com/go-kratos/kratos-layout/api/user/v1"
	"github.com/go-kratos/kratos-layout/internal/biz"
	"github.com/go-kratos/kratos-layout/internal/pkg/password"
	"github.com/go-kratos/kratos-layout/internal/pkg/validator"
	"github.com/gofrs/uuid/v5"

	"github.com/go-kratos/kratos/v2/errors"
)

type UserService struct {
	v1.UnimplementedUserServiceServer

	uc *biz.UserUsecase
}

func NewUserService(uc *biz.UserUsecase) *UserService {
	return &UserService{uc: uc}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, req *v1.CreateUserRequest) (*v1.CreateUserResponse, error) {
	// Validate email
	if err := validator.ValidateEmail(req.Email); err != nil {
		return nil, errors.BadRequest("INVALID_EMAIL", err.Error())
	}

	// Validate username
	if err := validator.ValidateUsername(req.Username); err != nil {
		return nil, errors.BadRequest("INVALID_USERNAME", err.Error())
	}

	// Validate password
	if err := validator.ValidatePassword(req.Password); err != nil {
		return nil, errors.BadRequest("INVALID_PASSWORD", err.Error())
	}

	// Hash password before saving
	passwordHash, err := password.Hash(req.Password)
	if err != nil {
		return nil, errors.InternalServer("PASSWORD_HASH_ERROR", "failed to hash password")
	}

	user := &biz.User{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: passwordHash,
		FullName:     req.FullName,
		Gender:       req.Gender,
		Role:         req.Role,
	}

	if req.Role == "" {
		user.Role = "user"
	}

	if req.DateOfBirth != "" {
		dob, err := time.Parse("2006-01-02", req.DateOfBirth)
		if err == nil {
			user.DateOfBirth = &dob
		}
	}

	createdUser, err := s.uc.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &v1.CreateUserResponse{
		User: toProtoUser(createdUser),
	}, nil
}

// UpdateUser updates a user
func (s *UserService) UpdateUser(ctx context.Context, req *v1.UpdateUserRequest) (*v1.UpdateUserResponse, error) {
	id, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, errors.BadRequest("INVALID_ID", "invalid user id")
	}

	user, err := s.uc.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	user.FullName = req.FullName
	user.Gender = req.Gender

	if req.DateOfBirth != "" {
		dob, err := time.Parse("2006-01-02", req.DateOfBirth)
		if err == nil {
			user.DateOfBirth = &dob
		}
	}

	updatedUser, err := s.uc.UpdateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return &v1.UpdateUserResponse{
		User: toProtoUser(updatedUser),
	}, nil
}

// DeleteUser deletes a user
func (s *UserService) DeleteUser(ctx context.Context, req *v1.DeleteUserRequest) (*v1.DeleteUserResponse, error) {
	id, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, errors.BadRequest("INVALID_ID", "invalid user id")
	}

	if err := s.uc.DeleteUser(ctx, id); err != nil {
		return nil, err
	}

	return &v1.DeleteUserResponse{Success: true}, nil
}

// ChangePassword changes user password
func (s *UserService) ChangePassword(ctx context.Context, req *v1.ChangePasswordRequest) (*v1.ChangePasswordResponse, error) {
	id, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, errors.BadRequest("INVALID_ID", "invalid user id")
	}

	// Validate new password
	if err := validator.ValidatePassword(req.NewPassword); err != nil {
		return nil, errors.BadRequest("INVALID_PASSWORD", err.Error())
	}

	// Hash password before updating
	passwordHash, err := password.Hash(req.NewPassword)
	if err != nil {
		return nil, errors.InternalServer("PASSWORD_HASH_ERROR", "failed to hash password")
	}

	if err := s.uc.ChangePassword(ctx, id, passwordHash); err != nil {
		return nil, err
	}

	return &v1.ChangePasswordResponse{Success: true}, nil
}

// GetUser gets a user by ID
func (s *UserService) GetUser(ctx context.Context, req *v1.GetUserRequest) (*v1.GetUserResponse, error) {
	id, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, errors.BadRequest("INVALID_ID", "invalid user id")
	}

	user, err := s.uc.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}

	return &v1.GetUserResponse{
		User: toProtoUser(user),
	}, nil
}

// GetUserByEmail gets a user by email
func (s *UserService) GetUserByEmail(ctx context.Context, req *v1.GetUserByEmailRequest) (*v1.GetUserByEmailResponse, error) {
	user, err := s.uc.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	return &v1.GetUserByEmailResponse{
		User: toProtoUser(user),
	}, nil
}

// GetUserByUsername gets a user by username
func (s *UserService) GetUserByUsername(ctx context.Context, req *v1.GetUserByUsernameRequest) (*v1.GetUserByUsernameResponse, error) {
	user, err := s.uc.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	return &v1.GetUserByUsernameResponse{
		User: toProtoUser(user),
	}, nil
}

// ListUsers lists users with filter
func (s *UserService) ListUsers(ctx context.Context, req *v1.ListUsersRequest) (*v1.ListUsersResponse, error) {
	filter := &biz.UserListFilter{
		Page:     req.Page,
		PageSize: req.PageSize,
		Search:   req.Search,
		Role:     req.Role,
		Status:   req.Status,
	}

	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.PageSize < 1 {
		filter.PageSize = 10
	}

	users, total, err := s.uc.ListUsers(ctx, filter)
	if err != nil {
		return nil, err
	}

	protoUsers := make([]*v1.User, len(users))
	for i, user := range users {
		protoUsers[i] = toProtoUser(user)
	}

	return &v1.ListUsersResponse{
		Users:    protoUsers,
		Total:    int32(total),
		Page:     filter.Page,
		PageSize: filter.PageSize,
	}, nil
}

// toProtoUser converts biz.User to proto User
func toProtoUser(user *biz.User) *v1.User {
	protoUser := &v1.User{
		Id:        user.ID.String(),
		Email:     user.Email,
		Username:  user.Username,
		FullName:  user.FullName,
		Gender:    user.Gender,
		Role:      user.Role,
		Status:    user.Status,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
	}

	if user.DateOfBirth != nil {
		protoUser.DateOfBirth = user.DateOfBirth.Format("2006-01-02")
	}

	if user.LastLoginAt != nil {
		protoUser.LastLoginAt = user.LastLoginAt.Format(time.RFC3339)
	}

	if user.LastLoginIP != "" {
		protoUser.LastLoginIp = user.LastLoginIP
	}

	return protoUser
}

