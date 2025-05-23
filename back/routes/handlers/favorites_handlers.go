package handlers

import (
	"strconv"

	"github.com/GiorgiUbiria/bachelor/config"
	"github.com/GiorgiUbiria/bachelor/models"
	"github.com/gofiber/fiber/v3"
)

// AddToFavorites adds a product to user's favorites
func AddToFavorites(c fiber.Ctx) error {
	var requestData struct {
		UserID    uint `json:"user_id"`
		ProductID uint `json:"product_id"`
	}

	if err := c.Bind().JSON(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Check if product exists
	var product models.Product
	if err := config.DB.First(&product, requestData.ProductID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	// Check if already in favorites
	var existingFavorite models.Favorite
	result := config.DB.Where("user_id = ? AND product_id = ?", requestData.UserID, requestData.ProductID).First(&existingFavorite)
	if result.Error == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Product already in favorites",
		})
	}

	// Create favorite
	favorite := models.Favorite{
		UserID:    requestData.UserID,
		ProductID: requestData.ProductID,
	}

	if err := config.DB.Create(&favorite).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add product to favorites",
		})
	}

	// Load the favorite with product details
	config.DB.Preload("Product").First(&favorite, favorite.ID)

	return c.Status(fiber.StatusCreated).JSON(favorite)
}

// RemoveFromFavorites removes a product from user's favorites
func RemoveFromFavorites(c fiber.Ctx) error {
	id := c.Params("id")

	// Find and delete the favorite
	var favorite models.Favorite
	result := config.DB.First(&favorite, id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Favorite not found",
		})
	}

	// Delete the favorite
	if err := config.DB.Delete(&favorite).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to remove from favorites",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// GetFavorites retrieves all favorites for a user with pagination
func GetFavorites(c fiber.Ctx) error {
	userID := c.Params("userId")

	// Pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset := (page - 1) * limit

	// Optional filters
	category := c.Query("category")
	sortBy := c.Query("sort_by", "created_at")
	sortOrder := c.Query("sort_order", "desc")

	var favorites []models.Favorite
	query := config.DB.Where("user_id = ?", userID).Preload("Product")

	// Apply category filter if specified
	if category != "" {
		query = query.Joins("JOIN products ON favorites.product_id = products.id").
			Where("products.category = ?", category)
	}

	// Apply sorting
	orderClause := sortBy
	if sortOrder == "desc" {
		orderClause += " DESC"
	}
	query = query.Order(orderClause)

	// Get total count
	var total int64
	query.Model(&models.Favorite{}).Count(&total)

	// Apply pagination and fetch results
	result := query.Offset(offset).Limit(limit).Find(&favorites)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user favorites",
		})
	}

	return c.JSON(fiber.Map{
		"favorites": favorites,
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// BulkAddToFavorites adds multiple products to favorites
func BulkAddToFavorites(c fiber.Ctx) error {
	var requestData struct {
		UserID     uint   `json:"user_id"`
		ProductIDs []uint `json:"product_ids"`
	}

	if err := c.Bind().JSON(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if len(requestData.ProductIDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No product IDs provided",
		})
	}

	// Validate all products exist
	var productCount int64
	config.DB.Model(&models.Product{}).Where("id IN ?", requestData.ProductIDs).Count(&productCount)
	if int(productCount) != len(requestData.ProductIDs) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "One or more products not found",
		})
	}

	// Get existing favorites to avoid duplicates
	var existingFavorites []models.Favorite
	config.DB.Where("user_id = ? AND product_id IN ?", requestData.UserID, requestData.ProductIDs).Find(&existingFavorites)

	existingProductIDs := make(map[uint]bool)
	for _, fav := range existingFavorites {
		existingProductIDs[fav.ProductID] = true
	}

	// Create new favorites
	var newFavorites []models.Favorite
	var skippedProducts []uint

	for _, productID := range requestData.ProductIDs {
		if existingProductIDs[productID] {
			skippedProducts = append(skippedProducts, productID)
			continue
		}

		newFavorites = append(newFavorites, models.Favorite{
			UserID:    requestData.UserID,
			ProductID: productID,
		})
	}

	// Bulk insert new favorites
	if len(newFavorites) > 0 {
		if err := config.DB.Create(&newFavorites).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to add products to favorites",
			})
		}
	}

	return c.JSON(fiber.Map{
		"message":          "Bulk add to favorites completed",
		"added_count":      len(newFavorites),
		"skipped_count":    len(skippedProducts),
		"skipped_products": skippedProducts,
	})
}

// BulkRemoveFromFavorites removes multiple products from favorites
func BulkRemoveFromFavorites(c fiber.Ctx) error {
	var requestData struct {
		UserID     uint   `json:"user_id"`
		ProductIDs []uint `json:"product_ids"`
	}

	if err := c.Bind().JSON(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if len(requestData.ProductIDs) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No product IDs provided",
		})
	}

	// Remove favorites
	result := config.DB.Where("user_id = ? AND product_id IN ?", requestData.UserID, requestData.ProductIDs).Delete(&models.Favorite{})
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to remove products from favorites",
		})
	}

	return c.JSON(fiber.Map{
		"message":       "Bulk remove from favorites completed",
		"removed_count": result.RowsAffected,
	})
}

// CheckFavoriteStatus checks if a product is in user's favorites
func CheckFavoriteStatus(c fiber.Ctx) error {
	userID := c.Params("userId")
	productID := c.Params("productId")

	var favorite models.Favorite
	result := config.DB.Where("user_id = ? AND product_id = ?", userID, productID).First(&favorite)

	isFavorite := result.Error == nil

	return c.JSON(fiber.Map{
		"user_id":     userID,
		"product_id":  productID,
		"is_favorite": isFavorite,
		"favorite_id": func() *uint {
			if isFavorite {
				return &favorite.ID
			}
			return nil
		}(),
	})
}

// GetFavoritesByCategory gets user's favorites grouped by category
func GetFavoritesByCategory(c fiber.Ctx) error {
	userID := c.Params("userId")

	var favorites []models.Favorite
	result := config.DB.Where("user_id = ?", userID).
		Preload("Product").
		Find(&favorites)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch favorites",
		})
	}

	// Group by category
	categoryMap := make(map[string][]models.Favorite)
	for _, favorite := range favorites {
		category := favorite.Product.Category
		if category == "" {
			category = "Uncategorized"
		}
		categoryMap[category] = append(categoryMap[category], favorite)
	}

	return c.JSON(fiber.Map{
		"user_id":               userID,
		"favorites_by_category": categoryMap,
		"total_favorites":       len(favorites),
	})
}

// ClearAllFavorites removes all favorites for a user
func ClearAllFavorites(c fiber.Ctx) error {
	userID := c.Params("userId")

	// Count favorites before deletion
	var count int64
	config.DB.Model(&models.Favorite{}).Where("user_id = ?", userID).Count(&count)

	// Delete all favorites for the user
	result := config.DB.Where("user_id = ?", userID).Delete(&models.Favorite{})
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to clear favorites",
		})
	}

	return c.JSON(fiber.Map{
		"message":       "All favorites cleared successfully",
		"removed_count": count,
	})
}
