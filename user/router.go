package user

import (
	"errors"
	"github.com/gofiber/fiber/v3"
	"log"
	"todo-api/config"
	"todo-api/utils"
)

func SetupRoutes(app *fiber.App, config *config.AppConfig, storage Storage, validator *utils.AppValidator) {
	app.Post("/register", Register(config, storage, validator))
	app.Post("/login", Login(config, storage, validator))
}

func Register(config *config.AppConfig, storage Storage, validator *utils.AppValidator) fiber.Handler {
	type RegistrationRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,gte=4,lte=255"`
		Name     string `json:"name" validate:"required,gte=2,lte=255"`
	}

	type RegistrationResponse struct {
		Token string `json:"token"`
	}

	return func(ctx fiber.Ctx) error {
		data := RegistrationRequest{}
		err := ctx.Bind().Body(&data)
		if err != nil {
			return err
		}
		err = validator.Validate(data)
		if err != nil {
			return err
		}

		passwordHash, err := utils.HashPassword(data.Password)
		if err != nil {
			return err
		}
		log.Printf("password:%s passwordHash: %s", data.Password, passwordHash)

		user, err := storage.Create(ctx.Context(), data.Email, passwordHash, data.Name)
		if err != nil {
			if errors.Is(err, AlreadyExists) {
				return fiber.NewError(fiber.StatusBadRequest, err.Error())
			}
			return fiber.ErrInternalServerError
		}

		token, err := GetToken(user, config.JwtSecret)
		if err != nil {
			return err
		}

		return ctx.JSON(RegistrationResponse{Token: token})
	}
}

func Login(config *config.AppConfig, storage Storage, validator *utils.AppValidator) fiber.Handler {
	type LoginRequest struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,gte=4,lte=255"`
	}

	type LoginResponse struct {
		Token string `json:"token"`
	}

	return func(ctx fiber.Ctx) error {
		data := LoginRequest{}
		err := ctx.Bind().Body(&data)
		if err != nil {
			return err
		}
		err = validator.Validate(data)
		if err != nil {
			return err
		}

		user, err := storage.GetUserByEmail(ctx.Context(), data.Email)
		if err != nil {
			if errors.Is(err, NotFound) {
				return fiber.NewError(fiber.StatusBadRequest, "wrong email or password")
			}
			return err
		}

		isValid, err := utils.ComparePassword(data.Password, user.passwordHash)
		if err != nil {
			return err
		}
		if !isValid {
			return fiber.NewError(fiber.StatusBadRequest, "wrong email or password")
		}

		token, err := GetToken(user, config.JwtSecret)
		if err != nil {
			return err
		}

		return ctx.JSON(LoginResponse{Token: token})
	}
}
