package repos

import (
	"shiplabs/schat/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepoInterface interface {
	Create(user *models.User) error
	FindByEmail(email string) (models.User, error)
	FindByID(id uuid.UUID) (models.User, error)
}

type UserRepo struct {
	DB gorm.DB
}

func NewUserRepo(db gorm.DB) UserRepoInterface {
	return &UserRepo{
		DB: db,
	}
}

func (u *UserRepo) Create(user *models.User) error {
	return u.DB.Create(user).Error
}

func (u *UserRepo) FindByEmail(email string) (models.User, error) {
	var user models.User
	err := u.DB.Where("email=?", email).First(&user).Error
	return user, err
}

func (u *UserRepo) FindByID(id uuid.UUID) (models.User, error) {
	var user models.User
	err := u.DB.Where("id=?", id).First(&user).Error
	return user, err
}
