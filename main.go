package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/sfqa-app/backend/database"
	"github.com/sfqa-app/backend/routes"
)

var port string

const defaultPort = "8080"

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	if port = os.Getenv("PORT"); port == "" {
		port = defaultPort
	}

	database.ConnectDb()
}

func main() {
	app := fiber.New()

	routes.SetUpRoutes(app)

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
