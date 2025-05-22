package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/GiorgiUbiria/bachelor/config"
	"github.com/GiorgiUbiria/bachelor/models"
	"github.com/GiorgiUbiria/bachelor/services"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
)

// Helper to set up a fresh Fiber app and DB for each test
func setupTestApp() *fiber.App {
	os.Setenv("ENV", "test")
	config.InitDB()
	app := fiber.New()
	app.Use(func(c fiber.Ctx) error {
		c.Locals("startTime", time.Now())
		return c.Next()
	})
	app.Post("/log-request", LogRequestHandler)
	app.Get("/analytics/requests/recent", GetRecentRequestLogs)
	return app
}

func TestNormalRequestLogging(t *testing.T) {
	app := setupTestApp()
	// Simulate a normal request
	req := httptest.NewRequest(http.MethodPost, "/log-request", nil)
	req.Header.Set("User-Agent", "test-agent")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Check log in DB
	var logs []models.RequestLog
	config.DB.Find(&logs)
	assert.NotEmpty(t, logs)
	assert.Equal(t, "normal", logs[len(logs)-1].Category)
}

func TestAnomalyDetectionAndBan(t *testing.T) {
	app := setupTestApp()
	// Patch services.DetectAnomaly to always return -1 (anomaly)
	origDetect := services.DetectAnomaly
	services.DetectAnomaly = func(features [][]float64) (int, error) { return -1, nil }
	defer func() { services.DetectAnomaly = origDetect }()

	req := httptest.NewRequest(http.MethodPost, "/log-request", nil)
	req.Header.Set("X-Forwarded-For", "1.2.3.4")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Check log in DB
	var log models.RequestLog
	config.DB.Last(&log)
	assert.Equal(t, "anomaly", log.Category)

	// Check ban in DB
	var ban models.BannedIP
	config.DB.Where("ip = ?", "1.2.3.4").First(&ban)
	assert.NotZero(t, ban.ID)
	assert.True(t, ban.BannedUntil.After(time.Now()))
}

func TestBannedIPCannotRequest(t *testing.T) {
	app := setupTestApp()
	// Insert ban manually
	ban := models.BannedIP{
		IP:          "5.6.7.8",
		BannedUntil: time.Now().Add(1 * time.Hour),
		Reason:      "test ban",
		CreatedAt:   time.Now(),
	}
	config.DB.Create(&ban)

	req := httptest.NewRequest(http.MethodPost, "/log-request", nil)
	req.Header.Set("X-Forwarded-For", "5.6.7.8")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
}

func TestAnalyticsRecentRequests(t *testing.T) {
	app := setupTestApp()
	// Add a log
	log := models.RequestLog{
		IP:           "9.9.9.9",
		Method:       "POST",
		Path:         "/log-request",
		Status:       200,
		UserAgent:    "test-agent",
		Category:     "normal",
		Details:      "test",
		ResponseTime: 10,
		CreatedAt:    time.Now(),
	}
	config.DB.Create(&log)

	req := httptest.NewRequest(http.MethodGet, "/analytics/requests/recent?limit=1", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var logs []models.RequestLog
	err = json.NewDecoder(resp.Body).Decode(&logs)
	assert.NoError(t, err)
	assert.Len(t, logs, 1)
	assert.Equal(t, "9.9.9.9", logs[0].IP)
}
