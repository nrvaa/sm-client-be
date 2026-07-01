//go:build ignore
// +build ignore

// File ini di-disable sementara karena WebSocket (gofiber/contrib/websocket)
// belum kompatibel dengan Fiber v3, dan Redis belum digunakan.
// Aktifkan kembali saat fitur real-time siap diimplementasikan.

package handlers

import (
	"context"
	"log"

	"github.com/gofiber/contrib/websocket"
	"sm-client-backend/config"
)

// WsHandler adalah handler untuk WebSocket connection
func WsHandler(c *websocket.Conn) {
	// Ambil slug/user_id dari parameter URL (misal: /ws/SM-BUDI)
	slug := c.Params("slug")

	// Subscribe ke channel Redis khusus user ini
	channelName := "update:" + slug
	pubsub := config.Redis.Subscribe(context.Background(), channelName)
	defer pubsub.Close()

	ch := pubsub.Channel()

	log.Printf("Client connected to WS: %s", slug)

	// Goroutine untuk membaca dari Redis dan kirim ke WebSocket
	go func() {
		for msg := range ch {
			log.Printf("Sending update to %s: %s", slug, msg.Payload)
			if err := c.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
				log.Println("Write error:", err)
				break
			}
		}
	}()

	// Loop untuk menjaga koneksi tetap hidup dan membaca pesan dari client (jika ada)
	var (
		mt  int
		msg []byte
		err error
	)
	for {
		if mt, msg, err = c.ReadMessage(); err != nil {
			log.Println("Client disconnected:", err)
			break
		}
		log.Printf("Received msg from %s (type %d): %s", slug, mt, msg)
	}
}
