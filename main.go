package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/proxy"
)

func main() {
	app := fiber.New()

	// Enable CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000, http://localhost:5000",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Authorization,Content-Type",
		AllowCredentials: true,
	}))

	// Log incoming requests
	app.Use(func(c *fiber.Ctx) error {
		log.Printf("Incoming request: %s %s", c.Method(), c.OriginalURL())
		return c.Next()
	})

	// Root route
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the API Gateway!")
	})

	// Health check route
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "UP"})
	})

	// Proxy route for the Auth service
	app.All("/api/auth/*", func(c *fiber.Ctx) error {
		target := "http://rr-auth:5000" + c.OriginalURL()[len("/api/auth"):] // Correct dynamic routing
		if c.OriginalURL() == "/api/auth/health" {
			target = "http://rr-auth:5000/health" // Correct health endpoint
		}
		log.Printf("Routing request to Auth service: %s", target)
		return proxy.Do(c, target)
	})

	// Proxy route for the E-Store service
	app.All("/api/estore/*", func(c *fiber.Ctx) error {
		targetURL := "http://rr-store:8080" + c.OriginalURL()[len("/api/estore"):]
		if c.OriginalURL() == "/api/estore/health" {
			targetURL = "http://rr-store:8080/actuator/health" // ✅ Fix incorrect routing
		}
		log.Printf("Routing request to E-Store service: %s", targetURL)
		return proxy.Do(c, targetURL)
	})

	// Proxy route for the Payments service
	app.All("/api/payments/*", func(c *fiber.Ctx) error {
		targetURL := "http://rr-payments:8082" + c.OriginalURL()[len("/api/payments"):]
		log.Printf("Routing request to Payments service: %s", targetURL)
		return proxy.Do(c, targetURL)
	})

	// Get port from environment variable or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081" // Default to 8081 if no env var is set
	}

	log.Printf("Starting API Gateway on port %s...", port)

	// Listen on all network interfaces (0.0.0.0)
	if err := app.Listen("0.0.0.0:" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
