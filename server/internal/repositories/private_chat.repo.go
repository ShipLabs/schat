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
	BeginDBTransaction() *gorm.DB
	CreatePrivateChat(txn *gorm.DB, chat *models.PrivateChat) error
	GetUserPrivateChats(userID uuid.UUID) ([]models.PrivateChat, error)
}

func NewPrivateChatRepo(db gorm.DB) PrivateChatRepoInterface {
	return &privateChatRepo{
		DB: db,
	}
}

func (p *privateChatRepo) BeginDBTransaction() *gorm.DB {
	return p.DB.Begin()
}

func (p *privateChatRepo) CreatePrivateChat(txn *gorm.DB, chat *models.PrivateChat) error {
	return txn.Create(chat).Error
}

func (p *privateChatRepo) GetUserPrivateChats(userID uuid.UUID) ([]models.PrivateChat, error) {
	var chats []models.PrivateChat
	err := p.DB.Where("first_member_id = ? OR second_member_id = ?", userID, userID).Find(&chats).Error
	return chats, err
}
