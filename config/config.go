package config

import (
	"errors"
	"fmt"
	"github.com/caarlos0/env/v11"
	"os"
)

type AppConfig struct {
	IsProduction bool   `env:"IS_PRODUCTION" envDefault:"false"`
	Port         int    `env:"PORT" envDefault:"3000"`
	SqliteDbPath string `env:"SQLITE_DB_PATH" envDefault:"./db.sqlite"`
	JwtSecret    string `env:"JWT_SECRET" envDefault:"mySecret"`
}

func (c *AppConfig) IsDev() bool {
	return !c.IsProduction
}

func (c *AppConfig) Address() string {
	return fmt.Sprintf(":%d", c.Port)
}

func (c *AppConfig) DebugString() string {
	return fmt.Sprintf("IsProduction: %v\nPort: %d\ndb: %s\n", c.IsProduction, c.Port, c.SqliteDbPath)
}

func (c *AppConfig) Validate() error {
	if c.Port < 1 || c.Port > 65535 {
		return errors.New("port must in range 1-65535")
	}

	if _, err := os.Stat(c.SqliteDbPath); os.IsNotExist(err) {
		return fmt.Errorf("sqlite db path %s does not exist", c.SqliteDbPath)
	}

	if c.IsProduction && c.JwtSecret == "" {
		return errors.New("jwt secret is required in production mode, set JWT_SECRET environment variable")
	}

	return nil
}

func FromEnv() (AppConfig, error) {
	config := AppConfig{}
	err := env.ParseWithOptions(&config, env.Options{
		RequiredIfNoDef: true,
		Prefix:          "",
	})
	if err != nil {
		return config, err
	}
	return config, nil
}
