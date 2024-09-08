package user

import (
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	jwtware "jwt"
)

const (
	tokenContextKey = "jwt"
	ContextKey      = "user"
)

var (
	SigningMethod = jwt.SigningMethodHS512
)

func GetToken(user User, secretKey string) (string, error) {
	token := jwt.New(SigningMethod)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = user.Id
	claims["email"] = user.Email
	claims["name"] = user.Name

	return token.SignedString([]byte(secretKey))
}

func ValidateAndExtractTokenMiddleware(secretKey string) func(fiber.Ctx) error {
	return jwtware.New(jwtware.Config{
		SigningKey:     jwtware.SigningKey{Key: []byte(secretKey)},
		SuccessHandler: successHandler,
		ErrorHandler:   authErrorHandler,
		ContextKey:     tokenContextKey,
	})
}

func successHandler(c fiber.Ctx) error {
	token, ok := c.Locals(tokenContextKey).(*jwt.Token)
	if !ok {
		return fiber.ErrUnauthorized
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fiber.ErrInternalServerError
	}

	user := User{
		Id:    claims["sub"].(Id),
		Email: claims["email"].(string),
		Name:  claims["name"].(string),
	}
	// todo: validate user

	c.Locals(ContextKey, &user)

	return c.Next()
}

func authErrorHandler(c fiber.Ctx, err error) error {
	if err.Error() == "missing or malformed JWT" {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusBadRequest,
			"message": "missing or malformed JWT",
		})
	} else {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"status":  fiber.StatusUnauthorized,
			"message": "invalid or expired auth token",
		})
	}
}
