package services

import (
	"errors"
	"shiplabs/schat/internal/models"
	"shiplabs/schat/internal/pkg/config"
	repos "shiplabs/schat/internal/repositories"
	"shiplabs/schat/pkg/shared"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthServiceInterface interface {
	Login(email string, password string) (string, error)
	SignUp(email string, password string) (string, error)
}

type AuthService struct {
	UserRepo repos.UserRepoInterface
}

const (
	ErrInvalidCredentials = "invalid credentials"
)

func NewAuthService(userRepo repos.UserRepoInterface) AuthServiceInterface {
	return &AuthService{
		UserRepo: userRepo,
	}
}

func (a *AuthService) Login(email string, password string) (string, error) {
	user, err := a.UserRepo.FindByEmail(email)
	if err != nil {
		return "", err
	}

	if !shared.VerifyDataHash(password, user.Password) {
		return "", errors.New(ErrInvalidCredentials)
	}

	return a.signJWT(user.ID)
}

func (a *AuthService) SignUp(email string, password string) (string, error) {
	user := &models.User{
		Email:    email,
		Password: shared.HashData(password),
	}

	if err := a.UserRepo.Create(user); err != nil {
		return "", err
	}

	return a.signJWT(user.ID)
}

func (a *AuthService) signJWT(userId uuid.UUID) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   userId.String(),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(700000)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.Configs.APP_SECRET))
}
