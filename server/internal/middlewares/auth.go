package middlewares

import (
	"errors"
	"shiplabs/schat/internal/pkg/config"
	"shiplabs/schat/internal/pkg/db"
	repos "shiplabs/schat/internal/repositories"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	ErrCredentialsRequired = "credentials required"
)

var userRepo = repos.NewUserRepo(*db.DB)

func Auth(ctx *gin.Context) {
	authT := ctx.GetHeader("Authorization")
	if authT == "" {
		ctx.JSON(401, gin.H{"error": "Authorization header required"})
		ctx.Abort()
		return
	}

	jwtToken, err := extractBearerToken(authT)
	if err != nil {
		ctx.JSON(401, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	claims := jwt.RegisteredClaims{}
	jwt.ParseWithClaims(jwtToken, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Configs.APP_SECRET), nil
	})

	userId := uuid.MustParse(claims.Subject)
	user, err := userRepo.FindByID(userId)
	if err != nil {
		ctx.JSON(401, gin.H{"error": "User not found"})
		ctx.Abort()
		return
	}

	ctx.Set("userID", user.ID)
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
