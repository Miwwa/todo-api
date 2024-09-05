package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todo-api/config"
	"todo-api/utils"
)

const (
	shutdownTimeout = 5 * time.Second
	idleTimeout     = 5 * time.Second
	readTimeout     = 5 * time.Second
	writeTimeout    = 5 * time.Second
)

func main() {
	appConfig := config.Default()

	app := setupApp()
	startWithGracefulShutdown(app, appConfig)
}

func setupApp() *fiber.App {
	app := fiber.New(fiber.Config{
		IdleTimeout:  idleTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		ErrorHandler: utils.JsonErrorHandler,
	})

	app.Use(recover.New())
	app.Use(logger.New())

	app.Get(healthcheck.DefaultLivenessEndpoint, healthcheck.NewHealthChecker())
	app.Get(healthcheck.DefaultReadinessEndpoint, healthcheck.NewHealthChecker(healthcheck.Config{Probe: func(ctx fiber.Ctx) bool {
		return true
	}}))

	app.Get("/", func(ctx fiber.Ctx) error {
		_, err := ctx.WriteString("Hello World")
		if err != nil {
			return err
		}
		return nil
	})

	app.Get("/501", func(ctx fiber.Ctx) error {
		return fiber.ErrNotImplemented
	})

	app.Use(utils.Json404)

	return app
}

func startWithGracefulShutdown(app *fiber.App, config config.AppConfig) {
	address := config.Address()
	fiberConfig := fiber.ListenConfig{EnablePrefork: config.IsProd()}

	go func() {
		if err := app.Listen(address, fiberConfig); err != nil {
			log.Panic(err)
		}
	}()

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)

	_ = <-sigChannel
	log.Println("Server shutdown...")
	_ = app.ShutdownWithTimeout(shutdownTimeout)

	// run db.close() etc.
}
