package repository

import (
	"fmt"
	"sm-client-backend/config"
	"sm-client-backend/models"
)

// DbUser adalah struct pemetaan dari tabel cars_user di DB
type DbUser struct {
	ID         int    `db:"id"`
	CarID      string `db:"car_id"`
	UserID     string `db:"user_id"`
	FullName   string `db:"full_name"`
	CarName    string `db:"car_name"`
	AccessCode string `db:"access_code"`
}

// GetUserByAccessCode mencari user berdasarkan access_code
func GetUserByAccessCode(accessCode string) (*models.User, error) {
	var dbUser DbUser
	// Cari user by access_code
	query := `SELECT id, car_id, user_id, full_name, car_name FROM cars_user WHERE access_code = ? LIMIT 1`
	
	row := config.DB.QueryRow(query, accessCode)
	err := row.Scan(&dbUser.ID, &dbUser.CarID, &dbUser.UserID, &dbUser.FullName, &dbUser.CarName)
	
	if err != nil {
		return nil, fmt.Errorf("invalid access code")
	}

	return &models.User{
		ID:       fmt.Sprintf("%d", dbUser.ID),
		Username: dbUser.FullName,
		Slug:     accessCode, // menggunakan accessCode sebagai identifier url_token
	}, nil
}

func GetVehicleBySlug(slug string) (*models.Vehicle, error) {
	var v models.Vehicle
	var coverPhoto string
	
	// JOIN antara sms_client.cars_user dan sms_db.cars
	query := `
		SELECT 
			u.car_id, 
			u.user_id, 
			u.full_name, 
			u.car_name,
			COALESCE(c.status, 'In Progress'),
			COALESCE(c.cover_photo_url, '')
		FROM sms_client.cars_user u
		LEFT JOIN sms_db.cars c ON u.car_id = c.id
		WHERE u.user_id = ? 
		LIMIT 1
	`
	
	row := config.DB.QueryRow(query, slug)
	err := row.Scan(&v.ID, &v.ClientCode, &v.ClientName, &v.Brand, &v.Status, &coverPhoto)
	if err != nil {
		return nil, err
	}

	v.BannerImage = coverPhoto
	if v.BannerImage == "" {
		v.BannerImage = "/images/bajaj_4.jpg" // Fallback banner
	}
	v.CompletionPercentage = 50 // Dummy untuk sekarang

	// Inject dummy Timeline
	v.Timeline = []models.TimelineEvent{
		{ID: "t1", Status: "Unit Tiba", Date: "2026-06-25", Description: "Kendaraan tiba di bengkel dan masuk antrian inspeksi.", Completed: true},
		{ID: "t2", Status: "Inspeksi & Pembongkaran", Date: "2026-06-28", Description: "Pembongkaran bodi dan pengecekan rangka.", Image: "/images/becak_1.jpg", Completed: true},
		{ID: "t3", Status: "Pengecatan", Date: "2026-07-02", Description: "Proses cat dasar dan epoxy.", Completed: false},
	}

	// Inject dummy Gallery
	v.Gallery = []models.ProgressPhoto{
		{ID: "g1", URL: "/images/bajaj_1.jpg", Title: "Kondisi Awal Depan", Date: "2026-06-25", Stage: "Inspeksi"},
		{ID: "g2", URL: "/images/bajaj_2.jpg", Title: "Interior", Date: "2026-06-25", Stage: "Inspeksi"},
		{ID: "g3", URL: "/images/becak_2.jpg", Title: "Pelepasan Panel", Date: "2026-06-28", Stage: "Pembongkaran"},
		{ID: "g4", URL: "/images/becak_3.jpg", Title: "Sasis Utama", Date: "2026-06-29", Stage: "Pembongkaran"},
	}

	return &v, nil
}
