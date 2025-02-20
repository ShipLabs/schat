package api

import (
	"shiplabs/schat/internal/base"
	"shiplabs/schat/internal/pkg/db"

	"github.com/gin-gonic/gin"
)

func RoutesHandler(e *gin.Engine) {
	app := base.New(db.DB).MountHandlers()

	v1 := e.Group("api/v1")
	// authRequired := v1.Group("").Use(middlewares.Auth)

	v1.POST("/register", app.AuthH.SignUp)
	v1.POST("/login", app.AuthH.Login)

}
