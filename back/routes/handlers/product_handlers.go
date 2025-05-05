package handlers

import (
	"github.com/GiorgiUbiria/bachelor/config"
	"github.com/GiorgiUbiria/bachelor/models"
	"github.com/gofiber/fiber/v3"
)

// GetProducts retrieves all products with optional filtering
func GetProducts(c fiber.Ctx) error {
	var products []models.Product
	result := config.DB.Find(&products)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch products",
		})
	}

	return c.JSON(products)
}

// GetProduct retrieves a single product by ID
func GetProduct(c fiber.Ctx) error {
	id := c.Params("id")
	var product models.Product

	result := config.DB.First(&product, id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	return c.JSON(product)
}

// GetProductsByCategory retrieves products filtered by category
func GetProductsByCategory(c fiber.Ctx) error {
	category := c.Params("category")
	var products []models.Product

	result := config.DB.Where("category = ?", category).Find(&products)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch products",
		})
	}

	return c.JSON(products)
}
