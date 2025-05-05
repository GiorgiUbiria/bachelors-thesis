package handlers

import (
	"time"

	"github.com/GiorgiUbiria/bachelor/config"
	"github.com/GiorgiUbiria/bachelor/models"
	"github.com/gofiber/fiber/v3"
)

// GetActivityAnalytics retrieves user activity analytics
func GetActivityAnalytics(c fiber.Ctx) error {
	var activities []models.UserActivity
	timeRange := c.Query("timeRange", "24h")

	// Calculate time range
	var startTime time.Time
	switch timeRange {
	case "24h":
		startTime = time.Now().Add(-24 * time.Hour)
	case "7d":
		startTime = time.Now().Add(-7 * 24 * time.Hour)
	case "30d":
		startTime = time.Now().Add(-30 * 24 * time.Hour)
	default:
		startTime = time.Now().Add(-24 * time.Hour)
	}

	// Get activities within time range
	result := config.DB.Where("created_at >= ?", startTime).Find(&activities)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch activity analytics",
		})
	}

	// Group activities by type
	activityCounts := make(map[string]int)
	for _, activity := range activities {
		activityCounts[activity.Type]++
	}

	return c.JSON(fiber.Map{
		"time_range": timeRange,
		"activities": activityCounts,
		"total":      len(activities),
	})
}

// GetRequestAnalytics retrieves request log analytics
func GetRequestAnalytics(c fiber.Ctx) error {
	var logs []models.RequestLog
	timeRange := c.Query("timeRange", "24h")

	// Calculate time range
	var startTime time.Time
	switch timeRange {
	case "24h":
		startTime = time.Now().Add(-24 * time.Hour)
	case "7d":
		startTime = time.Now().Add(-7 * 24 * time.Hour)
	case "30d":
		startTime = time.Now().Add(-30 * 24 * time.Hour)
	default:
		startTime = time.Now().Add(-24 * time.Hour)
	}

	// Get logs within time range
	result := config.DB.Where("created_at >= ?", startTime).Find(&logs)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch request analytics",
		})
	}

	// Group logs by category
	categoryCounts := make(map[string]int)
	for _, log := range logs {
		categoryCounts[log.Category]++
	}

	return c.JSON(fiber.Map{
		"time_range": timeRange,
		"categories": categoryCounts,
		"total":      len(logs),
	})
}

// GetPopularProducts retrieves most popular products
func GetPopularProducts(c fiber.Ctx) error {
	var activities []models.UserActivity
	timeRange := c.Query("timeRange", "24h")

	// Calculate time range
	var startTime time.Time
	switch timeRange {
	case "24h":
		startTime = time.Now().Add(-24 * time.Hour)
	case "7d":
		startTime = time.Now().Add(-7 * 24 * time.Hour)
	case "30d":
		startTime = time.Now().Add(-30 * 24 * time.Hour)
	default:
		startTime = time.Now().Add(-24 * time.Hour)
	}

	// Get product views within time range
	result := config.DB.Where("created_at >= ? AND type = ?", startTime, "view").Find(&activities)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch popular products",
		})
	}

	// Count product views
	productViews := make(map[uint]int)
	for _, activity := range activities {
		if activity.ProductID != nil {
			productViews[*activity.ProductID]++
		}
	}

	// Get top 10 products
	var topProducts []struct {
		ProductID uint
		Views     int
	}
	for productID, views := range productViews {
		topProducts = append(topProducts, struct {
			ProductID uint
			Views     int
		}{productID, views})
	}

	// Sort by views (descending)
	for i := 0; i < len(topProducts); i++ {
		for j := i + 1; j < len(topProducts); j++ {
			if topProducts[i].Views < topProducts[j].Views {
				topProducts[i], topProducts[j] = topProducts[j], topProducts[i]
			}
		}
	}

	// Limit to top 10
	if len(topProducts) > 10 {
		topProducts = topProducts[:10]
	}

	return c.JSON(fiber.Map{
		"time_range": timeRange,
		"products":   topProducts,
	})
}

// GetActiveUsers retrieves most active users
func GetActiveUsers(c fiber.Ctx) error {
	var activities []models.UserActivity
	timeRange := c.Query("timeRange", "24h")

	// Calculate time range
	var startTime time.Time
	switch timeRange {
	case "24h":
		startTime = time.Now().Add(-24 * time.Hour)
	case "7d":
		startTime = time.Now().Add(-7 * 24 * time.Hour)
	case "30d":
		startTime = time.Now().Add(-30 * 24 * time.Hour)
	default:
		startTime = time.Now().Add(-24 * time.Hour)
	}

	// Get activities within time range
	result := config.DB.Where("created_at >= ?", startTime).Find(&activities)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch active users",
		})
	}

	// Count user activities
	userActivities := make(map[uint]int)
	for _, activity := range activities {
		userActivities[activity.UserID]++
	}

	// Get top 10 users
	var topUsers []struct {
		UserID     uint
		Activities int
	}
	for userID, activities := range userActivities {
		topUsers = append(topUsers, struct {
			UserID     uint
			Activities int
		}{userID, activities})
	}

	// Sort by activities (descending)
	for i := 0; i < len(topUsers); i++ {
		for j := i + 1; j < len(topUsers); j++ {
			if topUsers[i].Activities < topUsers[j].Activities {
				topUsers[i], topUsers[j] = topUsers[j], topUsers[i]
			}
		}
	}

	// Limit to top 10
	if len(topUsers) > 10 {
		topUsers = topUsers[:10]
	}

	return c.JSON(fiber.Map{
		"time_range": timeRange,
		"users":      topUsers,
	})
}
