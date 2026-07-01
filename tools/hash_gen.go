package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run tools/hash_gen.go <kode_akses_asli>")
		fmt.Println("Contoh: go run tools/hash_gen.go MY-SECRET-123")
		os.Exit(1)
	}

	plainText := os.Args[1]

	// Menggunakan cost 12 sebagai standar yang baik (seimbang antara keamanan & kecepatan)
	hash, err := bcrypt.GenerateFromPassword([]byte(plainText), 12)
	if err != nil {
		fmt.Println("Error generating hash:", err)
		os.Exit(1)
	}

	fmt.Println("=== BCRYPT HASH GENERATOR ===")
	fmt.Println("Plaintext  :", plainText)
	fmt.Println("Bcrypt Hash:", string(hash))
	fmt.Println("\nCara pakai:")
	fmt.Printf("UPDATE sms_client.cars_user SET access_code = '%s' WHERE user_id = 'SM-BUDI';\n", string(hash))
}
