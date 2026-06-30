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
		c.Set("Access-Control-Allow-Origin", "http://localhost:5173") // or "*"
		c.Set("Access-Control-Allow-Methods", "GET, POST, HEAD, PUT, DELETE, PATCH, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, Access-Control-Request-Private-Network")
		c.Set("Access-Control-Allow-Credentials", "true") // important for cookies
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
	
	// Auth routes
	api.Post("/login", func(c fiber.Ctx) error {
		return handlers.Login(c)
	})
	api.Post("/logout", middleware.Protected(), func(c fiber.Ctx) error {
		return handlers.Logout(c)
	})
	api.Get("/auth/check", middleware.Protected(), handlers.CheckSession)

	// Protected routes
	vehicles := api.Group("/vehicles")
	vehicles.Get("/:slug", middleware.Protected(), handlers.GetVehicle)

	log.Println("Server is running on port 3000")
	log.Fatal(app.Listen(":3000"))
}
