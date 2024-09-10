package config

import (
	"fmt"
	"github.com/caarlos0/env/v11"
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
