package biz

import (
	"context"
	"time"

	"github.com/gofrs/uuid/v5"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	ErrUserAlreadyExists = errors.Conflict("USER_ALREADY_EXISTS", "user already exists")
	ErrInvalidPassword   = errors.Unauthorized("INVALID_PASSWORD", "invalid password")
)

// User là domain model cho User
type User struct {
	BaseEntity

	// Thông tin cơ bản
	Email        string `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Username     string `gorm:"type:varchar(100);uniqueIndex;not null" json:"username"`
	PasswordHash string `gorm:"type:varchar(255);not null" json:"-"` // Ẩn trong JSON

	// Thông tin cá nhân
	FullName string `gorm:"type:varchar(200);index" json:"full_name,omitempty"`

	// Thông tin bổ sung
	DateOfBirth *time.Time `gorm:"type:date" json:"date_of_birth,omitempty"`
	Gender      string     `gorm:"type:varchar(20)" json:"gender,omitempty"` // male, female, other

	// Last login tracking
	LastLoginAt *time.Time `gorm:"type:timestamp;index" json:"last_login_at,omitempty"`
	LastLoginIP string     `gorm:"type:varchar(45)" json:"last_login_ip,omitempty"`

	// Role và permissions
	Role string `gorm:"type:varchar(50);default:'user';index" json:"role"` // user, admin, moderator
}

// GetDisplayName returns display name (FullName or Username)
func (u *User) GetDisplayName() string {
	if u.FullName != "" {
		return u.FullName
	}
	return u.Username
}

// IsAdmin checks if user is admin
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

// UserCommandRepo là repository interface cho write operations
type UserCommandRepo interface {
	Save(context.Context, *User) (*User, error)
	Update(context.Context, *User) (*User, error)
	Delete(context.Context, uuid.UUID) error
	UpdatePassword(context.Context, uuid.UUID, string) error
	UpdateLastLogin(context.Context, uuid.UUID, string) error
}

// UserQueryRepo là repository interface cho read operations
type UserQueryRepo interface {
	FindByID(context.Context, uuid.UUID) (*User, error)
	FindByEmail(context.Context, string) (*User, error)
	FindByUsername(context.Context, string) (*User, error)
	List(context.Context, *UserListFilter) ([]*User, int64, error)
	Count(context.Context, *UserListFilter) (int64, error)
}

// UserListFilter cho pagination và filtering
type UserListFilter struct {
	Page     int32
	PageSize int32
	Search   string // Search by email, username, full_name
	Role     string // Filter by role
	Status   string // Filter by status
}

// UserUsecase là usecase cho User với CQRS pattern
type UserUsecase struct {
	commandRepo UserCommandRepo
	queryRepo   UserQueryRepo
	log         *log.Helper
}

// NewUserUsecase tạo UserUsecase mới
func NewUserUsecase(
	commandRepo UserCommandRepo,
	queryRepo UserQueryRepo,
	logger log.Logger,
) *UserUsecase {
	return &UserUsecase{
		commandRepo: commandRepo,
		queryRepo:   queryRepo,
		log:         log.NewHelper(logger),
	}
}

// CreateUser creates a new user (Command)
func (uc *UserUsecase) CreateUser(ctx context.Context, user *User) (*User, error) {
	uc.log.WithContext(ctx).Infof("CreateUser: %s", user.Email)

	// Set audit fields from context
	user.SetAuditFields(ctx, true)

	// Check if user already exists
	existing, _ := uc.queryRepo.FindByEmail(ctx, user.Email)
	if existing != nil {
		return nil, ErrUserAlreadyExists
	}

	existing, _ = uc.queryRepo.FindByUsername(ctx, user.Username)
	if existing != nil {
		return nil, ErrUserAlreadyExists
	}

	return uc.commandRepo.Save(ctx, user)
}

// UpdateUser updates a user (Command)
func (uc *UserUsecase) UpdateUser(ctx context.Context, user *User) (*User, error) {
	uc.log.WithContext(ctx).Infof("UpdateUser: %s", user.ID.String())
	
	// Set audit fields from context
	user.SetAuditFields(ctx, false)
	
	return uc.commandRepo.Update(ctx, user)
}

// DeleteUser deletes a user (Command)
func (uc *UserUsecase) DeleteUser(ctx context.Context, id uuid.UUID) error {
	uc.log.WithContext(ctx).Infof("DeleteUser: %s", id.String())
	return uc.commandRepo.Delete(ctx, id)
}

// ChangePassword changes user password (Command)
func (uc *UserUsecase) ChangePassword(ctx context.Context, id uuid.UUID, newPasswordHash string) error {
	uc.log.WithContext(ctx).Infof("ChangePassword: %s", id.String())
	return uc.commandRepo.UpdatePassword(ctx, id, newPasswordHash)
}

// UpdateLastLogin updates last login info (Command)
func (uc *UserUsecase) UpdateLastLogin(ctx context.Context, id uuid.UUID, ip string) error {
	return uc.commandRepo.UpdateLastLogin(ctx, id, ip)
}

// GetUser gets a user by ID (Query)
func (uc *UserUsecase) GetUser(ctx context.Context, id uuid.UUID) (*User, error) {
	return uc.queryRepo.FindByID(ctx, id)
}

// GetUserByEmail gets a user by email (Query)
func (uc *UserUsecase) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	return uc.queryRepo.FindByEmail(ctx, email)
}

// GetUserByUsername gets a user by username (Query)
func (uc *UserUsecase) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	return uc.queryRepo.FindByUsername(ctx, username)
}

// ListUsers lists users with filter (Query)
func (uc *UserUsecase) ListUsers(ctx context.Context, filter *UserListFilter) ([]*User, int64, error) {
	return uc.queryRepo.List(ctx, filter)
}

