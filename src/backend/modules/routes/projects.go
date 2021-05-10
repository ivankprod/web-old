package routes

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

func RouteProjectsIndex(c *fiber.Ctx) error {
	uAuth := c.Locals("user_auth")
	data := make(fiber.Map)

	if uAuth != nil {
		data = fiber.Map{"user": uAuth}
	}

	err := c.Render("projects", fiber.Map{
		"urlCanonical": c.BaseURL() + c.Path(),
		"pageTitle":    "Проекты - " + os.Getenv("INFO_TITLE_BASE"),
		"pageDesc":     os.Getenv("INFO_DESC_BASE"),
		"ogTags": fiber.Map{
			"title": "Проекты - " + os.Getenv("INFO_TITLE_BASE"),
		},
		"activeProjects": true,
		"data":           data,
	})
	if err == nil {
		return nil
	}

	return fiber.NewError(fiber.StatusNotFound, "Страница не найдена!")
}

func RouteProjectsView(c *fiber.Ctx) error {
	uAuth := c.Locals("user_auth")
	data := make(fiber.Map)

	if uAuth != nil {
		data = fiber.Map{"user": uAuth}
	}

	var path = "projects_" + c.Params("type")

	err := c.Render(path, fiber.Map{
		"urlCanonical": c.BaseURL() + c.Path(),
		"pageTitle":    "Проекты - " + os.Getenv("INFO_TITLE_BASE"),
		"pageDesc":     os.Getenv("INFO_DESC_BASE"),
		"ogTags": fiber.Map{
			"title": "Проекты - " + os.Getenv("INFO_TITLE_BASE"),
		},
		"activeProjects": true,
		"data":           data,
	})
	if err == nil {
		return nil
	}

	return fiber.NewError(fiber.StatusNotFound, "Страница не найдена!")
}
