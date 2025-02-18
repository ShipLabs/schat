package services

import (
	"errors"
	"shiplabs/schat/internal/models"
	repos "shiplabs/schat/internal/repositories"

	"github.com/google/uuid"
)

type CreateGroupDto struct {
	CreatorID   uuid.UUID
	GroupName   string
	Description string
	Members     []uuid.UUID
}

type groupService struct {
	userRepo  repos.UserRepoInterface
	groupRepo repos.GroupRepoInterface
}

type GroupServiceInterface interface {
	CreateGroup(data CreateGroupDto) error
	AddToGroup(groupID, adminID, newMemberID uuid.UUID) error
	RemoveFromGroup(groupID, adminID, memberID uuid.UUID) error
}

func NewGroupService(
	userRepo repos.UserRepoInterface,
	groupRepo repos.GroupRepoInterface,
) GroupServiceInterface {
	return &groupService{
		userRepo:  userRepo,
		groupRepo: groupRepo,
	}
}

var (
	ErrNotAdmin = errors.New("user not group admin")
)

func (g *groupService) CreateGroup(data CreateGroupDto) error {
	tx := g.groupRepo.BeginDBTx()
	group := models.Group{
		CreatorID:   data.CreatorID,
		Name:        data.GroupName,
		Description: &data.Description,
	}
	//if group is not modified id will be (0000-0000-0000), watch out for that
	if err := g.groupRepo.CreateGroup(&group); err != nil {
		tx.Rollback()
		return err
	}

	members := g.buildMembershipSlice(group.ID, data.CreatorID, data.Members)
	if err := g.groupRepo.CreateGroupMembership(tx, &members); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (g *groupService) buildMembershipSlice(groupID, adminID uuid.UUID, membersID []uuid.UUID) []models.GroupMember {
	var members []models.GroupMember

	admin := models.GroupMember{
		UserID:  adminID,
		GroupID: groupID,
		Role:    models.Admin,
	}
	members = append(members, admin)

	for _, memberID := range membersID {
		member := models.GroupMember{
			UserID:  memberID,
			GroupID: groupID,
			Role:    models.Member,
		}
		members = append(members, member)
	}

	return members
}

func (g *groupService) AddToGroup(groupID, adminID, newMemberID uuid.UUID) error {
	_, err := g.userRepo.FindByID(newMemberID)
	if err != nil {
		return err
	}
	_, err = g.groupRepo.FindByID(groupID)
	if err != nil {
		return err
	}

	if !g.isGroupAdmin(groupID, adminID) {
		return ErrNotAdmin
	}

	memberShip := models.GroupMember{
		UserID:  newMemberID,
		GroupID: groupID,
		Role:    models.Member,
	}

	if err := g.groupRepo.CreateGroupMembership(nil, &[]models.GroupMember{memberShip}); err != nil {
		return err
	}

	return nil
}

func (g *groupService) RemoveFromGroup(groupID, adminID, memberID uuid.UUID) error {
	if g.isGroupAdmin(groupID, adminID) {
		return g.groupRepo.RevokeMembership(groupID, memberID)
	}

	return ErrNotAdmin
}

func (g *groupService) isGroupAdmin(groupID, userID uuid.UUID) bool {
	membership, err := g.groupRepo.GetGroupMember(groupID, userID)
	if err != nil {
		return false
	}
	return membership.Role == models.Admin
}
