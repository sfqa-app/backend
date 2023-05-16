package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sfqa-app/backend/handlers"
)

func SetUpRoutes(app *fiber.App) {
	app.Get("/user/:id", handlers.UserGet)
	app.Put("/user", handlers.UserUpdate)
	app.Post("/user", handlers.UserCreate)
	app.Delete("/user", handlers.UserDelete)
}
