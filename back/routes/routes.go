package routes

import (
	"github.com/GiorgiUbiria/bachelor/routes/handlers"
	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(app *fiber.App) {
	// API group
	api := app.Group("/api")

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
