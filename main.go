package main

import (
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Create a new Fiber app
	app := fiber.New()

	// Root route (basic health check)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Welcome to the API Gateway!")
	})

	// Health check route
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("API Gateway is up and running!")
	})

	// Proxy route for the Auth service
	app.Get("/auth/*", func(c *fiber.Ctx) error {
		// Proxy the request to the Auth service running on localhost:5000
		resp, err := http.Get("http://localhost:5000" + c.OriginalURL())
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).SendString("Service Unavailable")
		}
		defer resp.Body.Close()

		// Copy the response from Auth service to the client
		return c.SendStream(resp.Body)
	})

	// Start the server on port 8080
	log.Fatal(app.Listen(":8080"))
}
