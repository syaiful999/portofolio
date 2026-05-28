package middleware

import (
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/contrib/jwt"
)

var (
	d struct{}
)

func Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:    jwtware.SigningKey{Key: []byte(os.Getenv("SERVICE_JWT_SECRET"))},
		ErrorHandler: jwtError,
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"IsSuccess": false,
			"Status":    http.StatusBadRequest,
			"Message":   err.Error(),
			"Data":      d,
		})
	}

	return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
		"IsSuccess": false,
		"Status":    http.StatusUnauthorized,
		"Message":   "Invalid or expired JWT",
		"Data":      d,
	})
}
