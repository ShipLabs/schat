package repos

import (
	"shiplabs/schat/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepoInterface interface {
	Create(user *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id uuid.UUID) (*models.User, error)
}

type UserRepo struct {
	db gorm.DB
}

func NewUserRepo(db gorm.DB) UserRepoInterface {
	return &UserRepo{
		db: db,
	}
}

func (u *UserRepo) Create(user *models.User) error {
	return u.db.Create(user).Error
}

func (u *UserRepo) FindByEmail(email string) (*models.User, error) {
	user := &models.User{}
	err := u.db.Where("email = ?", email).First(user).Error
	return user, err
}

func (u *UserRepo) FindByID(id uuid.UUID) (*models.User, error) {
	user := &models.User{}
	err := u.db.Where("id = ?", id).First(user).Error
	return user, err
}
