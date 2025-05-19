package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/GiorgiUbiria/bachelor/config"
	"github.com/GiorgiUbiria/bachelor/routes"
	"github.com/GiorgiUbiria/bachelor/routes/handlers"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	config.InitDB()

	app := fiber.New(fiber.Config{
		AppName: "Bachelor API v1",
		ErrorHandler: func(c fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	app.Use(logger.New())

	allowOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowOrigins == "" {
		allowOrigins = "http://localhost:3000,http://localhost:5173"
	}

	app.Use(cors.New(cors.Config{
		AllowOrigins: strings.Split(allowOrigins, ","),
		AllowMethods: []string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodHead,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
			fiber.MethodOptions,
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Requested-With",
			"X-Csrf-Token",
		},
		AllowCredentials: true,
		ExposeHeaders: []string{
			"Content-Length",
			"Content-Type",
			"Authorization",
			"X-Csrf-Token",
		},
		MaxAge: 3600,
	}))

	// Public endpoints (no CSRF required)
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
			"db":     config.DB != nil,
		})
	})

	// Auth endpoints (no CSRF required)
	auth := app.Group("/api/auth")
	auth.Post("/login", func(c fiber.Ctx) error {
		return handlers.Login(c, config.DB)
	})
	auth.Post("/register", func(c fiber.Ctx) error {
		return handlers.Register(c, config.DB)
	})

	// Setup CSRF protection for all other routes
	app.Use(csrf.New(csrf.Config{
		KeyLookup:      "header:X-Csrf-Token",
		CookieName:     "csrf_",
		CookieSameSite: "Lax",
		CookieSecure:   os.Getenv("APP_ENV") == "production",
		CookieHTTPOnly: true,
		IdleTimeout:    30 * time.Minute,
		ErrorHandler: func(c fiber.Ctx, err error) error {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "CSRF token missing or invalid",
			})
		},
	}))

	// CSRF token endpoint
	app.Get("/api/csrf-token", func(c fiber.Ctx) error {
		token := csrf.TokenFromContext(c)
		if token == "" {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not get CSRF token"})
		}
		return c.JSON(fiber.Map{"csrfToken": token})
	})

	// Protected routes
	routes.SetupRoutes(app, config.DB)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on http://localhost:%s", port)
	log.Fatal(app.Listen(":" + port))
}
