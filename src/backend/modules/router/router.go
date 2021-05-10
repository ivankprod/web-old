package router

import (
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"ivankprod.ru/src/backend/modules/models"
	"ivankprod.ru/src/backend/modules/routes"
	"ivankprod.ru/src/backend/modules/utils"
)

// All errors
func HandleError(c *fiber.Ctx, err error) error {
	uAuth := c.Locals("user_auth")
	data := make(fiber.Map)

	if !utils.IsEmptySctruct(uAuth) {
		data = fiber.Map{"user": uAuth}
	}

	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	c.Status(code)
	strCode := strconv.Itoa(code)

	return c.Render("error", fiber.Map{
		"pageTitle": strCode + " - " + os.Getenv("INFO_TITLE_BASE"),
		"pageDesc":  os.Getenv("INFO_DESC_BASE"),
		"error": fiber.Map{
			"code": strCode,
			"msg":  err.Error(),
		},
		"data": data,
	})
}

// Router
func Router(app *fiber.App) {
	// Authentication
	app.Use(func(c *fiber.Ctx) error {
		if c.Cookies("session") != "" {
			auth, err := models.IsAuthenticated(c.Cookies("session"), c.Get("user-agent"))
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, "Ошибка модуля аутентификации!")
			}

			if auth != nil {
				c.Locals("user_auth", (*auth))
			}
		}

		return c.Next()
	})

	app.Get("/", routes.RouteHomeIndex)
	app.Get("/projects/", routes.RouteProjectsIndex)
	app.Get("/projects/:type/", routes.RouteProjectsView)
	app.Get("/auth/", routes.RouteAuthIndex)

	// 404 error
	app.Use(func(c *fiber.Ctx) error {
		uAuth := c.Locals("user_auth")
		data := make(fiber.Map)

		if !utils.IsEmptySctruct(uAuth) {
			data = fiber.Map{"user": uAuth}
		}

		c.Status(fiber.StatusNotFound)

		return c.Render("error", fiber.Map{
			"pageTitle": "404 - " + os.Getenv("INFO_TITLE_BASE"),
			"pageDesc":  os.Getenv("INFO_DESC_BASE"),
			"error": fiber.Map{
				"code": fiber.StatusNotFound,
				"msg":  "Страница не найдена!",
			},
			"data": data,
		})
	})
}
