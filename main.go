package main

import (
	"database/sql"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"todo-api/config"
	"todo-api/user"
	"todo-api/utils"
)

const (
	shutdownTimeout = 5 * time.Second
	idleTimeout     = 5 * time.Second
	readTimeout     = 5 * time.Second
	writeTimeout    = 5 * time.Second
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	appConfig := config.FromEnv()
	if appConfig.IsProd() && appConfig.JwtSecret() == "" {
		log.Fatal("JWT_SECRET environment variable is required in production mode")
	}
	log.Printf("config loaded:\n%s\n", appConfig.DebugString())

	db, err := sql.Open("sqlite3", appConfig.DbConnectionString())
	if err != nil {
		log.Fatal(err)
	}

	app := setupApp(&appConfig, db)
	startWithGracefulShutdown(app, db, appConfig)
}

func setupApp(config *config.AppConfig, db *sql.DB) *fiber.App {
	app := fiber.New(fiber.Config{
		IdleTimeout:  idleTimeout,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		ErrorHandler: utils.JsonErrorHandler,
	})

	app.Use(recover.New())
	app.Use(logger.New())

	// healthcheck api
	app.Get(healthcheck.DefaultLivenessEndpoint, healthcheck.NewHealthChecker())
	app.Get(healthcheck.DefaultReadinessEndpoint, healthcheck.NewHealthChecker(healthcheck.Config{Probe: func(ctx fiber.Ctx) bool {
		return true
	}}))

	usersStorage := user.NewSqliteUsersStorage(db)

	// user api
	user.SetupRoutes(app, config, usersStorage)
	app.Get("/user", func(c fiber.Ctx) error {
		user := c.Locals(user.ContextKey).(*user.User)
		return c.JSON(user)
	}, user.ValidateAndExtractTokenMiddleware(config.JwtSecret()))

	app.Use(utils.Json404)

	return app
}

func startWithGracefulShutdown(app *fiber.App, db *sql.DB, config config.AppConfig) {
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

	err := db.Close()
	if err != nil {
		log.Panic(err)
	}
}
