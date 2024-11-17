package main

import (
	"log"

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
		return c.SendString("API Gateway is up and running!")
	})

	// Proxy route for the Auth service
	app.All("/api/auth/*", func(c *fiber.Ctx) error {
		target := "http://localhost:5000/api" + c.OriginalURL()[len("/api/auth"):]
		log.Printf("Routing request to Auth service: %s", target)
		return proxy.Do(c, target)
	})

	// Proxy route for the E-Store service
	app.All("/api/estore/*", func(c *fiber.Ctx) error {
		targetURL := "http://localhost:8081" + c.OriginalURL()[len("/api/estore"):]
		log.Printf("Routing request to E-Store service: %s", targetURL)
		return proxy.Do(c, targetURL)
	})

	// Proxy route for the Payments service
	app.All("/api/payments/*", func(c *fiber.Ctx) error {
		targetURL := "http://localhost:8082" + c.OriginalURL()[len("/api/payments"):]
		log.Printf("Routing request to Payments service: %s", targetURL)
		return proxy.Do(c, targetURL)
	})

	log.Println("Starting API Gateway on port 8080...")
	log.Fatal(app.Listen(":8080"))
}
