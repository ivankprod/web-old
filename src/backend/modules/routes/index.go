package routes

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

func RouteHomeIndex(c *fiber.Ctx) error {
	uAuth := c.Locals("user_auth")
	data := make(fiber.Map)

	if uAuth != nil {
		data = fiber.Map{"user": uAuth}
	}

	err := c.Render("index", fiber.Map{
		"urlCanonical": c.BaseURL() + c.Path(),
		"pageTitle":    "Главная - " + os.Getenv("INFO_TITLE_BASE"),
		"pageDesc":     os.Getenv("INFO_DESC_BASE"),
		"ogTags": fiber.Map{
			"title": os.Getenv("INFO_TITLE_BASE"),
		},
		"activeHome": true,
		"data":       data,
	})
	if err == nil {
		return nil
	}

	return fiber.NewError(fiber.StatusNotFound, "Страница не найдена!")
}
