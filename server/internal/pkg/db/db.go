package db

import (
	"fmt"

	"shiplabs/schat/internal/models"
	"shiplabs/schat/internal/pkg/config"
	"shiplabs/schat/pkg/shared"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dbConfigs := config.Configs.DB

	dbconnStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s",
		dbConfigs.Host,
		dbConfigs.Port,
		dbConfigs.User,
		dbConfigs.Password,
		dbConfigs.Name,
	)

	db, err := gorm.Open(postgres.Open(dbconnStr), &gorm.Config{
		NowFunc: shared.TimeNow,
	})
	if err != nil {
		fmt.Println("Error connecting to database: ", err)
		panic(err)
	}

	err = db.AutoMigrate(
		&models.User{}, &models.PrivateChat{}, &models.GroupMessage{},
		&models.PrivateMessage{}, &models.Group{}, &models.GroupMember{},
	)

	if err != nil {
		fmt.Println("Error auto migrating database: ", err)
		panic(err)
	}
	fmt.Println("Database connected")

	DB = db
}
