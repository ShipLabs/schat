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
	BeginDBTx() *gorm.DB
	CreateGroup(group *models.Group) error
	GetUserGroups(userID uuid.UUID) ([]models.Group, error)
	FindByID(groupID uuid.UUID) (models.Group, error)
	GetGroupMember(groupID, userID uuid.UUID) (models.GroupMember, error)
	CreateGroupMembership(tx *gorm.DB, membership *[]models.GroupMember) error
	RevokeMembership(groupID, userID uuid.UUID) error
}

func NewGroupRepo(db gorm.DB) GroupRepoInterface {
	return &groupRepo{
		DB: db,
	}
}

func (g *groupRepo) BeginDBTx() *gorm.DB {
	return g.DB.Begin()
}

func (g *groupRepo) CreateGroup(group *models.Group) error {
	return g.DB.Create(group).Error
}

func (g *groupRepo) GetUserGroups(userID uuid.UUID) ([]models.Group, error) {
	var memberships []models.GroupMember
	var groups []models.Group
	err := g.DB.Joins("Groups").Where("user_id = ?", userID).Find(&memberships).Error
	if err != nil {
		return groups, err
	}
	for _, m := range memberships {
		groups = append(groups, m.Group)
	}

	return groups, nil
}

// func (g *groupRepo) GetUserGroups(userID uuid.UUID) ([]models.Group, error) {
//     var groups []models.Group
//     err := g.DB.Joins("JOIN group_members ON group_members.group_id = groups.id").
//         Where("group_members.user_id = ?", userID).
//         Find(&groups).Error
//     if err != nil {
//         return nil, err
//     }
//     return groups, nil
// }

func (g *groupRepo) FindByID(groupID uuid.UUID) (models.Group, error) {
	var group models.Group
	err := g.DB.Where("id=?", groupID).First(group).Error
	return group, err
}

func (g *groupRepo) GetGroupMember(groupID, userID uuid.UUID) (models.GroupMember, error) {
	var groupMember models.GroupMember
	err := g.DB.Where("user_id=? AND group_id=?", userID, groupID).First(&groupMember).Error
	return groupMember, err
}

func (g *groupRepo) CreateGroupMembership(tx *gorm.DB, membership *[]models.GroupMember) error {
	if tx != nil {
		return tx.Create(membership).Error
	}

	return g.DB.Create(membership).Error
}

func (g *groupRepo) RevokeMembership(groupID, userID uuid.UUID) error {
	return g.DB.Unscoped().Where("user_id=? AND group_id=?", userID, groupID).Delete(&models.GroupMember{}).Error
}
