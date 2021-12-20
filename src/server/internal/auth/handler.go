package auth

import (
	"os"
	"time"

	"github.com/gofiber/fiber/v2"

	"ivankprod.ru/src/server/internal/domain"
	"ivankprod.ru/src/server/internal/services"
	"ivankprod.ru/src/server/pkg/utils"
)

type AuthHandler struct {
	userService services.UserService
}

func NewAuthHandler(s services.UserService) *AuthHandler {
	return &AuthHandler{
		userService: s,
	}
}

func (h *AuthHandler) HandlerIndex(c *fiber.Ctx) error {
	uAuth, ok := c.Locals("user_auth").(*domain.User)
	if !ok {
		uAuth = nil
	}

	data := make(fiber.Map)
	title := "Авторизация"

	if c.Query("code") != "" && c.Query("state") != "" {
		if c.Query("state") == "vk" {
			if err := authVK(h.userService, c, uAuth); err != nil {
				return err
			}
		} else if c.Query("state") == "facebook" {
			if err := authFacebook(h.userService, c, uAuth); err != nil {
				return err
			}
		} else if c.Query("state") == "yandex" {
			if err := authYandex(h.userService, c, uAuth); err != nil {
				return err
			}
		} else if c.Query("state") == "google" {
			if err := authGoogle(h.userService, c, uAuth); err != nil {
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

			userAccounts, err := h.userService.FindGroup(uAuth.Group)
			if err != nil {
				return err
			}

			if userAccounts != nil {
				data["user_accounts"] = userAccounts.GetCondsByType(uAuth.Type)
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

func (h *AuthHandler) HandlerLogout(c *fiber.Ctx) error {
	uAuth, ok := c.Locals("user_auth").(*domain.User)
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
