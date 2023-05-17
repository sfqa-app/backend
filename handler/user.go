package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/sfqa-app/backend/database"
	"github.com/sfqa-app/backend/models"
)

func UserGet(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	if res := database.DB.First(&user, id); res.Error != nil {
		c.Status(fiber.StatusBadRequest)
		return res.Error
	}

	return c.Status(fiber.StatusOK).JSON(user)
}

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
