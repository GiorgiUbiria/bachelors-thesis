package config

import (
	"log"
	"math/rand"

	"github.com/GiorgiUbiria/bachelor/models"
	"golang.org/x/crypto/bcrypt"
)

var categories = []string{"Electronics", "Clothing", "Books", "Home", "Sports", "Beauty", "Toys", "Food"}
var statuses = []string{"pending", "processing", "shipped", "delivered", "cancelled"}

func SeedDatabase() {
	log.Println("Starting database seeding...")

	// Create admin user
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	admin := models.User{
		Username: "admin",
		Email:    "admin@example.com",
		Password: string(hashedPassword),
	}
	DB.Create(&admin)

	// Create regular users
	users := []models.User{}
	for i := 1; i <= 10; i++ {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)
		user := models.User{
			Username: "user" + string(rune('0'+i)),
			Email:    "user" + string(rune('0'+i)) + "@example.com",
			Password: string(hashedPassword),
		}
		DB.Create(&user)
		users = append(users, user)
	}

	// Create products
	products := []models.Product{}
	for i := 1; i <= 50; i++ {
		product := models.Product{
			Name:        "Product " + string(rune('0'+i)),
			Description: "Description for product " + string(rune('0'+i)),
			Price:       float64(rand.Intn(1000)) + 10.0,
			Stock:       rand.Intn(100) + 1,
			Category:    categories[rand.Intn(len(categories))],
			ImageURL:    "https://example.com/images/product" + string(rune('0'+i)) + ".jpg",
		}
		DB.Create(&product)
		products = append(products, product)
	}

	// Create orders
	for i := 0; i < 20; i++ {
		user := users[rand.Intn(len(users))]
		order := models.Order{
			UserID: user.ID,
			Status: statuses[rand.Intn(len(statuses))],
			Total:  0,
		}
		DB.Create(&order)

		// Create order items
		numItems := rand.Intn(5) + 1
		for j := 0; j < numItems; j++ {
			product := products[rand.Intn(len(products))]
			quantity := rand.Intn(5) + 1
			orderItem := models.OrderItem{
				OrderID:   order.ID,
				ProductID: product.ID,
				Quantity:  quantity,
				Price:     product.Price,
			}
			DB.Create(&orderItem)
			order.Total += product.Price * float64(quantity)
		}
		DB.Save(&order)
	}

	// Create favorites
	for _, user := range users {
		numFavorites := rand.Intn(10) + 1
		for i := 0; i < numFavorites; i++ {
			product := products[rand.Intn(len(products))]
			favorite := models.Favorite{
				UserID:    user.ID,
				ProductID: product.ID,
			}
			DB.Create(&favorite)
		}
	}

	// Create carts
	for _, user := range users {
		cart := models.Cart{
			UserID: user.ID,
		}
		DB.Create(&cart)

		// Create cart items
		numItems := rand.Intn(5) + 1
		for i := 0; i < numItems; i++ {
			product := products[rand.Intn(len(products))]
			cartItem := models.CartItem{
				CartID:    cart.ID,
				ProductID: product.ID,
				Quantity:  rand.Intn(3) + 1,
			}
			DB.Create(&cartItem)
		}
	}

	// Create user activities
	activityTypes := []string{"view", "click", "search", "add_to_cart", "remove_from_cart", "add_to_favorites"}
	for _, user := range users {
		numActivities := rand.Intn(20) + 5
		for i := 0; i < numActivities; i++ {
			activity := models.UserActivity{
				UserID:    user.ID,
				Type:      activityTypes[rand.Intn(len(activityTypes))],
				ProductID: &products[rand.Intn(len(products))].ID,
				Details:   "Sample activity details",
				SessionID: "session_" + string(rune('0'+rand.Intn(10))),
			}
			DB.Create(&activity)
		}
	}

	// Create request logs
	paths := []string{"/api/products", "/api/cart", "/api/favorites", "/api/orders", "/api/users"}
	methods := []string{"GET", "POST", "PUT", "DELETE"}
	categories := []string{"normal", "warning", "anomaly"}
	for i := 0; i < 100; i++ {
		log := models.RequestLog{
			IP:           "192.168.1." + string(rune('0'+rand.Intn(255))),
			Method:       methods[rand.Intn(len(methods))],
			Path:         paths[rand.Intn(len(paths))],
			Status:       200 + rand.Intn(400),
			UserAgent:    "Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
			UserID:       &users[rand.Intn(len(users))].ID,
			Category:     categories[rand.Intn(len(categories))],
			Details:      "Sample request details",
			ResponseTime: float64(rand.Intn(1000)) + 50.0,
		}
		DB.Create(&log)
	}

	log.Println("Database seeding completed!")
}
