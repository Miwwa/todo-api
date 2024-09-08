package utils

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v3"
)

func Json404(ctx fiber.Ctx) error {
	ctx.Status(404)
	return ctx.JSON(fiber.Map{
		"status":  fiber.StatusNotFound,
		"message": "404 Not Found",
	})
}

func JsonErrorHandler(ctx fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError

	// Retrieve the custom status code if it's a *fiber.Error
	var e *fiber.Error
	if errors.As(err, &e) {
		code = e.Code
	}

	// Send error as json
	err = ctx.Status(code).JSON(fiber.Map{
		"status":  code,
		"message": err.Error(),
	})
	// Send error in plaintext as fallback
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Internal Server Error: %s", err.Error()))
	}

	return nil
}
