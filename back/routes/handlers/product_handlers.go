package handlers

import (
	"strconv"
	"strings"

	"github.com/GiorgiUbiria/bachelor/config"
	"github.com/GiorgiUbiria/bachelor/models"
	"github.com/gofiber/fiber/v3"
)

// GetProducts retrieves all products with optional filtering
func GetProducts(c fiber.Ctx) error {
	var products []models.Product

	// Get query parameters for filtering
	category := c.Query("category")
	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")
	search := c.Query("search")
	sortBy := c.Query("sort_by", "id")
	sortOrder := c.Query("sort_order", "asc")

	// Pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset := (page - 1) * limit

	query := config.DB.Model(&models.Product{})

	// Apply filters
	if category != "" {
		query = query.Where("category = ?", category)
	}

	if minPrice != "" {
		if price, err := strconv.ParseFloat(minPrice, 64); err == nil {
			query = query.Where("price >= ?", price)
		}
	}

	if maxPrice != "" {
		if price, err := strconv.ParseFloat(maxPrice, 64); err == nil {
			query = query.Where("price <= ?", price)
		}
	}

	if search != "" {
		searchTerm := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", searchTerm, searchTerm)
	}

	// Apply sorting
	orderClause := sortBy
	if sortOrder == "desc" {
		orderClause += " DESC"
	}
	query = query.Order(orderClause)

	// Get total count for pagination
	var total int64
	query.Count(&total)

	// Apply pagination and fetch results
	result := query.Offset(offset).Limit(limit).Find(&products)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch products",
		})
	}

	return c.JSON(fiber.Map{
		"products": products,
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
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

// CreateProduct creates a new product (Admin only)
func CreateProduct(c fiber.Ctx) error {
	var product models.Product

	if err := c.Bind().JSON(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if product.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Product name is required",
		})
	}

	if product.Price <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Product price must be greater than 0",
		})
	}

	if product.Stock < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Product stock cannot be negative",
		})
	}

	result := config.DB.Create(&product)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create product",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(product)
}

// UpdateProduct updates an existing product (Admin only)
func UpdateProduct(c fiber.Ctx) error {
	id := c.Params("id")
	var product models.Product

	// Find the product
	result := config.DB.First(&product, id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	// Parse update data
	var updateData models.Product
	if err := c.Bind().JSON(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate fields if they are being updated
	if updateData.Name != "" {
		product.Name = updateData.Name
	}

	if updateData.Description != "" {
		product.Description = updateData.Description
	}

	if updateData.Price > 0 {
		product.Price = updateData.Price
	}

	if updateData.Stock >= 0 {
		product.Stock = updateData.Stock
	}

	if updateData.Category != "" {
		product.Category = updateData.Category
	}

	if updateData.ImageURL != "" {
		product.ImageURL = updateData.ImageURL
	}

	// Save the updated product
	result = config.DB.Save(&product)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update product",
		})
	}

	return c.JSON(product)
}

// DeleteProduct deletes a product (Admin only)
func DeleteProduct(c fiber.Ctx) error {
	id := c.Params("id")
	var product models.Product

	// Check if product exists
	result := config.DB.First(&product, id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	// Check if product is referenced in any orders or carts
	var orderItemCount int64
	config.DB.Model(&models.OrderItem{}).Where("product_id = ?", id).Count(&orderItemCount)

	var cartItemCount int64
	config.DB.Model(&models.CartItem{}).Where("product_id = ?", id).Count(&cartItemCount)

	if orderItemCount > 0 || cartItemCount > 0 {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Cannot delete product that is referenced in orders or carts",
		})
	}

	// Delete the product
	result = config.DB.Delete(&product)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete product",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// SearchProducts performs full-text search on products
func SearchProducts(c fiber.Ctx) error {
	query := c.Query("q")
	if query == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Search query is required",
		})
	}

	// Pagination
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	offset := (page - 1) * limit

	var products []models.Product
	searchTerm := "%" + strings.ToLower(query) + "%"

	dbQuery := config.DB.Where(
		"LOWER(name) LIKE ? OR LOWER(description) LIKE ? OR LOWER(category) LIKE ?",
		searchTerm, searchTerm, searchTerm,
	)

	// Get total count
	var total int64
	dbQuery.Model(&models.Product{}).Count(&total)

	// Get results with pagination
	result := dbQuery.Offset(offset).Limit(limit).Find(&products)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to search products",
		})
	}

	return c.JSON(fiber.Map{
		"query":    query,
		"products": products,
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// UploadProductImage handles product image upload (Admin only)
func UploadProductImage(c fiber.Ctx) error {
	id := c.Params("id")
	var product models.Product

	// Find the product
	result := config.DB.First(&product, id)
	if result.Error != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	// Get the uploaded file
	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No image file provided",
		})
	}

	// Validate file type
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
	}

	if !allowedTypes[file.Header.Get("Content-Type")] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid file type. Only JPEG, PNG, and GIF are allowed",
		})
	}

	// Validate file size (max 5MB)
	maxSize := int64(5 * 1024 * 1024) // 5MB
	if file.Size > maxSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File size too large. Maximum size is 5MB",
		})
	}

	// Generate filename and save path
	filename := "product_" + id + "_" + file.Filename
	savePath := "./uploads/products/" + filename

	// Save the file
	if err := c.SaveFile(file, savePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save image file",
		})
	}

	// Update product with new image URL
	product.ImageURL = "/uploads/products/" + filename
	result = config.DB.Save(&product)
	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update product image URL",
		})
	}

	return c.JSON(fiber.Map{
		"message":   "Image uploaded successfully",
		"image_url": product.ImageURL,
		"product":   product,
	})
}
