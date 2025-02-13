package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PrivateChat struct {
	gorm.Model     `json:"-"`
	FirstMemberID  uuid.UUID `gorm:"not null;index" json:"first_member_id"`
	FirstMember    User      `gorm:"foreignKey:first_member_id" json:"-"`
	SecondMemberID uuid.UUID `gorm:"not null;index" json:"second_member_id"`
	SecondMember   User      `gorm:"foreignKey:second_member_id" json:"-"`
}
