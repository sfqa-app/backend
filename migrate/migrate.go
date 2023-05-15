package main

import (
	"github.com/sfqa-app/backend/database"
	"github.com/sfqa-app/backend/initializers"
	"github.com/sfqa-app/backend/models"
)

func init() {
  initializers.LoadEnv()
	database.ConnectDb()
}

func main() {
	database.Db.Db.AutoMigrate(&models.User{})
}
