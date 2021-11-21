package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/favicon"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/helmet/v2"
	"github.com/gofiber/template/handlebars"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/markbates/pkger"
	"github.com/markbates/pkger/pkging"
	"github.com/tarantool/go-tarantool"

	"ivankprod.ru/src/server/modules/db"
	"ivankprod.ru/src/server/modules/router"
	"ivankprod.ru/src/server/modules/utils"
)

var (
	MODE_DEV  bool
	MODE_PROD bool
)

// App struct
type App struct {
	*fiber.App

	DBM *sqlx.DB
	DBT *tarantool.Connection
}

// Sitemap JSON to HTML
func loadSitemap(fileSitemapJSON *pkging.File) *string {
	infoSitemapJSON, err := (*fileSitemapJSON).Stat()
	if err != nil {
		log.Fatalf("Error reading sitemap.json file: %v", err)
	}

	bytesSitemapJSON := make([]byte, infoSitemapJSON.Size())
	_, err = (*fileSitemapJSON).Read(bytesSitemapJSON)
	if err != nil {
		log.Fatalf("Error reading sitemap.json file: %v", err)
	}

	sitemap := &utils.Sitemap{}
	err = json.Unmarshal(bytesSitemapJSON, sitemap)
	if err != nil {
		log.Fatalf("Error unmarshalling sitemap.json file: %v", err)
	}

	return sitemap.Nest().ToHTMLString()
}

func main() {
	// Logging file
	f, err := os.OpenFile("./logs/"+utils.DateMSK_ToLocaleSepString()+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening log file: %v\n", err)
	}
	defer f.Close()

	// Server base logging
	log.SetFlags(0)
	log.SetPrefix(utils.TimeMSK_ToLocaleString() + " ")
	log.SetOutput(f)
	log.Println("-- Server starting...")

	// Load .env configuration
	err = godotenv.Load(".env")
	if err != nil {
		log.SetPrefix(utils.TimeMSK_ToLocaleString() + " ")
		log.Println("Error loading .env file")
		log.Fatalln("-- Server starting failed")
	}

	// Load STAGE_MODE configuration
	if os.Getenv("STAGE_MODE") == "dev" {
		MODE_DEV = true
		MODE_PROD = false
	} else if os.Getenv("STAGE_MODE") == "prod" {
		MODE_DEV = false
		MODE_PROD = true
	}

	// Open sitemap.json file for reading
	fileSitemapJSON, err := pkger.Open("/misc/sitemap.json")
	if err != nil {
		log.SetPrefix(utils.TimeMSK_ToLocaleString() + " ")
		log.Printf("Error opening sitemap.json file: %v\n", err)
		log.Fatalln("-- Server starting failed")
	}

	// Templates engine
	engine := handlebars.New("./views", ".hbs")

	// App
	app := App{
		App: fiber.New(fiber.Config{
			Prefork:       false,
			ErrorHandler:  router.HandleError,
			Views:         engine,
			StrictRouting: true,
		}),
	}

	// DB MySQL connect
	/*dbm, err := db.ConnectMySQL()
	if dbm == nil {
		log.SetPrefix(utils.TimeMSK_ToLocaleString() + " ")

		if err == nil {
			log.Println("Failed connecting to MySQL database")
		} else {
			log.Printf("Error connecting to MySQL database: %v\n", err)
		}

		app.fail("-- Server starting failed\n")
	} else {
		app.DBM = dbm
	}*/

	// DB Tarantool connect
	dbt, err := db.ConnectTarantool()
	if dbt == nil {
		log.SetPrefix(utils.TimeMSK_ToLocaleString() + " ")

		if err == nil {
			log.Println("Failed connecting to Tarantool database")
		} else {
			log.Printf("Error connecting to Tarantool database: %v\n", err)
		}

		app.fail("-- Server starting failed\n")
	} else {
		app.DBT = dbt
	}

	// Safe panic
	app.Use(recover.New())

	// Logger
	app.Use(logger.New(logger.Config{
		Format:     "${method} | IP: ${ip} | TIME: ${time} UTC | STATUS: ${status}\nURL: ${protocol}://${host}${url}\n\n",
		TimeFormat: "02.01.2006 15:04:05",
		TimeZone:   "Europe/Moscow",
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

	// Prometheus
	prometheus := fiberprometheus.New("ivankprodru_app")
	prometheus.RegisterAt(app.App, "/metrics")
	app.Use(prometheus.Middleware)

	// Static files
	app.Static("/static/", "./static", fiber.Static{Compress: true, MaxAge: 86400})

	// Setup router
	router.Router(app.App /*, app.DBM */, app.DBT, loadSitemap(&fileSitemapJSON))

	// HTTP listener
	go func() {
		log.SetPrefix(utils.TimeMSK_ToLocaleString() + " ")
		log.Printf("-- Attempt starting at %s:%s\n", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT_HTTP"))

		if err := app.Listen( /*os.Getenv("SERVER_HOST") +*/ ":" + os.Getenv("SERVER_PORT_HTTP")); err != nil {
			log.SetPrefix(utils.TimeMSK_ToLocaleString() + " ")
			log.Println(err)
			app.fail(fmt.Sprintf("-- Server starting at %s:%s failed\n\n", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT_HTTP")))
		}
	}()

	// HTTPS certs
	cer, err := tls.LoadX509KeyPair(os.Getenv("SERVER_SSL_CERT"), os.Getenv("SERVER_SSL_KEY"))
	if err != nil {
		log.SetPrefix(utils.TimeMSK_ToLocaleString() + " ")
		log.Println(err)
		app.fail("-- Server starting failed\n")
	}

	// HTTPS listener
	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	ln, err := tls.Listen("tcp" /*os.Getenv("SERVER_HOST")+*/, ":"+os.Getenv("SERVER_PORT_HTTPS"), config)
	if err != nil {
		log.SetPrefix(utils.TimeMSK_ToLocaleString() + " ")
		log.Println(err)
		app.fail("-- Server starting failed\n")
	}

	// Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-c
		app.exit()
	}()

	// LISTEN
	log.SetPrefix(utils.TimeMSK_ToLocaleString() + " ")
	log.Println("-- Attempt starting at " + os.Getenv("SERVER_HOST") + ":" + os.Getenv("SERVER_PORT_HTTPS") + "\n")

	if err = app.Listener(ln); err != nil {
		log.SetPrefix(utils.TimeMSK_ToLocaleString() + " ")
		log.Println(err)
		app.fail(fmt.Sprintf("-- Server starting at %s:%s failed\n\n", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT_HTTPS")))
	}
}

// App fail
func (app *App) fail(msg ...string) {
	if app.DBM != nil {
		app.DBM.Close()
	}
	if app.DBT != nil {
		app.DBT.Close()
	}

	log.SetPrefix(utils.TimeMSK_ToLocaleString() + " ")

	if len(msg) > 0 {
		log.Fatalln(msg[0])
	} else {
		os.Exit(1)
	}
}

// App exit
func (app *App) exit(msg ...string) {
	if app.DBM != nil {
		app.DBM.Close()
	}
	if app.DBT != nil {
		app.DBT.Close()
	}

	log.SetPrefix(utils.TimeMSK_ToLocaleString() + " ")

	if len(msg) > 0 {
		log.Println(msg[0])
	}

	log.Print("-- Server terminated\n\n")
	_ = app.Shutdown()
}
