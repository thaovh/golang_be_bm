package biz

import (
	"context"

	v1 "github.com/go-kratos/kratos-layout/api/helloworld/v1"
	"github.com/gofrs/uuid/v5"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

var (
	// ErrUserNotFound is user not found.
	ErrUserNotFound = errors.NotFound(v1.ErrorReason_USER_NOT_FOUND.String(), "user not found")
)

// Greeter is a Greeter model with BaseEntity.
type Greeter struct {
	BaseEntity
	Hello   string `gorm:"type:varchar(255);not null" json:"hello"`
	Message string `gorm:"type:text" json:"message,omitempty"`
}

// GreeterCommandRepo is a repository interface for write operations (Commands)
type GreeterCommandRepo interface {
	Save(context.Context, *Greeter) (*Greeter, error)
	Update(context.Context, *Greeter) (*Greeter, error)
	Delete(context.Context, uuid.UUID) error
}

// GreeterQueryRepo is a repository interface for read operations (Queries)
type GreeterQueryRepo interface {
	FindByID(context.Context, uuid.UUID) (*Greeter, error)
	ListByHello(context.Context, string) ([]*Greeter, error)
	ListAll(context.Context) ([]*Greeter, error)
}

// GreeterUsecase is a Greeter usecase using CQRS pattern
type GreeterUsecase struct {
	commandRepo GreeterCommandRepo
	queryRepo   GreeterQueryRepo
	log         *log.Helper
}

// NewGreeterUsecase new a Greeter usecase with CQRS pattern
func NewGreeterUsecase(
	commandRepo GreeterCommandRepo,
	queryRepo GreeterQueryRepo,
	logger log.Logger,
) *GreeterUsecase {
	return &GreeterUsecase{
		commandRepo: commandRepo,
		queryRepo:   queryRepo,
		log:         log.NewHelper(logger),
	}
}

// CreateGreeter creates a Greeter (Command operation)
func (uc *GreeterUsecase) CreateGreeter(ctx context.Context, g *Greeter) (*Greeter, error) {
	uc.log.WithContext(ctx).Infof("CreateGreeter: %v", g.Hello)
	return uc.commandRepo.Save(ctx, g)
}

// UpdateGreeter updates a Greeter (Command operation)
func (uc *GreeterUsecase) UpdateGreeter(ctx context.Context, g *Greeter) (*Greeter, error) {
	uc.log.WithContext(ctx).Infof("UpdateGreeter: %v", g.Hello)
	return uc.commandRepo.Update(ctx, g)
}

// DeleteGreeter deletes a Greeter (Command operation)
func (uc *GreeterUsecase) DeleteGreeter(ctx context.Context, id uuid.UUID) error {
	uc.log.WithContext(ctx).Infof("DeleteGreeter: %s", id.String())
	return uc.commandRepo.Delete(ctx, id)
}

// GetGreeter gets a Greeter by ID (Query operation)
func (uc *GreeterUsecase) GetGreeter(ctx context.Context, id uuid.UUID) (*Greeter, error) {
	return uc.queryRepo.FindByID(ctx, id)
}

// ListGreetersByHello lists Greeters by Hello (Query operation)
func (uc *GreeterUsecase) ListGreetersByHello(ctx context.Context, hello string) ([]*Greeter, error) {
	return uc.queryRepo.ListByHello(ctx, hello)
}

// ListAllGreeters lists all Greeters (Query operation)
func (uc *GreeterUsecase) ListAllGreeters(ctx context.Context) ([]*Greeter, error) {
	return uc.queryRepo.ListAll(ctx)
}
