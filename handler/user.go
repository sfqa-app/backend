package handler

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sfqa-app/backend/database"
	"github.com/sfqa-app/backend/models"
)

// UserGet get a user account
//
//	@Summary       Get user
//	@Description	Get user account
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			id 	path		int	true	"User ID"
//	@Success		200	{object}	models.User
//	@Failure		400	{object}	error
//	@Router			/user/{id} [get]
func UserGet(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	if res := database.DB.First(&user, id); res.Error != nil {
		c.Status(fiber.StatusBadRequest)
		return res.Error
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// UserCreate creates a new user account
//
//	@Summary		Create user
//	@Description	Create new user account
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Param			User 	body		models.UserInfo	true	"User info"
//	@Success		200	{object}	models.User
//	@Failure		400	{object}	error
//	@Router			/user [post]
func UserCreate(c *fiber.Ctx) error {
	var data map[string]string

	// parse request body
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("failed to parse user data")
	}

	// create new user
	user := models.NewUser(data["email"], data["password"])

	// validate email
	if validEmail := user.IsValidEmail(); !validEmail {
		return c.Status(fiber.StatusBadRequest).JSON("email not valid")
	}

	// check if email already exists and verified
	if res := database.DB.Where("email = ?", user.Email).First(&user); res.Error == nil {
		if user.IsEmailVerified {
			return c.Status(fiber.StatusBadRequest).JSON("email already exists and verified")
		}
	} else {
		// encrypt password
		if err := user.EncryptPassword(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON("failed to encrypt password")
		}

		// create user
		if res := database.DB.Create(user); res.Error != nil {
			log.Println(res.Error)
			return c.Status(fiber.StatusBadRequest).JSON("failed to create user")
		}
	}

	// send email verification link
	if err := SendEmailVerificationCode(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("failed to send email verification link")
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func SendEmailVerificationCode(user *models.User) error {
	// get jwt secret from env
	secret := os.Getenv("JWT_SECRET")

	// generate jwt token with user id
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	// sign jwt token
	token, err := claims.SignedString([]byte(secret))
	if err != nil {
		return err
	}

	emailTemplate := `Subject: Verify Your Email Address
Thank you for signing up for SFQA App! We're excited to have you on board.

Please click on the link below to confirm your email address:

%s

If you did not sign up for SFQA App, please disregard this message and do not click the link above.

Thank you,

SFQA App Team
`

	domain := os.Getenv("DOMAIN")
	emailBody := fmt.Sprintf(emailTemplate, domain+"/verify-email?token="+token)

	return SendEmail(emailBody, user.Email)
}

// UserDelete deletes a user account
//
//	@Summary		Delete user
//	@Description	Delete user account
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.User
//	@Failure		400	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/user [delete]
func UserDelete(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("failed to parse user data")
	}

	if res := database.DB.Delete(&user, user.ID); res.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON("failed to delete user")
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// UserUpdate updates a user account info
//
//	@Summary		Update user account info
//	@Description	Update user account info
//	@Tags			User
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	models.User
//	@Failure		400	{object}	error
//	@Param			id	path		int	true	"User ID"
//	@Param			User	body		models.UserInfo	true	"User"
//	@Security		ApiKeyAuth
//	@Router			/user/{id} [put]
func UserUpdate(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("failed to parse user data")
	}

	if res := database.DB.Model(&user).Updates(&user); res.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON("failed to update user")
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func isValid(field string) bool {
	return field != ""
}

// Login get user and password
func UserLogin(c *fiber.Ctx) error {
	type LoginInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("failed to parse login data")
	}

	if !isValid(input.Password) || !isValid(input.Email) {
		return c.Status(fiber.StatusBadRequest).JSON("invalid email or password")
	}

	var user models.User
	if res := database.DB.Where("email = ?", input.Email).First(&user); res.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON("email not found")
	}

	if !user.IsEmailVerified {
		return c.Status(fiber.StatusBadRequest).JSON("email not verified")
	}

	if !user.IsPasswordMatch(input.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "wrong password",
		})
	}

	expireDate := time.Now().Add(time.Hour * 7 * 24) // 7 days

	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: expireDate.Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	token, err := claims.SignedString([]byte(secret))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  expireDate,
		HTTPOnly: true,
	}

	c.Cookie(&cookie)

	return c.SendStatus(fiber.StatusOK)
}

func UserLogout(c *fiber.Ctx) error {
	c.ClearCookie()
	return c.SendStatus(fiber.StatusOK)
}
