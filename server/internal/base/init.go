package base

import (
	"shiplabs/schat/internal/handlers"

	"gorm.io/gorm"
)

type base struct {
	db *gorm.DB
}

type baseHandlers struct {
	AuthH handlers.AuthHandlerInterface
}

func New(db *gorm.DB) *base {
	return &base{
		db: db,
	}
}

func (b *base) MountHandlers() baseHandlers {
	var h baseHandlers

	h.AuthH = b.WithAuthController()

	return h
}
