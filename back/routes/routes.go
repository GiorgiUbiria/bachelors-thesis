package routes

import (
	"github.com/GiorgiUbiria/bachelor/middleware"
	"github.com/GiorgiUbiria/bachelor/routes/handlers"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	auth := app.Group("/api/auth")
	auth.Post("/login", func(c fiber.Ctx) error {
		return handlers.Login(c, db)
	})
	auth.Post("/register", func(c fiber.Ctx) error {
		return handlers.Register(c, db)
	})

	api := app.Group("/api", middleware.Protected())

	admin := api.Group("/admin", middleware.AdminOnly())
	admin.Get("/dashboard", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Admin dashboard",
		})
	})

	user := api.Group("/user")
	user.Get("/profile", func(c fiber.Ctx) error {
		user := c.Locals("user")
		return c.JSON(user)
	})

	// Product routes
	products := api.Group("/products")
	products.Get("/", handlers.GetProducts)
	products.Get("/search", handlers.SearchProducts)
	products.Get("/:id", handlers.GetProduct)
	products.Get("/category/:category", handlers.GetProductsByCategory)

	// Admin-only product routes
	adminProducts := admin.Group("/products")
	adminProducts.Post("/", handlers.CreateProduct)
	adminProducts.Put("/:id", handlers.UpdateProduct)
	adminProducts.Delete("/:id", handlers.DeleteProduct)
	adminProducts.Post("/:id/upload-image", handlers.UploadProductImage)

	users := api.Group("/users")
	users.Get("/", handlers.GetUsers)
	users.Get("/:id", handlers.GetUser)
	users.Get("/:id/activities", handlers.GetUserActivities)
	users.Get("/:id/favorites", handlers.GetUserFavorites)
	users.Get("/:id/cart", handlers.GetUserCart)
	users.Get("/:id/orders", handlers.GetUserOrders)

	cart := api.Group("/cart")
	cart.Get("/:id", handlers.GetCart)
	cart.Post("/:id/items", handlers.AddCartItem)
	cart.Put("/:id/items/:itemId", handlers.UpdateCartItem)
	cart.Delete("/:id/items/:itemId", handlers.RemoveCartItem)

	// Order routes
	orders := api.Group("/orders")
	orders.Get("/", handlers.GetOrders)
	orders.Get("/:id", handlers.GetOrder)
	orders.Post("/", handlers.CreateOrder)
	orders.Delete("/:id", handlers.CancelOrder)
	orders.Post("/:id/payment", handlers.ProcessPayment)
	orders.Get("/user/:userId", handlers.GetOrdersByUser)

	// Admin-only order routes
	adminOrders := admin.Group("/orders")
	adminOrders.Put("/:id/status", handlers.UpdateOrderStatus)

	// Activity tracking routes
	activities := api.Group("/activities")
	activities.Post("/view", handlers.TrackProductView)
	activities.Post("/click", handlers.TrackClick)
	activities.Post("/search", handlers.TrackSearch)
	activities.Post("/session", handlers.StartSession)
	activities.Post("/session/end", handlers.EndSession)
	activities.Post("/add-to-cart", handlers.TrackAddToCart)
	activities.Post("/remove-from-cart", handlers.TrackRemoveFromCart)
	activities.Post("/add-to-favorites", handlers.TrackAddToFavorites)
	activities.Get("/user/:userId/timeline", handlers.GetUserActivityTimeline)
	activities.Get("/session/:sessionId", handlers.GetSessionActivities)

	// Favorites management routes
	favorites := api.Group("/favorites")
	favorites.Post("/", handlers.AddToFavorites)
	favorites.Delete("/:id", handlers.RemoveFromFavorites)
	favorites.Get("/user/:userId", handlers.GetFavorites)
	favorites.Post("/bulk", handlers.BulkAddToFavorites)
	favorites.Delete("/bulk", handlers.BulkRemoveFromFavorites)
	favorites.Get("/user/:userId/status/:productId", handlers.CheckFavoriteStatus)
	favorites.Get("/user/:userId/categories", handlers.GetFavoritesByCategory)
	favorites.Delete("/user/:userId/clear", handlers.ClearAllFavorites)

	analytics := api.Group("/analytics")
	analytics.Get("/activities", handlers.GetActivityAnalytics)
	analytics.Get("/requests", handlers.GetRequestAnalytics)
	analytics.Get("/requests/recent", handlers.GetRecentRequestLogs)
	analytics.Get("/products/popular", handlers.GetPopularProducts)
	analytics.Get("/users/active", handlers.GetActiveUsers)

	// ML Integration routes
	ml := api.Group("/ml")
	ml.Get("/training-data", handlers.GetMLTrainingData)
	ml.Post("/retrain", handlers.TriggerMLRetraining)
	ml.Get("/recommendations/:userId", handlers.GetUserRecommendations)
	ml.Get("/user-behavior/:userId", handlers.GetUserBehaviorAnalysis)
	ml.Get("/trend-analysis", handlers.GetTrendAnalysis)
	ml.Get("/health", handlers.CheckMLServiceHealth)

	// Anomaly detection log endpoint
	api.Post("/log-request", handlers.LogRequestHandler)
}
