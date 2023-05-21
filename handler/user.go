package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sfqa-app/backend/database"
	"github.com/sfqa-app/backend/models"
)

// UserGet get a user account
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
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	if res := database.DB.Create(&user); res.Error != nil {
		c.Status(fiber.StatusBadRequest)
		return res.Error
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// UserDelete deletes a user account
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
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	if res := database.DB.Delete(&user, user.ID); res.Error != nil {
		c.Status(fiber.StatusBadRequest)
		return res.Error
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

// UserUpdate updates a user account info
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
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	if res := database.DB.Model(&user).Updates(&user); res.Error != nil {
		c.Status(fiber.StatusBadRequest)
		return res.Error
	}

	return c.Status(fiber.StatusOK).JSON(user)
}
