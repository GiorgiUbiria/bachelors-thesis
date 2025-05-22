package handlers

import (
	"net"
	"time"

	"github.com/GiorgiUbiria/bachelor/config"
	"github.com/GiorgiUbiria/bachelor/models"
	"github.com/GiorgiUbiria/bachelor/services"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm/clause"
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

// Helper to check if IP is banned
func isIPBanned(ip string) (bool, time.Time) {
	var ban models.BannedIP
	result := config.DB.Where("ip = ? AND banned_until > ?", ip, time.Now()).First(&ban)
	if result.Error == nil {
		return true, ban.BannedUntil
	}
	return false, time.Time{}
}

// LogRequestHandler logs a request and runs anomaly detection
func LogRequestHandler(c fiber.Ctx) error {
	ip := c.IP()
	if banned, until := isIPBanned(ip); banned {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":        "IP banned",
			"banned_until": until,
		})
	}

	features := extractRequestFeatures(c)
	prediction, err := services.DetectAnomaly([][]float64{features})
	category := "normal"
	if err == nil && prediction == -1 {
		category = "anomaly"
		// Preemptive action: ban IP for 1 hour
		ban := models.BannedIP{
			IP:          ip,
			BannedUntil: time.Now().Add(1 * time.Hour),
			Reason:      "Anomaly detected",
			CreatedAt:   time.Now(),
		}
		config.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&ban)
		// TODO: Notify admin (placeholder)
	}

	// Save to DB
	log := models.RequestLog{
		IP:           ip,
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
