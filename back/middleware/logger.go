package middleware

import (
	"time"

	"github.com/gofiber/fiber/v3"
)

func RequestTimer() fiber.Handler {
	return func(c fiber.Ctx) error {
		c.Locals("startTime", time.Now())
		return c.Next()
	}
}
