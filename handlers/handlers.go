package handlers

import (
	"log"

	"sm-client-backend/config"
	"sm-client-backend/models"

	"github.com/gofiber/fiber/v3"
)

// Login logic has been moved to handlers/auth.go

// GetVehicle retrieves vehicle data by slug
func GetVehicle(c fiber.Ctx) error {
	slug := c.Params("slug")

	// Query data kendaraan berdasarkan relasi tabel sms_client.cars_user dan sms_db.cars
	query := `
		SELECT 
			cu.car_id, cu.access_code, cu.full_name, cu.car_name, 
			COALESCE(c.plate_number, ''), COALESCE(c.status, ''), 
			COALESCE(c.restoration_type, ''), COALESCE(c.cover_photo_url, '')
		FROM sms_client.cars_user cu
		LEFT JOIN sms_db.cars c ON cu.car_id = c.id
		WHERE cu.access_code = ?
	`

	rows, err := config.DB.Query(query, slug)
	if err != nil {
		log.Printf("GetVehicle Query Error: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query database",
		})
	}
	defer rows.Close()

	var results []models.Vehicle

	for rows.Next() {
		var v models.Vehicle
		var carName string
		
		err := rows.Scan(
			&v.ID, 
			&v.ClientCode, 
			&v.ClientName, 
			&carName, 
			&v.LicensePlate, 
			&v.Status, 
			&v.RestorationType, 
			&v.BannerImage,
		)
		
		if err != nil {
			log.Printf("GetVehicle Scan Error: %v\n", err)
			continue
		}

		v.Brand = carName // Memetakan car_name ke Brand
		v.Timeline = []models.TimelineEvent{} // Dikosongkan sesuai feedback
		v.Gallery = []models.ProgressPhoto{}  // Dikosongkan sesuai feedback

		results = append(results, v)
	}

	if err = rows.Err(); err != nil {
		log.Printf("GetVehicle Rows Error: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed reading database rows",
		})
	}

	if len(results) > 0 {
		return c.JSON(results)
	}

	return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"error": "Vehicle not found",
	})
}
