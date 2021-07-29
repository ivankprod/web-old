package main

import (
	"crypto/tls"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/helmet/v2"
	"github.com/gofiber/template/handlebars"
	"github.com/joho/godotenv"
	"github.com/markbates/pkger"

	"ivankprod.ru/src/server/modules/router"
	"ivankprod.ru/src/server/modules/utils"
)

var (
	MODE_DEV  bool
	MODE_PROD bool
)

func main() {
	// Load .env configuration
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Load STAGE_MODE configuration
	if os.Getenv("STAGE_MODE") == "dev" {
		MODE_DEV = true
		MODE_PROD = false
	} else if os.Getenv("STAGE_MODE") == "prod" {
		MODE_DEV = false
		MODE_PROD = true
	}

	// Sitemap JSON file
	fileSitemapJSON, err := pkger.Open("/misc/sitemap.json")
	if err != nil {
		log.Fatalf("Error opening sitemap.json file: %v", err)
	}

	infoSitemapJSON, err := fileSitemapJSON.Stat()
	if err != nil {
		log.Fatalf("Error reading sitemap.json file: %v", err)
	}

	bytesSitemapJSON := make([]byte, infoSitemapJSON.Size())
	_, err = fileSitemapJSON.Read(bytesSitemapJSON)
	if err != nil {
		log.Fatalf("Error reading sitemap.json file: %v", err)
	}

	// App & template engine
	views := pkger.Dir("/views")
	engine := handlebars.NewFileSystem(views, ".hbs")
	app := fiber.New(fiber.Config{
		Prefork:       false,
		ErrorHandler:  router.HandleError,
		Views:         engine,
		StrictRouting: true,
	})

	// Safe panic
	app.Use(recover.New())

	// Logger
	f, err := os.OpenFile("./logs/"+utils.DateMSK_ToLocaleSepString()+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v", err)
	}
	defer f.Close()

	app.Use(logger.New(logger.Config{
		Format:     "IP: ${ip} | TIME: ${time} | STATUS: ${status}\nURL: ${protocol}://${host}${url}\n\n",
		TimeFormat: "02.01.2006 15:04:05",
		TimeZone:   "Russia/Moscow",
		Output:     f,
	}))

	// ContentSecurityPolicy
	var csp string

	if MODE_DEV {
		csp = "default-src 'self'; base-uri 'self'; block-all-mixed-content; font-src 'self' https: data:; frame-ancestors 'self'; img-src 'self' *.userapi.com *.fbsbx.com *.yandex.net *.googleusercontent.com data:; object-src 'none'; script-src 'self' 'unsafe-inline' 'unsafe-eval'; script-src-attr 'none'; style-src 'self' https: 'unsafe-inline'; upgrade-insecure-requests"
	} else if MODE_PROD {
		csp = "default-src 'self'; base-uri 'self'; block-all-mixed-content; font-src 'self' https: data:; frame-ancestors 'self'; img-src 'self' *.userapi.com *.fbsbx.com *.yandex.net *.googleusercontent.com data:; object-src 'none'; script-src 'self' 'unsafe-inline'; script-src-attr 'none'; style-src 'self' https: 'unsafe-inline'; upgrade-insecure-requests"
	}

	app.Use(helmet.New(helmet.Config{
		ContentSecurityPolicy: csp,
	}))

	// Compression
	app.Use(compress.New(compress.Config{Level: compress.LevelBestSpeed}))

	// HTTP->HTTPS, without www & sitemap.xml
	app.Use(func(c *fiber.Ctx) error {
		if c.Protocol() == "http" || (c.Subdomains() != nil && c.Subdomains(0)[0] == "www") {
			return c.Redirect("https://"+os.Getenv("SERVER_HOST")+c.OriginalURL(), 301)
		}

		if c.OriginalURL() == "/sitemap.xml" {
			return c.SendFile("./sitemap.xml", true)
		}

		return c.Next()
	})

	// favicon
	app.Use(favicon.New(favicon.Config{
		File: "./favicon.ico",
	}))

	// Static files
	app.Static("/static/", "./static", fiber.Static{Compress: true, MaxAge: 86400})

	// Setup router
	router.Router(app, &bytesSitemapJSON)

	// HTTP listener
	go func() {
		log.Fatal(app.Listen(os.Getenv("SERVER_HOST") + ":" + os.Getenv("SERVER_PORT_HTTP")))
	}()

	// HTTPS certs
	cer, err := tls.LoadX509KeyPair("./certs/ivankprod.ru/ivankprod.crt", "./certs/ivankprod.ru/ivankprod.key")
	if err != nil {
		log.Fatal(err)
	}

	// HTTPS listener
	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	ln, err := tls.Listen("tcp", os.Getenv("SERVER_HOST")+":"+os.Getenv("SERVER_PORT_HTTPS"), config)
	if err != nil {
		panic(err)
	}

	// LISTEN
	log.Fatal(app.Listener(ln))
}
