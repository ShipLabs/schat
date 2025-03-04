package services

import (
	"errors"
	"log"
	"shiplabs/schat/internal/models"
	repos "shiplabs/schat/internal/repositories"

	"github.com/google/uuid"
)

type MessageDto struct {
	Type    models.ValidMsgType `json:"type"`
	Content string              `json:"content"`
}

type PrivateMessageDto struct {
	MessageDto
	ReceiverID string `json:"receiver_id"`
}

type GroupMessageDto struct {
	MessageDto
	GroupID string `json:"group_id"`
}

type chatService struct {
	userRepo           repos.UserRepoInterface
	privateChatRepo    repos.PrivateChatRepoInterface
	groupRepo          repos.GroupRepoInterface
	groupMsgRepo       repos.GroupMessageRepoInterface
	privateMessageRepo repos.PrivateMessageRepoInterface
}

type ChatServiceInterface interface {
	SendPrivateMsg(userID uuid.UUID, data PrivateMessageDto) error
	SendMsgToGroup(userID uuid.UUID, data GroupMessageDto) error
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

func (c *chatService) SendPrivateMsg(userID uuid.UUID, data PrivateMessageDto) error {
	receiverUUID := uuid.MustParse(data.ReceiverID)
	chat, err := c.privateChatRepo.FindChat(receiverUUID, userID)
	if err == nil {
		pchat := models.PrivateMessage{
			BaseMessage: models.BaseMessage{
				Type:     data.Type,
				SenderID: userID,
				Content:  data.Content,
			},
			ChatID: chat.ID,
		}
		return c.privateMessageRepo.Create(nil, &pchat)
	}

	_, err = c.userRepo.FindByID(receiverUUID)
	if err != nil {
		return ErrUserNotFound
	}

	privateChat := &models.PrivateChat{
		FirstMemberID:  userID,
		SecondMemberID: receiverUUID,
	}

	// tx := c.privateChatRepo.BeginDBTx()
	if err := c.privateChatRepo.CreatePrivateChat(nil, privateChat); err != nil {
		// tx.Rollback()
		return ErrCreatingChat
	}

	privateMessage := &models.PrivateMessage{
		ChatID: privateChat.ID,
		BaseMessage: models.BaseMessage{
			Type:     data.Type,
			SenderID: userID,
			Content:  data.Content,
		},
	}
	if err := c.privateMessageRepo.Create(nil, privateMessage); err != nil {
		log.Println(err)
		// tx.Rollback()
		return ErrCreatingChat
	}

	// tx.Commit()

	return nil
}

func (c *chatService) SendMsgToGroup(userID uuid.UUID, data GroupMessageDto) error {
	groupUUID := uuid.MustParse(data.GroupID)
	_, err := c.groupRepo.FindByID(groupUUID)
	if err != nil {
		return err
	}

	msg := &models.GroupMessage{
		BaseMessage: models.BaseMessage{
			Type:     data.Type,
			SenderID: userID,
		},
		GroupID: groupUUID,
	}

	return c.groupMsgRepo.Create(msg)
}
