package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

var (
	DB *sql.DB
)

func ConnectDB() {
	// Kunci AES dari env
	hexKey := os.Getenv("HEX_KEY")
	ivKey := os.Getenv("IV_KEY")

	// Helper untuk dekripsi
	decryptOrFallback := func(envKey string) string {
		ct := os.Getenv(envKey)
		pt, err := DecryptAES(ct, hexKey, ivKey)
		if err != nil || pt == "" {
			return ct // fallback ke nilai asli jika gagal/kosong
		}
		return pt
	}

	dbIP := decryptOrFallback("DB_IP")
	dbPort := decryptOrFallback("DB_PORT")
	dbUsr := decryptOrFallback("DB_USR")
	dbPwd := decryptOrFallback("DB_PWD")
	dbNm := decryptOrFallback("DB_NM")

	// Jika tidak ada di env (misal development lokal), gunakan default
	if dbIP == "" {
		dbIP = "127.0.0.1"
	}
	if dbPort == "" {
		dbPort = "3306"
	}
	if dbUsr == "" {
		dbUsr = "root"
	}
	if dbNm == "" {
		dbNm = "sms_client"
	}

	log.Printf("Connecting to database at %s:%s as %s...\n", dbIP, dbPort, dbUsr)

	cfg := mysql.Config{
		User:                 dbUsr,
		Passwd:               dbPwd,
		Net:                  "tcp",
		Addr:                 fmt.Sprintf("%s:%s", dbIP, dbPort),
		DBName:               dbNm,
		AllowNativePasswords: true,
		ParseTime:            true,
		Timeout:              5 * time.Second,  // Batas waktu koneksi awal (5 detik)
		ReadTimeout:          5 * time.Second,   // Batas waktu baca
		WriteTimeout:         5 * time.Second,   // Batas waktu tulis
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatalln("Failed to init database:", err)
	}

	// Batasi pool koneksi
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Cek koneksi sesungguhnya
	if err := db.Ping(); err != nil {
		log.Println("⚠️  WARNING: Failed to ping database:", err)
		log.Println("⚠️  Server tetap berjalan, tapi fitur database belum aktif.")
		log.Println("⚠️  Pastikan IP, port, username, dan password di .env sudah benar.")
	} else {
		log.Println("✅ Database connected successfully")
	}
	DB = db
}

