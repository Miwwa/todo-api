package users

import (
	"errors"
	"github.com/gofiber/fiber/v3"
	"log"
	"todo-api/utils"
)

func SetupRoutes(app *fiber.App, storage UsersStorage) {
	app.Post("/register", Register(storage))
	app.Post("/login", Login(storage))
}

func Register(storage UsersStorage) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		data := RegistrationData{}
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

		err = ctx.JSON(user)
		if err != nil {
			return err
		}

		return nil
	}
}

func Login(storage UsersStorage) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		data := LoginData{}
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

		err = ctx.JSON(user)
		if err != nil {
			return err
		}

		return nil
	}
}
