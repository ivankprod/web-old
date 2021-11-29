package router

import (
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tarantool/go-tarantool"

	"ivankprod.ru/src/server/internal/admin"
	"ivankprod.ru/src/server/internal/auth"
	"ivankprod.ru/src/server/internal/models"
	"ivankprod.ru/src/server/internal/monitor"
	"ivankprod.ru/src/server/internal/routes"
	"ivankprod.ru/src/server/internal/utils"
)

// All errors
func HandleError(c *fiber.Ctx, err error) error {
	uAuth, ok := c.Locals("user_auth").(*models.User)
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

	if os.Getenv("STAGE_MODE") == "dev" {
		go utils.DevLogger(c.Request().URI().String(), c.IP(), code)
	}

	return rerr
}

// Router
func Router(app *fiber.App /*dbm *sqlx.DB,*/, dbt *tarantool.Connection, sitemap *string) {
	// User authentication
	app.Use(func(c *fiber.Ctx) error {
		if c.Cookies("session") != "" {
			auth, err := models.IsAuthenticated(dbt, c.Cookies("session"), c.Get("user-agent"))
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

				// Update user last access time
				go func(dbt *tarantool.Connection, id uint64) {
					_ = models.UpdateUserAccessTime(dbt, id)
				}(dbt, auth.ID)
			}
		}

		return c.Next()
	})

	// Admin
	adminGroup := app.Group("/admin/", auth.WebmasterAdministratorAccess)
	adminGroup.Get("/", admin.RouteAdminIndex)

	// Monitoring routes
	adminGroup.Group("/monitor/prometheus/", auth.WebmasterAccess, monitor.RoutePrometheus)
	adminGroup.Group("/monitor/grafana/", auth.WebmasterAccess, monitor.RouteGrafana)

	// Routes
	app.Get("/", routes.RouteHomeIndex)
	app.Get("/projects/", routes.RouteProjectsIndex)
	app.Get("/projects/:type/", routes.RouteProjectsView)
	app.Get("/services/", routes.RouteServicesIndex)
	app.Get("/blog/", routes.RouteBlogIndex)
	app.Get("/about/", routes.RouteAboutIndex)
	app.Get("/contacts/", routes.RouteContactsIndex)
	app.Get("/auth/", auth.RouteAuthIndex(dbt))
	app.Get("/auth/logout/", auth.RouteAuthLogout)

	// Sitemap route
	app.Get("/sitemap/", func(c *fiber.Ctx) error {
		uAuth, ok := c.Locals("user_auth").(*models.User)
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
		if err == nil {
			if os.Getenv("STAGE_MODE") == "dev" {
				go utils.DevLogger(c.Request().URI().String(), c.IP(), 200)
			}

			return nil
		}

		return fiber.NewError(fiber.StatusNotFound, "Запрашиваемая страница не найдена либо ещё не создана")
	})

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
				"code":   fiber.StatusNotFound,
				"status": utils.ErrorStatus["404"],
				"msg":    "Запрашиваемая страница не найдена либо ещё не создана",
			},
			"data": data,
		})

		if os.Getenv("STAGE_MODE") == "dev" {
			go utils.DevLogger(c.Request().URI().String(), c.IP(), fiber.StatusNotFound)
		}

		return err
	})
}