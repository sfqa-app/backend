package main

import (
	"github.com/sfqa-app/backend/database"
	"github.com/sfqa-app/backend/config"
	"github.com/sfqa-app/backend/models"
)

func init() {
  config.LoadEnv()
	database.ConnectDb()
}

func main() {
	database.DB.AutoMigrate(&models.User{})
}
