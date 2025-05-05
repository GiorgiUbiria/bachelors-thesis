package handlers

import (
	"github.com/GiorgiUbiria/bachelor/config"
	"github.com/GiorgiUbiria/bachelor/models"
	"github.com/gofiber/fiber/v3"
)

// GetOrders retrieves all orders
func GetOrders(c fiber.Ctx) error {
	var orders []models.Order

	result := config.DB.Preload("Items.Product").Find(&orders)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch orders",
		})
	}

	return c.JSON(orders)
}

// GetOrder retrieves a single order by ID
func GetOrder(c fiber.Ctx) error {
	id := c.Params("id")
	var order models.Order

	result := config.DB.Preload("Items.Product").First(&order, id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Order not found",
		})
	}

	return c.JSON(order)
}
