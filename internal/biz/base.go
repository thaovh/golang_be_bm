package biz

import (
	"time"

	"github.com/gofrs/uuid/v5"
	"gorm.io/gorm"
)

// BaseEntity là base entity cho tất cả các domain models với UUID v7
type BaseEntity struct {
	ID        uuid.UUID      `gorm:"type:uuid;primarykey" json:"id"`
	CreatedAt time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Audit fields
	CreatedBy *uuid.UUID `gorm:"type:uuid;index" json:"created_by,omitempty"`
	UpdatedBy *uuid.UUID `gorm:"type:uuid;index" json:"updated_by,omitempty"`

	// Optimistic locking
	Version int `gorm:"default:1" json:"version"`

	// Status: active, inactive, archived, deleted
	Status string `gorm:"type:varchar(20);default:'active';index" json:"status"`
}

// BeforeCreate hook để tự động generate UUID v7 trước khi tạo
func (b *BaseEntity) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		// Generate UUID v7 (time-ordered)
		id, err := uuid.NewV7()
		if err != nil {
			return err
		}
		b.ID = id
	}

	// Set default status nếu chưa có
	if b.Status == "" {
		b.Status = "active"
	}

	// Set default version
	if b.Version == 0 {
		b.Version = 1
	}

	return nil
}

// BeforeUpdate hook để tự động tăng version
func (b *BaseEntity) BeforeUpdate(tx *gorm.DB) error {
	b.Version++
	return nil
}

// IsActive checks if entity is active
func (b *BaseEntity) IsActive() bool {
	return b.Status == "active"
}

// IsDeleted checks if entity is soft deleted
func (b *BaseEntity) IsDeleted() bool {
	return b.DeletedAt.Valid
}

