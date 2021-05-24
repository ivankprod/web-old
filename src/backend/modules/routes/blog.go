package routes

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"ivankprod.ru/src/backend/modules/models"
)

func RouteBlogIndex(c *fiber.Ctx) error {
	uAuth, ok := c.Locals("user_auth").(*models.User)
	if !ok {
		uAuth = nil
	}

	data := make(fiber.Map)

	if uAuth != nil {
		data["user"] = *uAuth
	}

	err := c.Render("blog", fiber.Map{
		"urlBase":      c.BaseURL(),
		"urlCanonical": c.BaseURL() + c.Path(),
		"pageTitle":    "Блог - " + os.Getenv("INFO_TITLE_BASE"),
		"pageDesc":     os.Getenv("INFO_DESC_BASE"),
		"pageScope":    "blog",
		"ogTags": fiber.Map{
			"title": os.Getenv("INFO_TITLE_BASE"),
		},
		"activeBlog": true,
		"data":       data,
	})
	if err == nil {
		return nil
	}

	return fiber.NewError(fiber.StatusNotFound, "Страница не найдена!")
}