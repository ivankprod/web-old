package monitor

import (
	"crypto/tls"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

func RoutePrometheus(c *fiber.Ctx) error {
	proxy.WithTlsConfig(&tls.Config{InsecureSkipVerify: true})

	if err := proxy.Do(c, "https://prometheus:9090"+c.OriginalURL()); err != nil {
		return err
	}

	c.Response().Header.Del(fiber.HeaderServer)
	return nil
}
