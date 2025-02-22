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

type PrivateMessageDto struct {
	MessageDto
	ReceiverID uuid.UUID
}

type GroupMessageDto struct {
	MessageDto
	GroupID uuid.UUID
}

type chatService struct {
	userRepo           repos.UserRepoInterface
	privateChatRepo    repos.PrivateChatRepoInterface
	groupRepo          repos.GroupRepoInterface
	groupMsgRepo       repos.GroupMessageRepoInterface
	privateMessageRepo repos.PrivateMessageRepoInterface
}

type ChatServiceInterface interface {
	SendPrivateMsg(data PrivateMessageDto) error
	SendMsgToGroup(data GroupMessageDto) error
}

func NewChatService(
	userRepo repos.UserRepoInterface,
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
	ErrChat404      = errors.New("chat not found")
)

func (c *chatService) SendPrivateMsg(data PrivateMessageDto) error {
	chat, err := c.privateChatRepo.FindChat(data.ReceiverID, data.SenderID)
	if err == nil {
		pchat := models.PrivateMessage{
			BaseMessage: models.BaseMessage{
				Type:     data.Type,
				SenderID: data.SenderID,
				Content:  data.Content,
			},
			ChatID: chat.ID,
		}
		return c.privateMessageRepo.Create(nil, &pchat)
	}

	_, err = c.userRepo.FindByID(data.ReceiverID)
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

	tx.Commit()

	return nil
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
