package middleware

import (
	"fmt"
	"moyo-gateway-service/utils"

	"github.com/gofiber/fiber/v2"
)

func CustomLogger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Read request body
		var requestBody []byte
		if c.Body() != nil {
			requestBody = c.Body()
		}
		// Log request details
		message := fmt.Sprintf("[Request] %s %s %s", c.Method(), c.OriginalURL(), string(requestBody))
		utils.PushLogf("", message, "")

		// Continue to next middleware
		err := c.Next()
		if err != nil {
			// Log the error
			utils.PushLogf("", "[Error]", err.Error())
			return err
		}
		return nil
	}
}
