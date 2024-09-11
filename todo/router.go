package todo

import (
	"github.com/gofiber/fiber/v3"
	"todo-api/config"
	"todo-api/user"
	"todo-api/utils"
)

func SetupRoutes(app *fiber.App, config *config.AppConfig, storage Storage, validator *utils.AppValidator) {
	todoGroup := app.Group("/todos", user.ValidateAndExtractTokenMiddleware(config.JwtSecret))
	todoGroup.Post("/", CreateHandler(storage, validator))
	todoGroup.Get("/", ReadHandler(storage, validator))
	todoGroup.Put("/:id", UpdateHandler(storage, validator))
	todoGroup.Delete("/:id", DeleteHandler(storage))
}

type dto struct {
	Id          Id     `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func CreateHandler(storage Storage, validator *utils.AppValidator) fiber.Handler {
	type CreateRequest struct {
		Title       string `json:"title" validate:"required,gte=0,lte=255"`
		Description string `json:"description" validate:"lte=100000"`
	}

	type CreateResponse dto

	return func(ctx fiber.Ctx) error {
		u := user.FromContext(ctx)

		req := CreateRequest{}
		err := ctx.Bind().Body(&req)
		if err != nil {
			return fiber.ErrBadRequest
		}
		if err = validator.Validate(req); err != nil {
			return err
		}

		todo, err := storage.Create(ctx.Context(), u.Id, req.Title, req.Description)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return ctx.JSON(CreateResponse{
			Id:          todo.Id,
			Title:       todo.Title,
			Description: todo.Description,
		})
	}
}

func ReadHandler(storage Storage, validator *utils.AppValidator) fiber.Handler {
	type ReadRequest struct {
		Page  uint `json:"page" validate:"gt=0"`
		Limit uint `json:"limit" validate:"gt=0"`
	}

	type ReadResponse struct {
		Data  []dto
		Page  uint `json:"page"`
		Limit uint `json:"limit"`
		Total uint `json:"total"`
	}

	return func(ctx fiber.Ctx) error {
		u := user.FromContext(ctx)
		req := ReadRequest{
			Page:  1,
			Limit: 10,
		}
		err := ctx.Bind().Query(&req)
		if err != nil {
			return fiber.ErrBadRequest
		}
		if err = validator.Validate(req); err != nil {
			return err
		}

		todos, err := storage.GetByUserId(ctx.Context(), u.Id, req.Limit, (req.Page-1)*req.Limit)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		total, err := storage.Count(ctx.Context(), u.Id)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		response := ReadResponse{
			Data:  make([]dto, len(todos)),
			Page:  req.Page,
			Limit: req.Limit,
			Total: total,
		}

		for i, todo := range todos {
			response.Data[i] = dto{
				Id:          todo.Id,
				Title:       todo.Title,
				Description: todo.Description,
			}
		}

		return ctx.JSON(response)
	}
}

func UpdateHandler(storage Storage, validator *utils.AppValidator) fiber.Handler {
	type UpdateRequest struct {
		Id          Id     `json:"id" validate:"required"`
		Title       string `json:"title" validate:"required,gte=0,lte=255"`
		Description string `json:"description" validate:"lte=100000"`
	}

	type UpdateResponse dto

	return func(ctx fiber.Ctx) error {
		req := UpdateRequest{}
		req.Id = Id(ctx.Params("id", ""))
		err := ctx.Bind().Body(&req)
		if err != nil {
			return fiber.ErrBadRequest
		}
		if err = validator.Validate(req); err != nil {
			return err
		}

		err = validatePermission(ctx, storage, req.Id)
		if err != nil {
			return err
		}

		todo, err := storage.Update(ctx.Context(), req.Id, req.Title, req.Description)
		if err != nil {
			return fiber.ErrInternalServerError
		}
		if todo.Invalid() {
			return fiber.ErrForbidden
		}

		return ctx.JSON(UpdateResponse{
			Id:          todo.Id,
			Title:       todo.Title,
			Description: todo.Description,
		})
	}
}

func DeleteHandler(storage Storage) fiber.Handler {
	type DeleteRequest struct {
	}

	type DeleteResponse struct {
	}

	return func(ctx fiber.Ctx) error {
		todoId := Id(ctx.Params("id", ""))

		err := validatePermission(ctx, storage, todoId)
		if err != nil {
			return err
		}

		err = storage.Delete(ctx.Context(), todoId)
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return ctx.JSON(DeleteResponse{})
	}
}

func validatePermission(ctx fiber.Ctx, storage Storage, todoId Id) error {
	u := user.FromContext(ctx)
	todo, err := storage.GetById(ctx.Context(), todoId)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	if todo.UserId != u.Id {
		return fiber.ErrForbidden
	}
	return nil
}
