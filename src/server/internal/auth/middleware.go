package auth

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/ivankprod/ivankprod.ru/src/server/internal/domain"
	"github.com/ivankprod/ivankprod.ru/src/server/internal/services"
	"github.com/ivankprod/ivankprod.ru/src/server/pkg/utils"
)

// Page access middleware
func Access(roles ...uint64) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		uAuth, ok := c.Locals("user_auth").(*domain.User)
		if !ok {
			uAuth = nil
		}

		if uAuth == nil || !utils.Contains(uAuth.Role, roles...) {
			return fiber.NewError(fiber.StatusForbidden, "Доступ к запрашиваемой странице запрещен")
		}

		return c.Next()
	}
}

func Middleware(service services.UserService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Cookies("session") != "" {
			auth, err := service.IsAuthenticated(c.Cookies("session"), c.Get("user-agent"))
			if err != nil {
				return err
			}

			if auth != nil {
				c.Locals("user_auth", auth)
				c.Cookie(&fiber.Cookie{
					Name:     "session",
					Value:    c.Cookies("session"),
					Path:     "/",
					MaxAge:   86400 * 7,
					Expires:  time.Now().Add(time.Hour * 168),
					Secure:   true,
					HTTPOnly: true,
					SameSite: "Lax",
				})

				go func(s services.UserService, id uint64) {
					_, _ = s.UpdateLastAccessTime(id)
				}(service, auth.ID)
			}
		}

		return c.Next()
	}
}
