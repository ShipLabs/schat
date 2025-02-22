package middlewares

import (
	"errors"
	"net/http"
	"shiplabs/schat/internal/pkg/config"
	"shiplabs/schat/internal/pkg/db"
	repos "shiplabs/schat/internal/repositories"
	"shiplabs/schat/pkg/shared"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	ErrCredentialsRequired = "credentials required"
)

// var userRepo = repos.NewUserRepo(*db.DB) find how to instantiate the user repo and have it available

func Auth(ctx *gin.Context) {
	authT := ctx.GetHeader("Authorization")
	if authT == "" {
		shared.ErrorResponse(ctx, http.StatusUnauthorized, ErrCredentialsRequired)
		ctx.Abort()
		return
	}

	jwtToken, err := extractBearerToken(authT)
	if err != nil {
		shared.ErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		ctx.Abort()
		return
	}

	claims := jwt.RegisteredClaims{}
	jwt.ParseWithClaims(jwtToken, &claims, func(token *jwt.Token) (any, error) {
		return []byte(config.Configs.APP_SECRET), nil
	})

	userId := uuid.MustParse(claims.Subject)
	user, err := repos.NewUserRepo(*db.DB).FindByID(userId)
	if err != nil {
		shared.ErrorResponse(ctx, http.StatusUnauthorized, err.Error())
		ctx.Abort()
		return
	}

	ctx.Set("userID", user.ID.String())
	ctx.Next()
}

func extractBearerToken(header string) (string, error) {
	if header == "" {
		return "", errors.New(ErrCredentialsRequired)
	}

	jwtToken := strings.Split(header, " ")
	if len(jwtToken) != 2 {
		return "", errors.New(ErrCredentialsRequired)
	}

	return jwtToken[1], nil
}
