package handlers

import (
	"strconv"
	"time"

	"github.com/GiorgiUbiria/bachelor/config"
	"github.com/GiorgiUbiria/bachelor/models"
	"github.com/gofiber/fiber/v3"
)

// GetOrders retrieves all orders with optional filtering
func GetOrders(c fiber.Ctx) error {
	var orders []models.Order

	// Get query parameters for filtering
	status := c.Query("status")
	userID := c.Query("user_id")

	// Pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset := (page - 1) * limit

	query := config.DB.Model(&models.Order{}).Preload("Items.Product").Preload("User")

	// Apply filters
	if status != "" {
		query = query.Where("status = ?", status)
	}

	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	// Get total count
	var total int64
	query.Count(&total)

	// Apply pagination and fetch results
	result := query.Offset(offset).Limit(limit).Find(&orders)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch orders",
		})
	}

	return c.JSON(fiber.Map{
		"orders": orders,
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// GetOrder retrieves a single order by ID
func GetOrder(c fiber.Ctx) error {
	id := c.Params("id")
	var order models.Order

	result := config.DB.Preload("Items.Product").Preload("User").First(&order, id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Order not found",
		})
	}

	return c.JSON(order)
}

// CreateOrder creates a new order from user's cart
func CreateOrder(c fiber.Ctx) error {
	var requestData struct {
		UserID uint `json:"user_id"`
	}

	if err := c.Bind().JSON(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get user's cart
	var cart models.Cart
	result := config.DB.Where("user_id = ?", requestData.UserID).Preload("Items.Product").First(&cart)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Cart not found",
		})
	}

	if len(cart.Items) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cart is empty",
		})
	}

	// Start transaction
	tx := config.DB.Begin()

	// Create order
	order := models.Order{
		UserID: requestData.UserID,
		Status: "pending",
		Total:  0,
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create order",
		})
	}

	// Create order items from cart items
	var total float64
	for _, cartItem := range cart.Items {
		// Check stock availability
		if cartItem.Product.Stock < cartItem.Quantity {
			tx.Rollback()
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Insufficient stock for product: " + cartItem.Product.Name,
			})
		}

		orderItem := models.OrderItem{
			OrderID:   order.ID,
			ProductID: cartItem.ProductID,
			Quantity:  cartItem.Quantity,
			Price:     cartItem.Product.Price,
		}

		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create order item",
			})
		}

		// Update product stock
		cartItem.Product.Stock -= cartItem.Quantity
		if err := tx.Save(&cartItem.Product).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update product stock",
			})
		}

		total += cartItem.Product.Price * float64(cartItem.Quantity)
	}

	// Update order total
	order.Total = total
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update order total",
		})
	}

	// Clear cart items
	if err := tx.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to clear cart",
		})
	}

	// Commit transaction
	tx.Commit()

	// Reload order with items and product details
	config.DB.Preload("Items.Product").Preload("User").First(&order, order.ID)

	return c.Status(fiber.StatusCreated).JSON(order)
}

// UpdateOrderStatus updates the status of an order (Admin only)
func UpdateOrderStatus(c fiber.Ctx) error {
	id := c.Params("id")
	var order models.Order

	// Find the order
	result := config.DB.First(&order, id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Order not found",
		})
	}

	var requestData struct {
		Status string `json:"status"`
	}

	if err := c.Bind().JSON(&requestData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate status
	validStatuses := map[string]bool{
		"pending":    true,
		"processing": true,
		"shipped":    true,
		"delivered":  true,
		"cancelled":  true,
	}

	if !validStatuses[requestData.Status] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid status. Valid statuses: pending, processing, shipped, delivered, cancelled",
		})
	}

	// Update order status
	order.Status = requestData.Status
	result = config.DB.Save(&order)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update order status",
		})
	}

	return c.JSON(order)
}

// CancelOrder cancels an order
func CancelOrder(c fiber.Ctx) error {
	id := c.Params("id")
	var order models.Order

	// Find the order with items
	result := config.DB.Preload("Items.Product").First(&order, id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Order not found",
		})
	}

	// Check if order can be cancelled
	if order.Status == "delivered" || order.Status == "cancelled" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot cancel order with status: " + order.Status,
		})
	}

	// Start transaction
	tx := config.DB.Begin()

	// Restore product stock
	for _, item := range order.Items {
		var product models.Product
		if err := tx.First(&product, item.ProductID).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to find product for stock restoration",
			})
		}

		product.Stock += item.Quantity
		if err := tx.Save(&product).Error; err != nil {
			tx.Rollback()
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to restore product stock",
			})
		}
	}

	// Update order status to cancelled
	order.Status = "cancelled"
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to cancel order",
		})
	}

	// Commit transaction
	tx.Commit()

	return c.JSON(order)
}

// GetOrdersByUser retrieves all orders for a specific user
func GetOrdersByUser(c fiber.Ctx) error {
	userID := c.Params("userId")
	var orders []models.Order

	// Pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	query := config.DB.Where("user_id = ?", userID).Preload("Items.Product")

	// Get total count
	var total int64
	query.Model(&models.Order{}).Count(&total)

	// Apply pagination and fetch results
	result := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&orders)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user orders",
		})
	}

	return c.JSON(fiber.Map{
		"orders": orders,
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// ProcessPayment processes payment for an order
func ProcessPayment(c fiber.Ctx) error {
	id := c.Params("id")
	var order models.Order

	// Find the order
	result := config.DB.First(&order, id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Order not found",
		})
	}

	var paymentData struct {
		PaymentMethod string  `json:"payment_method"`
		Amount        float64 `json:"amount"`
		CardNumber    string  `json:"card_number,omitempty"`
		ExpiryDate    string  `json:"expiry_date,omitempty"`
		CVV           string  `json:"cvv,omitempty"`
	}

	if err := c.Bind().JSON(&paymentData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate payment amount
	if paymentData.Amount != order.Total {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Payment amount does not match order total",
		})
	}

	// Validate payment method
	validMethods := map[string]bool{
		"credit_card":   true,
		"debit_card":    true,
		"paypal":        true,
		"bank_transfer": true,
	}

	if !validMethods[paymentData.PaymentMethod] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid payment method",
		})
	}

	// Simulate payment processing
	// In a real application, this would integrate with a payment gateway
	paymentSuccess := true // Simulate successful payment

	if !paymentSuccess {
		return c.Status(fiber.StatusPaymentRequired).JSON(fiber.Map{
			"error": "Payment processing failed",
		})
	}

	// Update order status to processing
	order.Status = "processing"
	result = config.DB.Save(&order)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update order after payment",
		})
	}

	// Create payment record (you would need to create a Payment model)
	// For now, we'll just return success

	return c.JSON(fiber.Map{
		"message": "Payment processed successfully",
		"order":   order,
		"payment": fiber.Map{
			"method":       paymentData.PaymentMethod,
			"amount":       paymentData.Amount,
			"status":       "completed",
			"processed_at": time.Now(),
		},
	})
}
