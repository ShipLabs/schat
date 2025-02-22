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
	BeginDBTx() *gorm.DB
	CreatePrivateChat(txn *gorm.DB, chat *models.PrivateChat) error
	FindChat(mem1, mem2 uuid.UUID) (models.PrivateChat, error)
	GetUserPrivateChats(userID uuid.UUID) ([]models.PrivateChat, error)
}

func NewPrivateChatRepo(db gorm.DB) PrivateChatRepoInterface {
	return &privateChatRepo{
		DB: db,
	}
}

func (p *privateChatRepo) BeginDBTx() *gorm.DB {
	return p.DB.Begin()
}

func (p *privateChatRepo) CreatePrivateChat(txn *gorm.DB, chat *models.PrivateChat) error {
	if txn == nil {
		return p.DB.Create(chat).Error
	}
	return txn.Create(chat).Error
}

func (p *privateChatRepo) FindChat(mem1, mem2 uuid.UUID) (models.PrivateChat, error) {
	var chat models.PrivateChat
	err := p.DB.Where("(first_member_id = ? AND second_member_id = ?) OR (first_member_id = ? AND second_member_id = ?)", mem1, mem2, mem2, mem1).First(&chat).Error
	return chat, err
}

func (p *privateChatRepo) GetUserPrivateChats(userID uuid.UUID) ([]models.PrivateChat, error) {
	var chats []models.PrivateChat
	err := p.DB.Where("first_member_id=? OR second_member_id=?", userID, userID).Find(&chats).Error
	return chats, err
}
