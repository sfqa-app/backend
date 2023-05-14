package routes

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func SetUpRoutes(app *fiber.App) {
  app.Get("/hello", helloWorld)
  app.Get("/:user", helloUser)
}

func helloWorld(c *fiber.Ctx) error {
	return c.SendString("Hello, World!\n")
}

func helloUser(c *fiber.Ctx) error {
	return c.SendString(fmt.Sprintf("Hello, %v!\n", c.Params("user")))
}
