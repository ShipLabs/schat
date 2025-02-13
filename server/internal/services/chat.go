package services

import (
	"errors"
	"shiplabs/schat/internal/models"
	repos "shiplabs/schat/internal/repositories"

	"github.com/google/uuid"
)

type MessageDto struct {
	SenderID uuid.UUID
	Type     models.ValidMsgType
	Content  string
}

type CreateChatDto struct {
	MessageDto
	ReceiverID uuid.UUID
}

type PrivateMessageDto struct {
	MessageDto
	ChatID uuid.UUID
}

type GroupMessageDto struct {
	MessageDto
	GroupID uuid.UUID
}

type chatService struct {
	userRepo           repos.UserRepo
	privateChatRepo    repos.PrivateChatRepoInterface
	groupRepo          repos.GroupRepoInterface
	groupMsgRepo       repos.GroupMessageRepoInterface
	privateMessageRepo repos.PrivateMessageRepoInterface
}

type ChatServiceInterface interface {
	CreatePrivateChat(data CreateChatDto) error
	SendPrivateMsg(data PrivateMessageDto) error
	SendMsgToGroup(data GroupMessageDto) error
}

func NewChatService(
	userRepo repos.UserRepo,
	privateChatRepo repos.PrivateChatRepoInterface,
	groupRepo repos.GroupRepoInterface,
	groupMsgRepo repos.GroupMessageRepoInterface,
	privateMessageRepo repos.PrivateMessageRepoInterface,
) ChatServiceInterface {
	return &chatService{
		userRepo:           userRepo,
		privateChatRepo:    privateChatRepo,
		groupRepo:          groupRepo,
		groupMsgRepo:       groupMsgRepo,
		privateMessageRepo: privateMessageRepo,
	}
}

var (
	ErrUserNotFound = errors.New("user not found")
	ErrCreatingChat = errors.New("error creating chat")
)

func (c *chatService) CreatePrivateChat(data CreateChatDto) error {
	_, err := c.userRepo.FindByID(data.ReceiverID)
	if err != nil {
		return ErrUserNotFound
	}

	privateChat := &models.PrivateChat{
		FirstMemberID:  data.SenderID,
		SecondMemberID: data.ReceiverID,
	}

	tx := c.privateChatRepo.BeginDBTx()
	if err := c.privateChatRepo.CreatePrivateChat(tx, privateChat); err != nil {
		tx.Rollback()
		return ErrCreatingChat
	}

	privateMessage := &models.PrivateMessage{
		ChatID: privateChat.ID,
		BaseMessage: models.BaseMessage{
			Type:     data.Type,
			SenderID: data.SenderID,
			Content:  data.Content,
		},
	}
	if err := c.privateMessageRepo.Create(tx, privateMessage); err != nil {
		tx.Rollback()
		return ErrCreatingChat
	}

	return nil
}

func (c *chatService) SendPrivateMsg(data PrivateMessageDto) error {
	_, err := c.privateChatRepo.FindByID(data.ChatID)
	if err != nil {
		return err
	}
	pchat := models.PrivateMessage{
		ChatID: data.ChatID,
		BaseMessage: models.BaseMessage{
			Type:     data.Type,
			SenderID: data.SenderID,
			Content:  data.Content,
		},
	}
	return c.privateMessageRepo.Create(nil, &pchat)
}

func (c *chatService) SendMsgToGroup(data GroupMessageDto) error {
	_, err := c.groupRepo.FindByID(data.GroupID)
	if err != nil {
		return err
	}

	msg := &models.GroupMessage{
		BaseMessage: models.BaseMessage{
			Type:     data.Type,
			SenderID: data.SenderID,
		},
		GroupID: data.GroupID,
	}

	return c.groupMsgRepo.Create(msg)
}
