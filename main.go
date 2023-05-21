package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/sfqa-app/backend/database"
	"github.com/sfqa-app/backend/router"
)

//	@title			sfqa-app docs
//	@version		1.0
//	@description	sfqa-app api documentation
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	GPL-3.0
//	@license.url	https://www.gnu.org/licenses/gpl-3.0.en.html
//	@host			localhost:8080
//	@BasePath		/

var port string

const defaultPort = "8080"

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	if port = os.Getenv("PORT"); port == "" {
		port = defaultPort
	}

	database.ConnectDb()
}

func main() {
	app := fiber.New()

	app.Use(
		cors.New(cors.Config{
			AllowCredentials: true}),
		// encryptcookie.New(encryptcookie.Config{
		// 	Key: os.Getenv("COOKIE_SECRET"),
		)

	router.SetUpRoutes(app)

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
