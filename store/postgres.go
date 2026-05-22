package storage

import (
	"calculationengine/constants"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func Connect() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		constants.AppConfig.DBHost,
		constants.AppConfig.DBUser,
		constants.AppConfig.DBPassword,
		constants.AppConfig.DBName,
		constants.AppConfig.DBPort,
		constants.AppConfig.SSLMode,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

}
