package repos

import (
	"shiplabs/schat/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type privateChatRepo struct {
	DB gorm.DB
}

type PrivateChatRepoInterface interface {
	CreatePrivateChat(chat *models.PrivateChat) error
	GetUserPrivateChats(userID uuid.UUID) ([]models.PrivateChat, error)
}

func NewPrivateChatRepo(db gorm.DB) PrivateChatRepoInterface {
	return &privateChatRepo{
		DB: db,
	}
}

func (p *privateChatRepo) CreatePrivateChat(chat *models.PrivateChat) error {
	return p.DB.Create(chat).Error
}

func (p *privateChatRepo) GetUserPrivateChats(userID uuid.UUID) ([]models.PrivateChat, error) {
	var chats []models.PrivateChat
	err := p.DB.Where("first_member_id = ? OR second_member_id = ?", userID, userID).Find(&chats).Error
	return chats, err
}
