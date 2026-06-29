package main

import (
	"log"

	"sm-client-backend/handlers"
	"sm-client-backend/middleware"

	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()

	// Enable CORS for frontend explicitly
	app.Use(func(c fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, HEAD, PUT, DELETE, PATCH, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, Access-Control-Request-Private-Network")
		c.Set("Access-Control-Allow-Private-Network", "true")

		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}
		return c.Next()
	})

	// Public routes
	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("SM-Client API")
	})

	api := app.Group("/api")
	api.Post("/login", handlers.Login)

	// Protected routes
	vehicles := api.Group("/vehicles")
	vehicles.Get("/:slug", middleware.Protected(), handlers.GetVehicle)

	log.Println("Server is running on port 3000")
	log.Fatal(app.Listen(":3000"))
}
