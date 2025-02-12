package base

import "gorm.io/gorm"

type base struct {
	db *gorm.DB
}

type baseHandlers struct {
}

func New(db *gorm.DB) *base {
	return &base{
		db: db,
	}
}

func (b *base) MountHandlers() baseHandlers {
	var h baseHandlers

	return h
}
