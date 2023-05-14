package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/sfqa-app/backend/routes"
)

var port string
const defaultPort = "8080"

func init() {
  port = os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
}

func main() {
	app := fiber.New()

  routes.SetUpRoutes(app)

  err := app.Listen(":" + port)
  if err != nil {
    log.Fatalf("Error starting server: %v", err)
  }
}
