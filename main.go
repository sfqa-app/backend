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
	// load env variables from an .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// load env variables
	port = os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

  database.ConnectDb()
}

func main() {
	app := fiber.New()

	routes.SetUpRoutes(app)

	err := app.Listen(":" + port)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
