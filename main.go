package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	configuration "github.com/jessicapeter01/piggy_bank/config"
	"github.com/jessicapeter01/piggy_bank/database"
)

type App struct {
	*fiber.App

	DB      *database.Database
	Session *session.Store
}

func main() {
	config := configuration.GetInstance()
	app := setupApp(config)

	// Close any connections on interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		app.exit()
	}()

	// Start listening on the specified address
	if err := app.Listen(config.GetString("APP_ADDR")); err != nil {
		log.Panic(err)
	}
}

func setupApp(config *configuration.Config) App {
	database.Setup()

	app := App{
		App:     fiber.New(*config.GetFiberConfig()),
		Session: session.New(config.GetSessionConfig()),
	}

	app.DB = (&database.Database{
		DB: database.DBConn,
	})

	database.SessionStore = app.Session
	app.Session.RegisterType("")
	var typeUint uint = 1
	app.Session.RegisterType(typeUint)
	var typeBool bool = false
	app.Session.RegisterType(typeBool)

	// setup routes
	setupRoutes(app.App)

	return app
}

func setupRoutes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		// Render index
		return c.Render("index", fiber.Map{
			"Title": "Hello, World!",
		})
	})

	// api group
	api := app.Group("/api")

	// give response when at /api
	api.Get("", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"message": "You are at the api endpoint 😉",
		})
	})

	// connect api routes
	// routes.ProductRoute(api.Group("/products"))
}

// Stop the Fiber application
func (app *App) exit() {
	_ = app.Shutdown()
}
