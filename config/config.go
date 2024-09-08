package config

import (
	"fmt"
	"os"
	"strconv"
)

type AppConfig struct {
	port         int
	environment  string
	isProd       bool
	dbConnection string
	jwtSecret    string
}

func Default() AppConfig {
	return AppConfig{
		port:        3000,
		environment: "development",
		isProd:      true,
	}
}

func (c *AppConfig) Port() int {
	return c.port
}

func (c *AppConfig) Address() string {
	return fmt.Sprintf(":%d", c.port)
}

func (c *AppConfig) IsDev() bool {
	return !c.isProd
}

func (c *AppConfig) IsProd() bool {
	return c.isProd
}

func (c *AppConfig) DbConnectionString() string {
	return c.dbConnection
}

func (c *AppConfig) DebugString() string {
	return fmt.Sprintf("port:%d\nenvironment:%s\nisDev:%t\nisProd:%t", c.port, c.environment, c.IsDev(), c.IsProd())
}

func (c *AppConfig) JwtSecret() string {
	return c.jwtSecret
}

func FromEnv() AppConfig {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 3000
	}

	env := os.Getenv("ENV")

	dbConnection := os.Getenv("DB_CONNECTION")
	if dbConnection == "" {
		dbConnection = "./db.sqlite"
	}

	jwtSecret := os.Getenv("JWT_SECRET")

	return AppConfig{
		port:         port,
		environment:  env,
		isProd:       env == "production",
		dbConnection: dbConnection,
		jwtSecret:    jwtSecret,
	}
}
