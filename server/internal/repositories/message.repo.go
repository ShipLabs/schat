package repos

import (
	"shiplabs/schat/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PrivateMessageRepoInterface interface {
	Create(txn *gorm.DB, message *models.PrivateMessage) error
	GetChatMessages(chatID uuid.UUID) ([]models.PrivateMessage, error)
}

type GroupMessageRepoInterface interface {
	Create(message *models.GroupMessage) error
	GetGroupMessages(groupID uuid.UUID) ([]models.GroupMessage, error)
}

type privateMessageRepo struct {
	DB gorm.DB
}

type groupMessageRepo struct {
	DB gorm.DB
}

func NewPrivateMessageRepo(db gorm.DB) PrivateMessageRepoInterface {
	return &privateMessageRepo{
		DB: db,
	}
}

func (p *privateMessageRepo) Create(txn *gorm.DB, message *models.PrivateMessage) error {
	if txn == nil {
		return p.DB.Create(message).Error
	}
	return txn.Create(message).Error
}

func (p *privateMessageRepo) GetChatMessages(chatID uuid.UUID) ([]models.PrivateMessage, error) {
	var messages []models.PrivateMessage
	err := p.DB.Where("chat_id=?", chatID).Find(&messages).Error
	return messages, err
}

func NewGroupMessageRepo(db gorm.DB) GroupMessageRepoInterface {
	return &groupMessageRepo{
		DB: db,
	}
}

func (g *groupMessageRepo) Create(message *models.GroupMessage) error {
	return g.DB.Create(&message).Error
}

func (g *groupMessageRepo) GetGroupMessages(groupID uuid.UUID) ([]models.GroupMessage, error) {
	var messages []models.GroupMessage
	err := g.DB.Where("group_id=?", groupID).Find(&messages).Error
	return messages, err
}
