package todo

import (
	"github.com/gofiber/fiber/v3"
	"todo-api/config"
	"todo-api/user"
)

func SetupRoutes(app *fiber.App, config *config.AppConfig, storage Storage) {
	todoGroup := app.Group("/todos", user.ValidateAndExtractTokenMiddleware(config.JwtSecret()))
	todoGroup.Post("/", CreateHandler(config, storage))
	todoGroup.Get("/", ReadHandler(config, storage))
	todoGroup.Put("/:id", UpdateHandler(config, storage))
	todoGroup.Delete("/:id", DeleteHandler(config, storage))
}

type Dto struct {
	Id          Id     `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func CreateHandler(config *config.AppConfig, storage Storage) fiber.Handler {
	type CreateRequest struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	type CreateResponse Dto

	return func(ctx fiber.Ctx) error {
		return fiber.ErrNotImplemented
	}
}

func ReadHandler(config *config.AppConfig, storage Storage) fiber.Handler {
	type ReadRequest struct {
		Page  uint `json:"page"`
		Limit uint `json:"limit"`
	}

	type ReadResponse struct {
		Data  []Dto
		Page  uint `json:"page"`
		Limit uint `json:"limit"`
		Total uint `json:"total"`
	}

	return func(ctx fiber.Ctx) error {
		return fiber.ErrNotImplemented
	}
}

func UpdateHandler(config *config.AppConfig, storage Storage) fiber.Handler {
	type UpdateRequest Dto

	type UpdateResponse Dto

	return func(ctx fiber.Ctx) error {
		return fiber.ErrNotImplemented
	}
}

func DeleteHandler(config *config.AppConfig, storage Storage) fiber.Handler {
	type DeleteRequest struct {
		Id Id `json:"id"`
	}

	type DeleteResponse struct {
	}

	return func(ctx fiber.Ctx) error {
		return fiber.ErrNotImplemented
	}
}
