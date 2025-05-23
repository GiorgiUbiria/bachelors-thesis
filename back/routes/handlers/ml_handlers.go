package handlers

import (
	"strconv"

	"github.com/GiorgiUbiria/bachelor/config"
	"github.com/GiorgiUbiria/bachelor/models"
	"github.com/GiorgiUbiria/bachelor/services"
	"github.com/gofiber/fiber/v3"
)

// GetMLTrainingData provides formatted training data for ML models
func GetMLTrainingData(c fiber.Ctx) error {
	data, err := services.GetMLTrainingData()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to fetch ML training data",
			"details": err.Error(),
		})
	}

	return c.JSON(data)
}

// TriggerMLRetraining triggers retraining of ML models
func TriggerMLRetraining(c fiber.Ctx) error {
	err := services.TriggerMLRetraining()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to trigger ML retraining",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "ML retraining triggered successfully",
		"status":  "success",
	})
}

// GetUserRecommendations gets product recommendations for a specific user
func GetUserRecommendations(c fiber.Ctx) error {
	userIDStr := c.Params("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	recommendations, err := services.GetRecommendations(uint(userID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get recommendations",
			"details": err.Error(),
		})
	}

	// Get product details for recommendations
	var products []models.Product
	if len(recommendations) > 0 {
		// Convert uint slice to interface slice for GORM
		ids := make([]interface{}, len(recommendations))
		for i, id := range recommendations {
			ids[i] = id
		}

		// Query products with the recommended IDs
		result := config.DB.Where("id IN ?", ids).Find(&products)
		if result.Error != nil {
			// If we can't get product details, just return the IDs
			return c.JSON(fiber.Map{
				"user_id":         userID,
				"recommendations": recommendations,
				"products":        []models.Product{},
			})
		}
	}

	return c.JSON(fiber.Map{
		"user_id":         userID,
		"recommendations": recommendations,
		"products":        products,
	})
}

// GetUserBehaviorAnalysis analyzes user behavior patterns
func GetUserBehaviorAnalysis(c fiber.Ctx) error {
	userIDStr := c.Params("userId")
	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	behavior, err := services.AnalyzeUserBehavior(uint(userID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to analyze user behavior",
			"details": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"user_id":        userID,
		"cluster":        behavior.Cluster,
		"cluster_center": behavior.ClusterCenter,
		"analysis":       getClusterDescription(behavior.Cluster),
	})
}

// GetTrendAnalysis provides sales trend analysis data
func GetTrendAnalysis(c fiber.Ctx) error {
	// Get query parameters
	category := c.Query("category")
	timeRange := c.Query("timeRange", "30d")

	data, err := services.GetMLTrainingData()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to fetch trend data",
			"details": err.Error(),
		})
	}

	// Filter by category if specified
	trends := data.SalesTrends
	if category != "" {
		filteredTrends := make([]services.SalesTrendFeatures, 0)
		for _, trend := range trends {
			if trend.Category == category {
				filteredTrends = append(filteredTrends, trend)
			}
		}
		trends = filteredTrends
	}

	// Aggregate trends by date
	trendMap := make(map[string]float64)
	for _, trend := range trends {
		dateStr := trend.Date.Format("2006-01-02")
		trendMap[dateStr] += trend.Revenue
	}

	return c.JSON(fiber.Map{
		"category":      category,
		"time_range":    timeRange,
		"trends":        trendMap,
		"total_records": len(trends),
	})
}

// CheckMLServiceHealth checks the health of the ML service
func CheckMLServiceHealth(c fiber.Ctx) error {
	err := services.CheckMLServiceHealth()
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"status": "unhealthy",
			"error":  err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "healthy",
		"message": "ML service is operational",
	})
}

// Helper function to provide cluster descriptions
func getClusterDescription(cluster int) string {
	descriptions := map[int]string{
		0: "Casual browsers - low engagement, occasional purchases",
		1: "Active shoppers - high engagement, frequent purchases",
		2: "Window shoppers - high browsing, low conversion",
		3: "Loyal customers - consistent purchasing patterns",
		4: "New users - limited activity, potential for growth",
	}

	if desc, exists := descriptions[cluster]; exists {
		return desc
	}
	return "Unknown cluster pattern"
}
