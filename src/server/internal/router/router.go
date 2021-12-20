package router

import (
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/tarantool/go-tarantool"

	"ivankprod.ru/src/server/internal/admin"
	"ivankprod.ru/src/server/internal/auth"
	"ivankprod.ru/src/server/internal/domain"
	"ivankprod.ru/src/server/internal/handlers"
	"ivankprod.ru/src/server/internal/monitor"
	"ivankprod.ru/src/server/internal/repositories"
	"ivankprod.ru/src/server/internal/services"
	"ivankprod.ru/src/server/pkg/utils"
)

// All errors
func HandleError(c *fiber.Ctx, err error) error {
	uAuth, ok := c.Locals("user_auth").(*domain.User)
	if !ok {
		uAuth = nil
	}

	data := make(fiber.Map)
	msgPrefix := ""
	msgStatus := "Неизвестная ошибка"

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

	if s, ok := utils.ErrorStatus[strCode]; ok {
		msgStatus = s
	}

	rerr := c.Render("error", fiber.Map{
		"pageTitle": strCode + " - " + os.Getenv("INFO_TITLE_BASE"),
		"pageDesc":  os.Getenv("INFO_DESC_BASE"),
		"error": fiber.Map{
			"code":   strCode,
			"status": msgStatus,
			"msg":    msgPrefix + err.Error(),
		},
		"data": data,
	})

	return rerr
}

// Router
func Router(app *fiber.App, dbt *tarantool.Connection, sitemap *string) {
	// User service
	userRepository := repositories.NewUserRepository(dbt)
	userService := services.NewUserService(userRepository)

	// User authentication middleware
	app.Use(auth.Middleware(userService))

	// Monitoring Grafana WebSocket handler
	app.All("/admin/monitor/grafana/api/live/ws", auth.Access(domain.USER_ROLE_WEBMASTER), monitor.HandlerGrafanaWSProxy)

	// Admin handlers
	adminGroup := app.Group("/admin/", auth.Access(domain.USER_ROLE_ADMINISTRATOR, domain.USER_ROLE_WEBMASTER))
	adminGroup.Get("/", admin.RouteAdminIndex)

	// Monitoring handlers
	adminGroup.Group("/monitor/prometheus/", auth.Access(domain.USER_ROLE_WEBMASTER), monitor.HandlerPrometheus)
	adminGroup.Group("/monitor/grafana/", auth.Access(domain.USER_ROLE_WEBMASTER), monitor.HandlerGrafana)

	// Handlers
	app.Get("/", handlers.HandlerHomeIndex)
	app.Get("/projects/", handlers.HandlerProjectsIndex)
	app.Get("/projects/:type/", handlers.HandlerProjectsView)
	app.Get("/services/", handlers.HandlerServicesIndex)
	app.Get("/blog/", handlers.HandlerBlogIndex)
	app.Get("/about/", handlers.HandlerAboutIndex)
	app.Get("/contacts/", handlers.HandlerContactsIndex)
	app.Get("/sitemap/", handlers.HandlerSitemapIndex(sitemap))

	authHandler := auth.NewAuthHandler(userService)
	app.Get("/auth/", authHandler.HandlerIndex)
	app.Get("/auth/logout/", authHandler.HandlerLogout)

	// 404 error
	app.Use(func(c *fiber.Ctx) error {
		uAuth, ok := c.Locals("user_auth").(*domain.User)
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
				"code":   fiber.StatusNotFound,
				"status": utils.ErrorStatus["404"],
				"msg":    "Запрашиваемая страница не найдена либо ещё не создана",
			},
			"data": data,
		})

		return err
	})
}
