package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ValidMsgType string

const (
	TEXT  ValidMsgType = "text"
	IMAGE ValidMsgType = "image"
	VIDEO ValidMsgType = "video"
	AUDIO ValidMsgType = "audio"
)

type BaseMessage struct {
	gorm.Model `json:"-"`
	ID         uuid.UUID    `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()" json:"id"`
	SenderID   uuid.UUID    `gorm:"not null;index" json:"sender_id"`
	Sender     User         `gorm:"foreignKey:sender_id" json:"-"`
	Type       ValidMsgType `gorm:"not null" json:"type"`
	Content    string       `gorm:"not null" json:"content"`
}

type Message struct {
	BaseMessage
	ReceiverID uuid.UUID `gorm:"not null;index" json:"receiver_id"`
	Receiver   User      `gorm:"foreignKey:receiver_id" json:"-"`
}

type GroupMessage struct {
	BaseMessage
	GroupID uuid.UUID `gorm:"not null;index" json:"group_id"`
	Group   Group     `gorm:"foreignKey:group_id" json:"-"`
}
