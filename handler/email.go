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

func NewResetPasswordEmailMessage(to, token string) (msg string) {
	domain := os.Getenv("DOMAIN")
	from := os.Getenv("SMTP_EMAIL")

	return fmt.Sprintf(`From: SFQA App <%s>
To: %s
Subject: Reset Your Password
Use the link below to reset your password:

%s

If you did not request to reset your password, please disregard this message and do not click the link above.

Thank you,

SFQA App Team
`,
		from,
		to,
		domain+"/reset-password/"+token)
}

func EmailSend(msg, to string) error {
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

	claims, err := ParseJwtToken(c, t)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	userID := claims.Issuer

	var user models.User

	if res := database.DB.First(&user, userID); res.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON("user not found")
	}

  if user.EmailVerified {
    return c.Status(fiber.StatusBadRequest).JSON("email already verified")
  }

	user.EmailVerified = true

	if res := database.DB.Save(&user); res.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON("error verifying email")
	}

	domain := os.Getenv("DOMAIN")

	return c.Redirect(domain, fiber.StatusTemporaryRedirect, fiber.StatusOK)
}
