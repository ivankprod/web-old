package routes

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"ivankprod.ru/src/server/modules/models"
	"ivankprod.ru/src/server/modules/utils"
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
			"title": "Блог - " + os.Getenv("INFO_TITLE_BASE"),
			"type":  "website",
		},
		"activeBlog": true,
		"data":       data,
	})

	if err == nil {
		if os.Getenv("STAGE_MODE") == "dev" {
			go utils.DevLogger(c.Request().URI().String(), c.IP(), 200)
		}
	}

	return err
}
