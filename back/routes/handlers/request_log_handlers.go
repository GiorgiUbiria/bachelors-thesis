package handlers

import (
	"net"
	"time"

	"github.com/GiorgiUbiria/bachelor/config"
	"github.com/GiorgiUbiria/bachelor/models"
	"github.com/GiorgiUbiria/bachelor/services"
	"github.com/gofiber/fiber/v3"
)

// Extract features from the request for anomaly detection
func extractRequestFeatures(c fiber.Ctx) []float64 {
	// Example: IP (as int), method (as int), path length, status, response time
	ip := ipToFloat64(c.IP())
	method := methodToFloat64(c.Method())
	pathLen := float64(len(c.Path()))
	status := float64(c.Response().StatusCode())
	responseTime := float64(time.Since(c.Locals("startTime").(time.Time)).Milliseconds())
	return []float64{ip, method, pathLen, status, responseTime}
}

func ipToFloat64(ip string) float64 {
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return 0
	}
	v4 := parsed.To4()
	if v4 == nil {
		return 0
	}
	return float64(uint32(v4[0])<<24 | uint32(v4[1])<<16 | uint32(v4[2])<<8 | uint32(v4[3]))
}

func methodToFloat64(method string) float64 {
	switch method {
	case fiber.MethodGet:
		return 1
	case fiber.MethodPost:
		return 2
	case fiber.MethodPut:
		return 3
	case fiber.MethodDelete:
		return 4
	case fiber.MethodPatch:
		return 5
	default:
		return 0
	}
}

// LogRequestHandler logs a request and runs anomaly detection
func LogRequestHandler(c fiber.Ctx) error {
	features := extractRequestFeatures(c)
	prediction, err := services.DetectAnomaly([][]float64{features})
	category := "normal"
	if err == nil && prediction == -1 {
		category = "anomaly"
		// Optionally: trigger auto-action, e.g., block IP
	}

	// Save to DB
	log := models.RequestLog{
		IP:           c.IP(),
		Method:       c.Method(),
		Path:         c.Path(),
		Status:       c.Response().StatusCode(),
		UserAgent:    string(c.Request().Header.UserAgent()),
		Category:     category,
		Details:      "Logged via API",
		ResponseTime: float64(time.Since(c.Locals("startTime").(time.Time)).Milliseconds()),
	}
	config.DB.Create(&log)

	return c.JSON(fiber.Map{
		"category": category,
		"log_id":   log.ID,
	})
}
