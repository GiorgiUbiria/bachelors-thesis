package handlers

import (
	"github.com/GiorgiUbiria/bachelor/config"
	"github.com/GiorgiUbiria/bachelor/models"
	"github.com/gofiber/fiber/v3"
)

// GetCart retrieves a cart by ID
func GetCart(c fiber.Ctx) error {
	id := c.Params("id")
	var cart models.Cart

	result := config.DB.Preload("Items.Product").First(&cart, id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Cart not found",
		})
	}

	return c.JSON(cart)
}

// AddCartItem adds a new item to the cart
func AddCartItem(c fiber.Ctx) error {
	id := c.Params("id")
	var cart models.Cart

	// Find the cart
	if err := config.DB.First(&cart, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Cart not found",
		})
	}

	// Parse request body
	var item models.CartItem
	if err := c.Bind().JSON(&item); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Set cart ID
	item.CartID = cart.ID

	// Create the item
	if err := config.DB.Create(&item).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to add item to cart",
		})
	}

	return c.JSON(item)
}

// UpdateCartItem updates an existing cart item
func UpdateCartItem(c fiber.Ctx) error {
	cartID := c.Params("id")
	itemID := c.Params("itemId")
	var item models.CartItem

	// Find the item
	if err := config.DB.Where("cart_id = ? AND id = ?", cartID, itemID).First(&item).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Cart item not found",
		})
	}

	// Parse request body
	var updateData struct {
		Quantity int `json:"quantity"`
	}
	if err := c.Bind().JSON(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Update the item
	item.Quantity = updateData.Quantity
	if err := config.DB.Save(&item).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update cart item",
		})
	}

	return c.JSON(item)
}

// RemoveCartItem removes an item from the cart
func RemoveCartItem(c fiber.Ctx) error {
	cartID := c.Params("id")
	itemID := c.Params("itemId")

	// Delete the item
	result := config.DB.Where("cart_id = ? AND id = ?", cartID, itemID).Delete(&models.CartItem{})
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to remove item from cart",
		})
	}

	if result.RowsAffected == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Cart item not found",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
