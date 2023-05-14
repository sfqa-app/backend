package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dbinstance struct {
	Db *gorm.DB
}

var Db Dbinstance

func ConnectDb() {
	dsn := "host=localhost user=postgres password='' dbname=go-db port=5432 sslmode=disable TimeZone=Asia/Beirut"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
	}

	log.Println("connected")
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("running migrations")

	db.AutoMigrate(&models.User{})

	Db = Dbinstance{
		Db: db,
	}
}
