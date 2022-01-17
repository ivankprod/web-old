package monitor

import (
	"crypto/tls"
	"sync"

	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/proxy"
	wsproxy "github.com/yeqown/fasthttp-reverse-proxy/v2"
)

var (
	proxyServer *wsproxy.WSReverseProxy
	once        sync.Once
)

func HandlerGrafana(c *fiber.Ctx) error {
	proxy.WithTlsConfig(&tls.Config{InsecureSkipVerify: true})

	if err := proxy.Do(c, "https://grafana:3000"+c.OriginalURL()); err != nil {
		return err
	}

	c.Response().Header.Del(fiber.HeaderServer)
	return nil
}

func HandlerGrafanaWSProxy(c *fiber.Ctx) error {
	once.Do(func() {
		var err error

		dialer := &websocket.Dialer{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}

		proxyServer, err = wsproxy.NewWSReverseProxyWith(
			wsproxy.WithURL_OptionWS("wss://grafana:3000/admin/monitor/grafana/api/live/ws"),
			wsproxy.WithDialer_OptionWS(dialer),
		)

		if err != nil {
			panic(err)
		}
	})

	proxyServer.ServeHTTP(c.Context())
	return nil
}
