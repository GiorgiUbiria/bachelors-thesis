package routes

import (
	"github.com/GiorgiUbiria/bachelor/middleware"
	"github.com/GiorgiUbiria/bachelor/routes/handlers"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	// Auth routes
	auth := app.Group("/api/auth")
	auth.Post("/login", func(c fiber.Ctx) error {
		return handlers.Login(c, db)
	})
	auth.Post("/register", func(c fiber.Ctx) error {
		return handlers.Register(c, db)
	})

	// Protected routes
	api := app.Group("/api", middleware.Protected())

	// Admin routes
	admin := api.Group("/admin", middleware.AdminOnly())
	admin.Get("/dashboard", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Admin dashboard",
		})
	})

	// User routes
	user := api.Group("/user")
	user.Get("/profile", func(c fiber.Ctx) error {
		user := c.Locals("user")
		return c.JSON(user)
	})

	// Product routes
	products := api.Group("/products")
	products.Get("/", handlers.GetProducts)
	products.Get("/:id", handlers.GetProduct)
	products.Get("/category/:category", handlers.GetProductsByCategory)

	// User routes
	users := api.Group("/users")
	users.Get("/", handlers.GetUsers)
	users.Get("/:id", handlers.GetUser)
	users.Get("/:id/activities", handlers.GetUserActivities)
	users.Get("/:id/favorites", handlers.GetUserFavorites)
	users.Get("/:id/cart", handlers.GetUserCart)
	users.Get("/:id/orders", handlers.GetUserOrders)

	// Cart routes
	cart := api.Group("/cart")
	cart.Get("/:id", handlers.GetCart)
	cart.Post("/:id/items", handlers.AddCartItem)
	cart.Put("/:id/items/:itemId", handlers.UpdateCartItem)
	cart.Delete("/:id/items/:itemId", handlers.RemoveCartItem)

	// Order routes
	orders := api.Group("/orders")
	orders.Get("/", handlers.GetOrders)
	orders.Get("/:id", handlers.GetOrder)

	// Analytics routes
	analytics := api.Group("/analytics")
	analytics.Get("/activities", handlers.GetActivityAnalytics)
	analytics.Get("/requests", handlers.GetRequestAnalytics)
	analytics.Get("/products/popular", handlers.GetPopularProducts)
	analytics.Get("/users/active", handlers.GetActiveUsers)
}
