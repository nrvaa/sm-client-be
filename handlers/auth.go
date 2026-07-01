package handlers

import (
	"time"

	"sm-client-backend/auth"
	"sm-client-backend/data"
	"sm-client-backend/models"

	"strings"

	"github.com/gofiber/fiber/v3"
)

type LoginBody struct {
	URLToken   string `json:"url_token"`
	AccessCode string `json:"access_code"`
	Username   string `json:"username"`
}

// 1. ENDPOINT LOGIN
func Login(c fiber.Ctx) error {
	var req LoginBody
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	loginIdentifier := req.AccessCode
	if loginIdentifier == "" {
		loginIdentifier = req.Username
	}

	// nyocokin url_token dan access_code ke mock data (tar migrasi ke MySQL)
	// contoh validasi sederhana sesuai alur dokumen
	var selectedUser *models.User // Ganti dengan tipe user mock 
	for _, user := range data.MockUsers {
		// Mock mapping: asumsikan username digunakan atau kamu punya field url_token nanti
		if strings.EqualFold(user.Username, loginIdentifier) || strings.EqualFold(user.Slug, loginIdentifier) {
			u := user // copy
			selectedUser = &u
			break
		}
	}

	if selectedUser == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	// Generate JWT Token
	tokenString, err := auth.GenerateToken(selectedUser.ID, selectedUser.Slug)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	// SET COOKIE HTTPONLY (Biar ga mental pas di-refresh)
	c.Cookie(&fiber.Cookie{
		Name:     "sm_jwt",
		Value:    tokenString,
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		Path:     "/",
		HTTPOnly: true,
		Secure:   false, // Set jadi true pas production (HTTPS)
		SameSite: "Strict",
	})

	return c.JSON(fiber.Map{
		"slug":  selectedUser.Slug,
		"token": tokenString, // Tetap dikembalikan untuk opsi frontend
	})
}

// 2. ENDPOINT CHECK SESSION (Menghindari mental logout)
func CheckSession(c fiber.Ctx) error {
	// Ambil data slug yang sudah diekstrak dan diset oleh Middleware Protected()
	slug := c.Locals("slug").(string)

	// Jika sampai di sini, artinya cookie dibaca, JWT valid, dan sesi Redis aktif
	return c.JSON(fiber.Map{
		"slug": slug,
	})
}

// 3. ENDPOINT LOGOUT
func Logout(c fiber.Ctx) error {
	// Bersihkan Cookie di Browser dengan memundurkan waktu expired
	c.Cookie(&fiber.Cookie{
		Name:    "sm_jwt",
		Value:   "",
		Expires: time.Now().Add(-1 * time.Hour),
		Path:    "/",
	})

	return c.JSON(fiber.Map{"message": "Successfully logged out"})
}
