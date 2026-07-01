package main

import (
	"log"

	"sm-client-backend/config"
	"sm-client-backend/handlers"
	"sm-client-backend/middleware"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables dari .env (jika ada)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on OS environment variables")
	}

	// Initialize Database Connection
	config.ConnectDB()

	app := fiber.New()

	// Enable CORS for frontend explicitly
	app.Use(func(c fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
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

	// WebSocket Route (misal: /ws/SM-BUDI)
	// TODO: Di-disable sementara karena github.com/gofiber/contrib/websocket
	// masih menggunakan Fiber v2, yang menyebabkan panic di Fiber v3.
	// app.Get("/ws/:slug", websocket.New(handlers.WsHandler))

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
