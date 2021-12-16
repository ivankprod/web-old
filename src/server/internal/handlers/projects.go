package handlers

import (
	"os"

	"github.com/gofiber/fiber/v2"

	"ivankprod.ru/src/server/internal/models"
)

func HandlerProjectsIndex(c *fiber.Ctx) error {
	uAuth, ok := c.Locals("user_auth").(*models.User)
	if !ok {
		uAuth = nil
	}

	data := make(fiber.Map)

	if uAuth != nil {
		data["user"] = *uAuth
	}

	err := c.Render("projects", fiber.Map{
		"urlBase":      c.BaseURL(),
		"urlCanonical": c.BaseURL() + c.Path(),
		"pageTitle":    "Проекты - " + os.Getenv("INFO_TITLE_BASE"),
		"pageDesc":     os.Getenv("INFO_DESC_BASE"),
		"pageScope":    "projects",
		"ogTags": fiber.Map{
			"title": "Проекты - " + os.Getenv("INFO_TITLE_BASE"),
			"type":  "website",
		},
		"activeProjects": true,
		"data":           data,
	})

	return err
}

func HandlerProjectsView(c *fiber.Ctx) error {
	uAuth, ok := c.Locals("user_auth").(*models.User)
	if !ok {
		uAuth = nil
	}

	data := make(fiber.Map)

	if uAuth != nil {
		data["user"] = *uAuth
	}

	title := "Проекты"

	switch c.Params("type") {
	case "it":
		title += ": IT-технологии"
	}

	path := "projects_" + c.Params("type")

	if err := c.Render(path, fiber.Map{
		"urlBase":      c.BaseURL(),
		"urlCanonical": c.BaseURL() + c.Path(),
		"pageTitle":    title + " - " + os.Getenv("INFO_TITLE_BASE"),
		"pageDesc":     os.Getenv("INFO_DESC_BASE"),
		"pageScope":    "projects",
		"ogTags": fiber.Map{
			"title": title + " - " + os.Getenv("INFO_TITLE_BASE"),
			"type":  "website",
		},
		"activeProjects": true,
		"data":           data,
	}); err == nil {
		return nil
	}

	return fiber.NewError(fiber.StatusNotFound, "Запрашиваемая страница не найдена либо ещё не создана")
}
