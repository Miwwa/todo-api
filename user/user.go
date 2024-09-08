package user

import "github.com/gofiber/fiber/v3"

type Id string

type User struct {
	Id           Id     `json:"id"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	passwordHash string
}

func FromContext(c fiber.Ctx) User {
	return c.Locals(userContextKey).(User)
}
