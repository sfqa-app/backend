package handler

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sfqa-app/backend/models"
)

func isValid(field string) bool {
	return field != ""
}

// Login get user and password
func Login(c *fiber.Ctx) error {
	var input models.UserInfo
	if err := c.BodyParser(&input); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if !isValid(input.Username) ||
		!isValid(input.Password) ||
		!isValid(input.Email) ||
		!isValid(input.Name) {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uname"] = input.Username
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	secret := os.Getenv("JWT_SECRET")
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(t)
}
