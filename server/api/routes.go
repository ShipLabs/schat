package api

import (
	"shiplabs/schat/internal/base"
	"shiplabs/schat/internal/middlewares"
	"shiplabs/schat/internal/pkg/db"
	"shiplabs/schat/internal/pkg/store"

	"github.com/gin-gonic/gin"
)

func RoutesHandler(e *gin.Engine) {
	app := base.New(db.DB, store.WebsocketStore).MountHandlers()

	v1 := e.Group("api/v1")
	authRequired := v1.Group("").Use(middlewares.Auth)

	v1.POST("/register", app.AuthH.SignUp)
	v1.POST("/login", app.AuthH.Login)

	authRequired.GET("/connect", app.ChatH.EstablishConnection)
	authRequired.GET("/chat", app.ChatH.HandlePrivateChat)
	authRequired.GET("/group/create", app.ChatH.GroupCreationHandler)
	authRequired.GET("/group/message", app.ChatH.HandleGroupChat)
	authRequired.GET("/group/manage", app.ChatH.HandleMembership)
}
