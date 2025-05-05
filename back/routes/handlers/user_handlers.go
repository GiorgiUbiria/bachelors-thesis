package handlers

import (
	"github.com/GiorgiUbiria/bachelor/config"
	"github.com/GiorgiUbiria/bachelor/models"
	"github.com/gofiber/fiber/v3"
)

// GetUsers retrieves all users
func GetUsers(c fiber.Ctx) error {
	var users []models.User
	result := config.DB.Find(&users)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}

	return c.JSON(users)
}

// GetUser retrieves a single user by ID
func GetUser(c fiber.Ctx) error {
	id := c.Params("id")
	var user models.User

	result := config.DB.First(&user, id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(user)
}

// GetUserActivities retrieves all activities for a user
func GetUserActivities(c fiber.Ctx) error {
	id := c.Params("id")
	var activities []models.UserActivity

	result := config.DB.Where("user_id = ?", id).Find(&activities)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user activities",
		})
	}

	return c.JSON(activities)
}

// GetUserFavorites retrieves all favorites for a user
func GetUserFavorites(c fiber.Ctx) error {
	id := c.Params("id")
	var favorites []models.Favorite

	result := config.DB.Where("user_id = ?", id).Preload("Product").Find(&favorites)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user favorites",
		})
	}

	return c.JSON(favorites)
}

// GetUserCart retrieves a user's cart with items
func GetUserCart(c fiber.Ctx) error {
	id := c.Params("id")
	var cart models.Cart

	result := config.DB.Where("user_id = ?", id).Preload("Items.Product").First(&cart)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Cart not found",
		})
	}

	return c.JSON(cart)
}

// GetUserOrders retrieves all orders for a user
func GetUserOrders(c fiber.Ctx) error {
	id := c.Params("id")
	var orders []models.Order

	result := config.DB.Where("user_id = ?", id).Preload("Items.Product").Find(&orders)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user orders",
		})
	}

	return c.JSON(orders)
}
