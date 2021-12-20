package auth

import (
	"time"

	"github.com/gofiber/fiber/v2"

	"ivankprod.ru/src/server/internal/services"
)

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
