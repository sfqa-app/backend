package handler

import (
	"errors"
	"net/smtp"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sfqa-app/backend/database"
	"github.com/sfqa-app/backend/models"
)

func EmailSend(msg string, to string) error {
	from := os.Getenv("SMTP_EMAIL")
	pass := os.Getenv("SMTP_PASSWORD")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	auth := smtp.PlainAuth("", from, pass, host)

	err := smtp.SendMail(host+":"+port, auth, from, []string{to}, []byte(msg))
	if err != nil {
		return errors.New("error sending email")
	}
	return nil
}

func EmailVerify(c *fiber.Ctx) error {
	t := c.Params("token")
	token, err := jwt.ParseWithClaims(t, &jwt.StandardClaims{}, func(t *jwt.Token) (interface{}, error) {
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("error parsing token")
	}

	claims, ok := token.Claims.(*jwt.StandardClaims)
	if !ok || !token.Valid || claims.ExpiresAt < time.Now().Unix() {
		return c.Status(fiber.StatusBadRequest).JSON("not valid token")
	}

	userID := claims.Issuer

	var user models.User

	if res := database.DB.First(&user, userID); res.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON("user not found")
	}

	user.IsEmailVerified = true

	if res := database.DB.Save(&user); res.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON("error verifying email")
	}

	return c.Status(fiber.StatusOK).JSON("email verified")
}
