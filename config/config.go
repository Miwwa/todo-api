package config

import "fmt"

type AppConfig struct {
	port        int
	environment string
	isDev       bool
}

func Default() AppConfig {
	return AppConfig{
		port:        3000,
		environment: "development",
		isDev:       true,
	}
}

func (c *AppConfig) Port() int {
	return c.port
}

func (c *AppConfig) Address() string {
	return fmt.Sprintf(":%d", c.port)
}

func (c *AppConfig) IsDev() bool {
	return c.isDev
}

func (c *AppConfig) IsProd() bool {
	return !c.isDev
}

// todo: read config from dotenv
