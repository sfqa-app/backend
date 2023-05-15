package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dbinstance struct {
	Db *gorm.DB
}

var Db Dbinstance

func ConnectDb() {
	log.Println("connecting to database...")

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, pass, name, port,
	)

	var db *gorm.DB
	var err error
	var reconnectSecondsInterval time.Duration = 5

	for {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})

		if err != nil {
			log.Printf("Failed to connect to database, retrying in %v seconds\n", reconnectSecondsInterval)
			time.Sleep(reconnectSecondsInterval * time.Second)
		} else {
			log.Println("connected to database")
			break
		}
	}

	Db = Dbinstance{
		Db: db,
	}
}
