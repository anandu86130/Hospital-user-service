package db

import (
	"log"
	"os"

	"github.com/anandu86130/Hospital-user-service/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() *gorm.DB {
	dsn, ok := os.LookupEnv("DSN")
	if !ok {
		log.Fatal("Failed to load database")
	}

	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	DB.AutoMigrate(&model.UserModel{})
	return DB
}
