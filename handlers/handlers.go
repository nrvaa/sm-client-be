package handlers

import (
	"fmt"
	"sm-client-backend/auth"
	"sm-client-backend/data"
	"sm-client-backend/models"

	"github.com/gofiber/fiber/v3"
)

type LoginRequest struct {
	Username string `json:"username"`
	// Password string `json:"password"`
}

// Login handles user authentication and JWT generation
func Login(c fiber.Ctx) error {
	var req LoginRequest
	if err := c.Bind().JSON(&req); err != nil {
		fmt.Println("Bind Error:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}
	fmt.Printf("Received Login Request - Username: '%s'", req.Username)

	// Verify against mock users
	for _, user := range data.MockUsers {
		if user.Username == req.Username {
			// Valid credentials, generate token
			token, err := auth.GenerateToken(user.ID, user.Slug)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to generate token",
				})
			}
			return c.JSON(fiber.Map{
				"token": token,
				"slug":  user.Slug,
			})
		}
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": "Invalid credentials",
	})
}

// GetVehicle retrieves vehicle data by slug
func GetVehicle(c fiber.Ctx) error {
	slug := c.Params("slug")

	var results []models.Vehicle
	for _, v := range data.InitialVehicles {
		if v.ClientCode == slug {
			results = append(results, v)
		}
	}

	if len(results) > 0 {
		return c.JSON(results)
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": "Vehicle not found",
	})
}
