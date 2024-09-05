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
)

const (
	shutdownTimeout = 5 * time.Second
)

func main() {
	app := setupApp()

	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)

	var serverShutdown sync.WaitGroup

	go func() {
		_ = <-sigChannel
		serverShutdown.Add(1)
		defer serverShutdown.Done()
		_ = app.ShutdownWithTimeout(shutdownTimeout)
	}()

	// todo: get port from env
	log.Println("Server starting...")
	if err := app.Listen(":3000", fiber.ListenConfig{EnablePrefork: false}); err != nil {
		log.Panic(err)
	}
	serverShutdown.Wait()
	log.Println("Server shutdown...")
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

/*
func startWithGracefulShutdown(app *fiber.App) {
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("Http server error: %v", err)
			}
		}
		log.Println("Stopped serving new connections.")
	}()
	log.Print("Server Started")

	<-sigChannel

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), shutdownTimeout)
	defer shutdownRelease()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
	log.Println("Graceful shutdown complete.")
}
*/
