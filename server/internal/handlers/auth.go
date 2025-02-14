package handlers

import (
	"shiplabs/schat/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthHandlerInterface interface {
}

type authHandler struct {
	authService services.AuthServiceInterface
}

func NewAuthHandler() AuthHandlerInterface {
	return &authHandler{}
}

func (h *authHandler) SignUp(ctx *gin.Context) {

}

func (h *authHandler) Login(ctx *gin.Context) {

}
