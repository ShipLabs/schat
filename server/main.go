package main

import (
	"shiplabs/schat/api"
	"shiplabs/schat/internal/pkg/config"
	"shiplabs/schat/internal/pkg/db"

	"github.com/gin-gonic/gin"
)

func init() {
	config.Load()
	db.Connect()
}

func main() {
	s := gin.New()
	s.Use(gin.Recovery())
	api.RoutesHandler(s)

	if err := s.Run(":" + config.Configs.Port); err != nil {
		panic(err)
	}
}
