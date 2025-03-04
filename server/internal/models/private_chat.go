package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PrivateChat struct {
	gorm.Model     `json:"-"`
	ID             uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	FirstMemberID  uuid.UUID `gorm:"not null;index" json:"first_member_id"`
	FirstMember    User      `gorm:"foreignKey:first_member_id" json:"-"`
	SecondMemberID uuid.UUID `gorm:"not null;index" json:"second_member_id"`
	SecondMember   User      `gorm:"foreignKey:second_member_id" json:"-"`
	CreatedAt      time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt      time.Time `gorm:"not null" json:"updated_at"`
}
