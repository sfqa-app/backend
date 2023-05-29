package handler

import (
	"errors"
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

	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("failed to parse user data")
	}

	user := models.NewUser(data["email"], data["password"])

	if validEmail := user.IsValidEmail(); !validEmail {
		return c.Status(fiber.StatusBadRequest).JSON("email not valid")
	}

	res := database.DB.Where("email = ?", user.Email).First(&user, user.Email)
	if res.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON("failed to find user")
	}

	if user.EmailVerified {
		return c.Status(fiber.StatusBadRequest).JSON("email already exists and verified")
	}

	if err := user.EncryptPassword(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("failed to encrypt password")
	}

	if res := database.DB.Create(user); res.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON("failed to create user")
	}

	if err := sendEmailVerificationLink(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("failed to send email verification link")
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

func sendEmailVerificationLink(user *models.User) error {
	claims := &jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}

	token, err := GenerateToken(claims)
	if err != nil {
		return err
	}

	msg := NewVerificationEmailMessage(user.Email, token)

	return EmailSend(msg, user.Email)
}

func UserPasswordReset(c *fiber.Ctx) error {
	var user models.User

	claims, err := getJWTCookieClaims(c)
	if err != nil {
		return err
	}

	userID := claims.Issuer

	if res := database.DB.First(&user, userID); res.Error != nil ||
		user.LoginMethod != "email" {
		return c.Status(fiber.StatusBadRequest).JSON("email not found")
	}

	if !user.EmailVerified {
		return c.Status(fiber.StatusBadRequest).JSON("email not verified")
	}

	if err := sendResetPasswordLink(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("failed to send reset password link")
	}

	return c.Status(fiber.StatusOK).JSON("password reset link sent")
}

func sendResetPasswordLink(user *models.User) error {
	claims := &jwt.StandardClaims{
		Issuer:    strconv.Itoa(int(user.ID)),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	}

	token, err := GenerateToken(claims)
	if err != nil {
		return err
	}

	msg := NewResetPasswordEmailMessage(user.Email, token)

	return EmailSend(msg, user.Email)
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

	claims, err := getJWTCookieClaims(c)
	if err != nil {
		return err
	}

	userID := claims.Issuer

	if res := database.DB.Delete(&user, userID); res.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON("failed to delete user")
	}

	return c.Status(fiber.StatusOK).JSON("user deleted")
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

	if err := userMe(c, &user); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(err.Error())
	}

	if res := database.DB.Model(&user).Updates(&user); res.Error != nil {
		return c.Status(fiber.StatusBadRequest).JSON("failed to update user")
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// field is valid if it is not empty
func isValidField(field string) bool {
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

	if !isValidField(input.Password) || !isValidField(input.Email) {
		return c.Status(fiber.StatusBadRequest).JSON("invalid email or password")
	}

	var user *models.User
	if res := database.DB.Where("email = ?", input.Email).First(user); res.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON("email not found")
	}

	if !user.EmailVerified {
		return c.Status(fiber.StatusBadRequest).JSON("email not verified")
	}

	if !user.IsPasswordMatch(input.Password) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "wrong password",
		})
	}

  SetUserCookie(c, user)

	return c.SendStatus(fiber.StatusOK)
}

func SetUserCookie(c *fiber.Ctx, user *models.User) error {
  expireDate := time.Now().Add(time.Hour * 7 * 24)

  claims := &jwt.StandardClaims{
    Issuer:    strconv.Itoa(int(user.ID)),
    ExpiresAt: expireDate.Unix(),
  }

  token, err := GenerateToken(claims)
  if err != nil {
    return err
  }

  cookie := fiber.Cookie{
    Name:     "jwt",
    Value:    token,
    Expires:  expireDate,
    HTTPOnly: true,
  }

  c.Cookie(&cookie)

  return nil
}

// UserLogout logout user
func UserLogout(c *fiber.Ctx) error {
	c.ClearCookie()
	return c.SendStatus(fiber.StatusOK)
}

// user is not allowed to mess with other user's account
func userMe(c *fiber.Ctx, user *models.User) error {
	claims, err := getJWTCookieClaims(c)
	if err != nil {
		return err
	}

	userID, err := strconv.Atoi(claims.Issuer)
	if err != nil {
		return errors.New("invalid user id")
	}

	if user.ID != uint(userID) {
		return errors.New("unauthorized")
	}

	return nil
}
