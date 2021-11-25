package routes

import (
	"os"

	"github.com/gofiber/fiber/v2"

	"ivankprod.ru/src/server/internal/models"
	"ivankprod.ru/src/server/internal/utils"
)

func RouteContactsIndex(c *fiber.Ctx) error {
	uAuth, ok := c.Locals("user_auth").(*models.User)
	if !ok {
		uAuth = nil
	}

	data := make(fiber.Map)

	if uAuth != nil {
		data["user"] = *uAuth
	}

	err := c.Render("contacts", fiber.Map{
		"urlBase":      c.BaseURL(),
		"urlCanonical": c.BaseURL() + c.Path(),
		"pageTitle":    "Контакты - " + os.Getenv("INFO_TITLE_BASE"),
		"pageDesc":     os.Getenv("INFO_DESC_BASE"),
		"pageScope":    "contacts",
		"ogTags": fiber.Map{
			"title": "Контакты - " + os.Getenv("INFO_TITLE_BASE"),
			"type":  "website",
		},
		"activeContacts": true,
		"data":           data,
	})

	if err == nil {
		if os.Getenv("STAGE_MODE") == "dev" {
			go utils.DevLogger(c.Request().URI().String(), c.IP(), 200)
		}
	}

	return err
}
