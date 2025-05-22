package middleware

import (
	"time"

	"github.com/GiorgiUbiria/bachelor/routes/handlers"
	"github.com/gofiber/fiber/v3"
)

func RequestTimer() fiber.Handler {
	return func(c fiber.Ctx) error {
		c.Locals("startTime", time.Now())
		return c.Next()
	}
}

func RequestLoggerAndAnomaly() fiber.Handler {
	return func(c fiber.Ctx) error {
		// Skip static files and health check
		if c.Path() == "/health" || c.Path() == "/favicon.ico" || c.Path() == "/robots.txt" {
			return c.Next()
		}
		// Log and check anomaly
		handlers.LogRequestHandler(c)
		return c.Next()
	}
}
