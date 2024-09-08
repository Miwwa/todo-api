package users

import (
	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	jwtware "jwt"
	"todo-api/utils"
)

const (
	ContextKey = "user"
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
		SigningKey:   jwtware.SigningKey{Key: []byte(secretKey)},
		ErrorHandler: utils.AuthErrorHandler,
		ContextKey:   ContextKey,
	})
}
