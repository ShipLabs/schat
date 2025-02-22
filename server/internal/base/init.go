package base

import (
	"shiplabs/schat/internal/handlers"
	"shiplabs/schat/internal/pkg/store"

	"gorm.io/gorm"
)

type base struct {
	db      *gorm.DB
	wsStore store.ConnectionStoreInterface
}

type baseHandlers struct {
	AuthH handlers.AuthHandlerInterface
	ChatH handlers.WsHandlerInterface
}

func New(db *gorm.DB, store store.ConnectionStoreInterface) *base {
	return &base{
		db:      db,
		wsStore: store,
	}
}

func (b *base) MountHandlers() baseHandlers {
	var h baseHandlers

	h.AuthH = b.WithAuthController()
	h.ChatH = b.WithChatController()

	return h
}
