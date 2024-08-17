package middleware

import (
	"VK-test/application/tokens"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Protected protect routes
func Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwt.SigningMethodHS256.Name,
			Key:    tokens.JwtKey},
		ErrorHandler: jwtError,
		TokenLookup:  "header:Authorization",
		AuthScheme:   "Bearer",
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	return c.SendStatus(fiber.StatusUnauthorized)
}
