package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sfqa-app/backend/handler"
	"github.com/sfqa-app/backend/middleware"
)

func SetUpRoutes(app *fiber.App) {
  // auth
  app.Post("/login", handler.Login)

  // user routes
  user := app.Group("/user")
	user.Post("/", handler.UserCreate)
	user.Get("/:id<int>", handler.UserGet)
	user.Put("/", middleware.Protected(), handler.UserUpdate)
  user.Delete("/", middleware.Protected(), handler.UserDelete)
}
