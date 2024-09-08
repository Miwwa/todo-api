package users

import (
	"errors"
	"github.com/gofiber/fiber/v3"
	"log"
	"todo-api/config"
	"todo-api/utils"
)

func SetupRoutes(app *fiber.App, config *config.AppConfig, storage UsersStorage) {
	app.Post("/register", Register(config, storage))
	app.Post("/login", Login(config, storage))
}

func Register(config *config.AppConfig, storage UsersStorage) fiber.Handler {
	type RegistrationRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
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

		passwordHash, err := utils.HashPassword(data.Password)
		if err != nil {
			return err
		}
		log.Printf("password:%s passwordHash: %s", data.Password, passwordHash)

		user, err := storage.Create(ctx.Context(), data.Email, passwordHash, data.Name)
		if err != nil {
			if errors.Is(err, UserAlreadyExists) {
				return fiber.NewError(fiber.StatusBadRequest, err.Error())
			}
			return fiber.ErrInternalServerError
		}

		token, err := GetToken(user, config.JwtSecret())
		if err != nil {
			return err
		}

		return ctx.JSON(RegistrationResponse{Token: token})
	}
}

func Login(config *config.AppConfig, storage UsersStorage) fiber.Handler {
	type LoginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
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

		user, err := storage.GetUserByEmail(ctx.Context(), data.Email)
		if err != nil {
			if errors.Is(err, UserNotFound) {
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

		token, err := GetToken(user, config.JwtSecret())
		if err != nil {
			return err
		}

		return ctx.JSON(LoginResponse{Token: token})
	}
}
