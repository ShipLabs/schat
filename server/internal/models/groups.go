package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GroupRole string

const (
	Admin  GroupRole = "admin"
	Member GroupRole = "member"
)

type Group struct {
	gorm.Model  `json:"-"`
	ID          uuid.UUID `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	CreatorID   uuid.UUID `gorm:"not null" json:"creator_id"`
	Creator     User      `gorm:"foreignKey:creator_id" json:"-"`
	Description *string   `json:"description"`
}

type GroupMember struct {
	gorm.Model `json:"-"`
	UserID     uuid.UUID `gorm:"primaryKey;type:uuid" json:"user_id"`
	User       User      `gorm:"foreignKey:user_id" json:"-"`
	GroupID    uuid.UUID `gorm:"primaryKey;type:uuid" json:"group_id"`
	Group      Group     `gorm:"foreignKey:group_id" json:"-"`
	Role       GroupRole `gorm:"not null;default:member" json:"role"`
}
