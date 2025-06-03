package main

import (
	"log"
	"net/http"
	"os"

	delivery "github.com/Soyuen/go-redis-chat-server/internal/delivery/http"
)

func main() {

	// Initialize the Google OAuth configuration
	r := delivery.NewRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	err := r.Run("0.0.0.0:" + port)
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("failed to run server: %v", err)
	}

}
