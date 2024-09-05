package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	recover2 "github.com/gofiber/fiber/v3/middleware/recover"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"todo-api/config"
)

const (
	shutdownTimeout = 5 * time.Second
)

func main() {
	appConfig := config.Default()

	app := setupApp()
	startWithGracefulShutdown(app, appConfig)
}

func setupApp() *fiber.App {
	app := fiber.New(fiber.Config{})

	app.Use(recover2.New())
	app.Use(logger.New())

	app.Get("/", func(ctx fiber.Ctx) error {
		_, err := ctx.WriteString("Hello World")
		if err != nil {
			return err
		}
		return nil
	})

	return app
}

func startWithGracefulShutdown(app *fiber.App, config config.AppConfig) {
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)

	var serverShutdown sync.WaitGroup

	go func() {
		_ = <-sigChannel
		serverShutdown.Add(1)
		defer serverShutdown.Done()
		_ = app.ShutdownWithTimeout(shutdownTimeout)
	}()

	log.Println("Server starting...")

	address := config.Address()
	fiberConfig := fiber.ListenConfig{EnablePrefork: config.IsProd()}
	if err := app.Listen(address, fiberConfig); err != nil {
		log.Panic(err)
	}
	serverShutdown.Wait()

	log.Println("Server shutdown...")
}
