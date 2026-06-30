package handlers

import (
	"sm-client-backend/data"
	"sm-client-backend/models"

	"github.com/gofiber/fiber/v3"
)

// Login logic has been moved to handlers/auth.go

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
