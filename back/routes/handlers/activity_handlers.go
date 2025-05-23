package handlers

import (
	"strconv"
	"time"

	"github.com/GiorgiUbiria/bachelor/config"
	"github.com/GiorgiUbiria/bachelor/models"
	"github.com/gofiber/fiber/v3"
)

// TrackProductView tracks when a user views a product
func TrackProductView(c fiber.Ctx) error {
	var requestData struct {
		UserID    uint   `json:"user_id"`
		ProductID uint   `json:"product_id"`
		SessionID string `json:"session_id"`
	}

	if err := c.Bind().JSON(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate that product exists
	var product models.Product
	if err := config.DB.First(&product, requestData.ProductID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	// Create activity record
	activity := models.UserActivity{
		UserID:    requestData.UserID,
		Type:      "view",
		ProductID: &requestData.ProductID,
		Details:   "Product viewed: " + product.Name,
		SessionID: requestData.SessionID,
	}

	if err := config.DB.Create(&activity).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to track product view",
		})
	}

	return c.JSON(fiber.Map{
		"message":     "Product view tracked successfully",
		"activity_id": activity.ID,
	})
}

// TrackClick tracks when a user clicks on something
func TrackClick(c fiber.Ctx) error {
	var requestData struct {
		UserID    uint   `json:"user_id"`
		ProductID *uint  `json:"product_id,omitempty"`
		SessionID string `json:"session_id"`
		Element   string `json:"element"` // e.g., "add_to_cart", "buy_now", "product_link"
		Details   string `json:"details,omitempty"`
	}

	if err := c.Bind().JSON(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Create activity record
	activity := models.UserActivity{
		UserID:    requestData.UserID,
		Type:      "click",
		ProductID: requestData.ProductID,
		Details:   "Clicked: " + requestData.Element + " - " + requestData.Details,
		SessionID: requestData.SessionID,
	}

	if err := config.DB.Create(&activity).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to track click",
		})
	}

	return c.JSON(fiber.Map{
		"message":     "Click tracked successfully",
		"activity_id": activity.ID,
	})
}

// TrackSearch tracks when a user performs a search
func TrackSearch(c fiber.Ctx) error {
	var requestData struct {
		UserID    uint   `json:"user_id"`
		SessionID string `json:"session_id"`
		Query     string `json:"query"`
		Results   int    `json:"results"` // Number of results returned
	}

	if err := c.Bind().JSON(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if requestData.Query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Search query is required",
		})
	}

	// Create activity record
	activity := models.UserActivity{
		UserID:    requestData.UserID,
		Type:      "search",
		Details:   "Search query: " + requestData.Query + " (Results: " + strconv.Itoa(requestData.Results) + ")",
		SessionID: requestData.SessionID,
	}

	if err := config.DB.Create(&activity).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to track search",
		})
	}

	return c.JSON(fiber.Map{
		"message":     "Search tracked successfully",
		"activity_id": activity.ID,
	})
}

// StartSession creates or updates a user session
func StartSession(c fiber.Ctx) error {
	var requestData struct {
		UserID    uint   `json:"user_id"`
		SessionID string `json:"session_id"`
		UserAgent string `json:"user_agent,omitempty"`
		IPAddress string `json:"ip_address,omitempty"`
	}

	if err := c.Bind().JSON(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if requestData.SessionID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Session ID is required",
		})
	}

	// Create session start activity
	activity := models.UserActivity{
		UserID:    requestData.UserID,
		Type:      "session_start",
		Details:   "Session started - IP: " + requestData.IPAddress + " - UA: " + requestData.UserAgent,
		SessionID: requestData.SessionID,
	}

	if err := config.DB.Create(&activity).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to track session start",
		})
	}

	return c.JSON(fiber.Map{
		"message":     "Session started successfully",
		"session_id":  requestData.SessionID,
		"activity_id": activity.ID,
	})
}

// EndSession ends a user session
func EndSession(c fiber.Ctx) error {
	var requestData struct {
		UserID    uint   `json:"user_id"`
		SessionID string `json:"session_id"`
		Duration  int    `json:"duration"` // Session duration in seconds
	}

	if err := c.Bind().JSON(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Create session end activity
	activity := models.UserActivity{
		UserID:    requestData.UserID,
		Type:      "session_end",
		Details:   "Session ended - Duration: " + strconv.Itoa(requestData.Duration) + " seconds",
		SessionID: requestData.SessionID,
	}

	if err := config.DB.Create(&activity).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to track session end",
		})
	}

	return c.JSON(fiber.Map{
		"message":     "Session ended successfully",
		"activity_id": activity.ID,
	})
}

// GetUserActivityTimeline retrieves a user's activity timeline
func GetUserActivityTimeline(c fiber.Ctx) error {
	userID := c.Params("userId")

	// Pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset := (page - 1) * limit

	// Optional filters
	activityType := c.Query("type")
	sessionID := c.Query("session_id")
	since := c.Query("since") // ISO date string

	var activities []models.UserActivity
	query := config.DB.Where("user_id = ?", userID).Preload("Product")

	// Apply filters
	if activityType != "" {
		query = query.Where("type = ?", activityType)
	}

	if sessionID != "" {
		query = query.Where("session_id = ?", sessionID)
	}

	if since != "" {
		if sinceTime, err := time.Parse(time.RFC3339, since); err == nil {
			query = query.Where("created_at >= ?", sinceTime)
		}
	}

	// Get total count
	var total int64
	query.Model(&models.UserActivity{}).Count(&total)

	// Apply pagination and fetch results
	result := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&activities)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user activity timeline",
		})
	}

	return c.JSON(fiber.Map{
		"activities": activities,
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// GetSessionActivities retrieves all activities for a specific session
func GetSessionActivities(c fiber.Ctx) error {
	sessionID := c.Params("sessionId")

	var activities []models.UserActivity
	result := config.DB.Where("session_id = ?", sessionID).
		Preload("Product").
		Order("created_at ASC").
		Find(&activities)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch session activities",
		})
	}

	// Calculate session statistics
	var sessionStats fiber.Map
	if len(activities) > 0 {
		startTime := activities[0].CreatedAt
		endTime := activities[len(activities)-1].CreatedAt
		duration := endTime.Sub(startTime)

		// Count activity types
		typeCounts := make(map[string]int)
		for _, activity := range activities {
			typeCounts[activity.Type]++
		}

		sessionStats = fiber.Map{
			"start_time":       startTime,
			"end_time":         endTime,
			"duration":         duration.Seconds(),
			"total_activities": len(activities),
			"activity_counts":  typeCounts,
		}
	}

	return c.JSON(fiber.Map{
		"session_id": sessionID,
		"activities": activities,
		"statistics": sessionStats,
	})
}

// TrackAddToCart tracks when a user adds a product to cart
func TrackAddToCart(c fiber.Ctx) error {
	var requestData struct {
		UserID    uint   `json:"user_id"`
		ProductID uint   `json:"product_id"`
		Quantity  int    `json:"quantity"`
		SessionID string `json:"session_id"`
	}

	if err := c.Bind().JSON(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Create activity record
	activity := models.UserActivity{
		UserID:    requestData.UserID,
		Type:      "add_to_cart",
		ProductID: &requestData.ProductID,
		Details:   "Added to cart - Quantity: " + strconv.Itoa(requestData.Quantity),
		SessionID: requestData.SessionID,
	}

	if err := config.DB.Create(&activity).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to track add to cart",
		})
	}

	return c.JSON(fiber.Map{
		"message":     "Add to cart tracked successfully",
		"activity_id": activity.ID,
	})
}

// TrackRemoveFromCart tracks when a user removes a product from cart
func TrackRemoveFromCart(c fiber.Ctx) error {
	var requestData struct {
		UserID    uint   `json:"user_id"`
		ProductID uint   `json:"product_id"`
		SessionID string `json:"session_id"`
	}

	if err := c.Bind().JSON(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Create activity record
	activity := models.UserActivity{
		UserID:    requestData.UserID,
		Type:      "remove_from_cart",
		ProductID: &requestData.ProductID,
		Details:   "Removed from cart",
		SessionID: requestData.SessionID,
	}

	if err := config.DB.Create(&activity).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to track remove from cart",
		})
	}

	return c.JSON(fiber.Map{
		"message":     "Remove from cart tracked successfully",
		"activity_id": activity.ID,
	})
}

// TrackAddToFavorites tracks when a user adds a product to favorites
func TrackAddToFavorites(c fiber.Ctx) error {
	var requestData struct {
		UserID    uint   `json:"user_id"`
		ProductID uint   `json:"product_id"`
		SessionID string `json:"session_id"`
	}

	if err := c.Bind().JSON(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Create activity record
	activity := models.UserActivity{
		UserID:    requestData.UserID,
		Type:      "add_to_favorites",
		ProductID: &requestData.ProductID,
		Details:   "Added to favorites",
		SessionID: requestData.SessionID,
	}

	if err := config.DB.Create(&activity).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to track add to favorites",
		})
	}

	return c.JSON(fiber.Map{
		"message":     "Add to favorites tracked successfully",
		"activity_id": activity.ID,
	})
}
