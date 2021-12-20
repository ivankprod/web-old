package handlers

import (
	"os"

	"github.com/gofiber/fiber/v2"

	"ivankprod.ru/src/server/internal/domain"
)

func HandlerHomeIndex(c *fiber.Ctx) error {
	uAuth, ok := c.Locals("user_auth").(*domain.User)
	if !ok {
		uAuth = nil
	}

	data := make(fiber.Map)

	if uAuth != nil {
		data["user"] = *uAuth
	}

	err := c.Render("index", fiber.Map{
		"urlBase":      c.BaseURL(),
		"urlCanonical": c.BaseURL() + c.Path(),
		"pageTitle":    "Главная - " + os.Getenv("INFO_TITLE_BASE"),
		"pageDesc":     os.Getenv("INFO_DESC_BASE"),
		"pageScope":    "home",
		"ogTags": fiber.Map{
			"title": os.Getenv("INFO_TITLE_BASE"),
			"type":  "website",
		},
		"activeHome": true,
		"data":       data,
	})

	return err
}
