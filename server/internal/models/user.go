package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model `json:"-"`
	ID         uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Name       string    `gorm:"not null" json:"name"`
	Email      string    `gorm:"not null;uniqueIndex" json:"email"`
	Password   string    `gorm:"not null" json:"password"`
	CreatedAt  time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt  time.Time `gorm:"not null" json:"updated_at"`
}
