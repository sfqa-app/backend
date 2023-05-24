package handler

import (
	"errors"
	"fmt"
	"net/smtp"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/sfqa-app/backend/database"
	"github.com/sfqa-app/backend/models"
)

// return verification email message body with token
func NewVerificationEmailMessage(to, token string) (msg string) {
	domain := os.Getenv("DOMAIN")
	from := os.Getenv("SMTP_EMAIL")

	return fmt.Sprintf(`From: SFQA App <%s>
To: %s
Subject: Verify Your Email Address
Thank you for signing up for SFQA App! We're excited to have you on board.

Please click on the link below to confirm your email address:

%s

If you did not sign up for SFQA App, please disregard this message and do not click the link above.

Thank you,

SFQA App Team
`,
		from,
		to,
		domain+"/verify/"+token,
	)
}

func EmailSend(msg, to string) error {
  // get email credentials from env
	from := os.Getenv("SMTP_EMAIL")
	pass := os.Getenv("SMTP_PASSWORD")
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")

  // authenticate to smtp server
	auth := smtp.PlainAuth("", from, pass, host)

  // send email
	err := smtp.SendMail(host+":"+port, auth, from, []string{to}, []byte(msg))
	if err != nil {
		return errors.New("error sending email")
	}

	return nil
}

func EmailVerify(c *fiber.Ctx) error {
  // get token from url params
	t := c.Params("token")

  // parse token
  claims, err := ParseJwtToken(c, t)
  if err != nil {
    return c.Status(fiber.StatusBadRequest).JSON(err.Error())
  }

  // parse userID from token
	userID := claims.Issuer

	var user models.User

  // get user from db
	if res := database.DB.First(&user, userID); res.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON("user not found")
	}

  // set user email verified to true
	user.IsEmailVerified = true

  // save user to db
	if res := database.DB.Save(&user); res.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON("error verifying email")
	}

  // get domain from env
	domain := os.Getenv("DOMAIN")

  // redirect to frontend
  return c.Redirect(domain, fiber.StatusTemporaryRedirect)
}
