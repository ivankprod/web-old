package routes

import (
	"os"

	"github.com/gofiber/fiber/v2"

	"ivankprod.ru/src/server/internal/models"
)

func RouteAboutIndex(c *fiber.Ctx) error {
	uAuth, ok := c.Locals("user_auth").(*models.User)
	if !ok {
		uAuth = nil
	}

	data := make(fiber.Map)

	if uAuth != nil {
		data["user"] = *uAuth
	}

	err := c.Render("about", fiber.Map{
		"urlBase":      c.BaseURL(),
		"urlCanonical": c.BaseURL() + c.Path(),
		"pageTitle":    "О компании - " + os.Getenv("INFO_TITLE_BASE"),
		"pageDesc":     os.Getenv("INFO_DESC_BASE"),
		"pageScope":    "about",
		"ogTags": fiber.Map{
			"title": "О компании - " + os.Getenv("INFO_TITLE_BASE"),
			"type":  "website",
		},
		"activeAbout": true,
		"data":        data,
	})

	return err
}
