package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/sfqa-app/backend/database"
	"github.com/sfqa-app/backend/models"
)

func UserGet(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	if res := database.DB.Db.First(&user, id); res.Error != nil {
		c.Status(http.StatusBadRequest)
		return res.Error
	}

	return c.Status(http.StatusOK).JSON(user)
}

func UserCreate(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	if res := database.DB.Db.Create(&user); res.Error != nil {
		c.Status(http.StatusBadRequest)
		return res.Error
	}

	return c.Status(http.StatusOK).JSON(user)
}

func UserDelete(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	if res := database.DB.Db.Delete(&user, user.ID); res.Error != nil {
		c.Status(http.StatusBadRequest)
		return res.Error
	}

	return c.Status(http.StatusOK).JSON(user)
}

func UserUpdate(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(err.Error())
	}

	if res := database.DB.Db.Model(&user).Updates(&user); res.Error != nil {
		c.Status(http.StatusBadRequest)
		return res.Error
	}

	return c.Status(http.StatusOK).JSON(user)
}
