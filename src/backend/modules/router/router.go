package router

import (
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"ivankprod.ru/src/backend/modules/models"
	"ivankprod.ru/src/backend/modules/routes"
	"ivankprod.ru/src/backend/modules/utils"
)

// All errors
func HandleError(c *fiber.Ctx, err error) error {
	uAuth, ok := c.Locals("user_auth").(*models.User)
	if !ok {
		uAuth = nil
	}

	data := make(fiber.Map)
	msgPrefix := ""

	if uAuth != nil {
		data["user"] = *uAuth
	}

	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	if code != 404 {
		msgPrefix = "Ошибка: "
	}

	c.Status(code)
	strCode := strconv.Itoa(code)

	rerr := c.Render("error", fiber.Map{
		"pageTitle": strCode + " - " + os.Getenv("INFO_TITLE_BASE"),
		"pageDesc":  os.Getenv("INFO_DESC_BASE"),
		"error": fiber.Map{
			"code": strCode,
			"msg":  msgPrefix + err.Error(),
		},
		"data": data,
	})

	go utils.Logger(c.Request().URI().String(), c.IP(), code)

	return rerr
}

// Router
func Router(app *fiber.App) {
	// Authentication
	app.Use(func(c *fiber.Ctx) error {
		if c.Cookies("session") != "" {
			auth, err := models.IsAuthenticated(c.Cookies("session"), c.Get("user-agent"))
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

				// Update access time
				go func(id int64) {
					models.UpdateUserAccessTime(id)
				}((*auth).ID)
			}
		}

		return c.Next()
	})

	// Routes
	app.Get("/", routes.RouteHomeIndex)
	app.Get("/projects/", routes.RouteProjectsIndex)
	app.Get("/projects/:type/", routes.RouteProjectsView)
	app.Get("/services/", routes.RouteServicesIndex)
	app.Get("/blog/", routes.RouteBlogIndex)
	app.Get("/about/", routes.RouteAboutIndex)
	app.Get("/contacts/", routes.RouteContactsIndex)
	app.Get("/auth/", routes.RouteAuthIndex)
	app.Get("/auth/logout/", routes.RouteAuthLogout)

	// 404 error
	app.Use(func(c *fiber.Ctx) error {
		uAuth, ok := c.Locals("user_auth").(*models.User)
		if !ok {
			uAuth = nil
		}

		data := make(fiber.Map)

		if uAuth != nil {
			data["user"] = *uAuth
		}

		c.Status(fiber.StatusNotFound)

		err := c.Render("error", fiber.Map{
			"pageTitle": "404 - " + os.Getenv("INFO_TITLE_BASE"),
			"pageDesc":  os.Getenv("INFO_DESC_BASE"),
			"error": fiber.Map{
				"code": fiber.StatusNotFound,
				"msg":  "Страница не найдена!",
			},
			"data": data,
		})

		go utils.Logger(c.Request().URI().String(), c.IP(), fiber.StatusNotFound)

		return err
	})
}
