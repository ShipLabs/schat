package handlers

import (
	"net/http"
	"shiplabs/schat/internal/services"
	"shiplabs/schat/pkg/shared"

	"github.com/gin-gonic/gin"
)

type LoginDto struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthHandlerInterface interface {
	SignUp(ctx *gin.Context)
	Login(ctx *gin.Context)
}

type authHandler struct {
	authService services.AuthServiceInterface
}

func NewAuthHandler(authS services.AuthServiceInterface) AuthHandlerInterface {
	return &authHandler{
		authService: authS,
	}
}

func (a *authHandler) SignUp(ctx *gin.Context) {
	var b LoginDto
	if !shared.ParseBody(ctx, &b) {
		return
	}

	token, err := a.authService.SignUp(b.Email, b.Password)
	if err != nil {
		shared.ErrorResponse(ctx, http.StatusUnprocessableEntity, err.Error())
	}

	shared.SuccessResponse(ctx, http.StatusOK, shared.SUCCESS, map[string]string{token: token})
}

func (a *authHandler) Login(ctx *gin.Context) {
	var b LoginDto
	if !shared.ParseBody(ctx, &b) {
		return
	}

	token, err := a.authService.Login(b.Email, b.Password)
	if err != nil {
		shared.ErrorResponse(ctx, http.StatusUnprocessableEntity, err.Error())
	}

	shared.SuccessResponse(ctx, http.StatusOK, shared.SUCCESS, map[string]string{token: token})
}
