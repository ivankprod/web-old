package auth

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tarantool/go-tarantool"

	"ivankprod.ru/src/server/internal/models"
	"ivankprod.ru/src/server/pkg/utils"
)

func HandlerAuthIndex(db *tarantool.Connection) fiber.Handler {
	return func(c *fiber.Ctx) error {
		uAuth, ok := c.Locals("user_auth").(*models.User)
		if !ok {
			uAuth = nil
		}

		data := make(fiber.Map)
		title := "Авторизация"

		if c.Query("code") != "" && c.Query("state") != "" {
			if c.Query("state") == "vk" {
				if err := authVK(c, db, uAuth); err != nil {
					return err
				}
			} else if c.Query("state") == "facebook" {
				if err := authFacebook(c, db, uAuth); err != nil {
					return err
				}
			} else if c.Query("state") == "yandex" {
				if err := authYandex(c, db, uAuth); err != nil {
					return err
				}
			} else if c.Query("state") == "google" {
				if err := authGoogle(c, db, uAuth); err != nil {
					return err
				}
			}

			c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
			return c.SendString("<!DOCTYPE html><html><head><script>window.location.href=\"/auth/\"</script></head><body></body></html>")
		} else {
			if uAuth == nil {
				data["links"] = utils.GetAuthLinks()
			} else {
				data["user"] = *uAuth
				data["links"] = utils.GetAuthLinks()
				title = "Личный кабинет"

				userAccounts, err := models.GetUsersGroup(db, uAuth.Group)
				if err != nil {
					return err
				}

				if userAccounts != nil {
					data["user_accounts"] = (*userAccounts).GetCondsByType(uAuth.Type)
				}
			}
		}

		err := c.Render("auth", fiber.Map{
			"urlBase":      c.BaseURL(),
			"urlCanonical": c.BaseURL() + c.Path(),
			"pageTitle":    title + " - " + os.Getenv("INFO_TITLE_BASE"),
			"pageDesc":     os.Getenv("INFO_DESC_BASE"),
			"ogTags": fiber.Map{
				"title": os.Getenv("INFO_TITLE_BASE"),
				"type":  "website",
			},
			"data": data,
		})

		return err
	}
}

func HandlerAuthLogout(c *fiber.Ctx) error {
	uAuth, ok := c.Locals("user_auth").(*models.User)
	if !ok || uAuth == nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Для того, чтобы выйти из системы, Вы должны быть авторизованы")
	}

	c.Cookie(&fiber.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		MaxAge:   0,
		Expires:  time.Now().Add(-(time.Hour * 1)),
		Secure:   true,
		HTTPOnly: true,
		SameSite: "Lax",
	})

	return c.Redirect("/auth/", 303)
}
