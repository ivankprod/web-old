package handlers

import (
	"os"

	"github.com/gofiber/fiber/v2"

	"github.com/ivankprod/ivankprod.ru/src/server/internal/domain"
)

func HandlerSitemapIndex(sitemap *string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		uAuth, ok := c.Locals("user_auth").(*domain.User)
		if !ok {
			uAuth = nil
		}

		data := make(fiber.Map)

		if uAuth != nil {
			data["user"] = *uAuth
		}

		data["sitemap"] = *sitemap

		err := c.Render("sitemap", fiber.Map{
			"urlBase":      c.BaseURL(),
			"urlCanonical": c.BaseURL() + c.Path(),
			"pageTitle":    "Карта сайта - " + os.Getenv("INFO_TITLE_BASE"),
			"pageDesc":     os.Getenv("INFO_DESC_BASE"),
			"pageScope":    "sitemap",
			"ogTags": fiber.Map{
				"title": os.Getenv("INFO_TITLE_BASE"),
			},
			"data": data,
		})

		return err
	}
}
