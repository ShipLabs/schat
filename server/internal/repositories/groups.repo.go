package repos

import (
	"shiplabs/schat/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type groupRepo struct {
	DB gorm.DB
}

type GroupRepoInterface interface {
	CreateGroup(group *models.Group) error
	GetUserGroups(userID uuid.UUID) ([]models.Group, error)
}

func NewGroupRepo(db gorm.DB) GroupRepoInterface {
	return &groupRepo{
		DB: db,
	}
}

func (g *groupRepo) CreateGroup(group *models.Group) error {
	return g.DB.Create(group).Error
}

func (g *groupRepo) GetUserGroups(userID uuid.UUID) ([]models.Group, error) {
	var groups []models.Group
	//TODO
	return groups, nil
}
