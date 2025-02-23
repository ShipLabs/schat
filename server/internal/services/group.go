package services

import (
	"errors"
	"log"
	"shiplabs/schat/internal/models"
	repos "shiplabs/schat/internal/repositories"

	"github.com/google/uuid"
)

type CreateGroupDto struct {
	GroupName   string   `json:"group_name"`
	Description string   `json:"description"`
	Members     []string `json:"members"`
}

type GroupMembershipAction string

const (
	Add    GroupMembershipAction = "add"
	Remove GroupMembershipAction = "remove"
)

type GroupMembershipDto struct {
	MemberID string                `json:"member_id"`
	Action   GroupMembershipAction `json:"action"`
}

type groupService struct {
	userRepo  repos.UserRepoInterface
	groupRepo repos.GroupRepoInterface
}

type GroupServiceInterface interface {
	CreateGroup(userID uuid.UUID, data CreateGroupDto) error
	GetGroupMembers(groupID uuid.UUID) ([]models.GroupMember, error)
	HandleMembership(groupID, adminID, memberID uuid.UUID, action GroupMembershipAction) error
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

func (g *groupService) CreateGroup(userID uuid.UUID, data CreateGroupDto) error {
	tx := g.groupRepo.BeginDBTx()
	group := models.Group{
		CreatorID:   userID,
		Name:        data.GroupName,
		Description: &data.Description,
	}
	//if group is not modified id will be (0000-0000-0000), watch out for that
	if err := g.groupRepo.CreateGroup(&group); err != nil {
		tx.Rollback()
		return err
	}

	members := g.buildMembershipSlice(group.ID, userID, data.Members)
	if err := g.groupRepo.CreateGroupMembership(tx, &members); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (g *groupService) buildMembershipSlice(groupID, adminID uuid.UUID, membersID []string) []models.GroupMember {
	var members []models.GroupMember

	admin := models.GroupMember{
		UserID:  adminID,
		GroupID: groupID,
		Role:    models.Admin,
	}
	members = append(members, admin)

	//should probably check if the users being added exist. but will research efficient ways to do that
	for _, memberID := range membersID {
		mUUID, err := uuid.Parse(memberID)
		if err != nil {
			log.Println(err)
			continue
		}
		member := models.GroupMember{
			UserID:  mUUID,
			GroupID: groupID,
			Role:    models.Member,
		}
		members = append(members, member)
	}

	return members
}

func (g *groupService) addToGroup(groupID, adminID, newMemberID uuid.UUID) error {
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

func (g *groupService) GetGroupMembers(groupID uuid.UUID) ([]models.GroupMember, error) {
	return g.groupRepo.GetGroupMembers(groupID)
}

func (g *groupService) HandleMembership(groupID, adminID, memberID uuid.UUID, action GroupMembershipAction) error {
	switch action {
	case Add:
		return g.addToGroup(groupID, adminID, memberID)
	case Remove:
		return g.removeFromGroup(groupID, adminID, memberID)
	default:
		return errors.New("invalid action")
	}
}

func (g *groupService) removeFromGroup(groupID, adminID, memberID uuid.UUID) error {
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
