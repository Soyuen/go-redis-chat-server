package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Soyuen/go-redis-chat-server/internal/config"
	delivery "github.com/Soyuen/go-redis-chat-server/internal/delivery/http"
	"github.com/Soyuen/go-redis-chat-server/internal/infrastructure/redis"
)

func main() {
	// Load Redis configuration from environment variables
	redisCfg := config.LoadRedisConfigFromEnv()

	// Initialize Redis Adapter
	redisAdapter, err := redis.NewRedisAdapter(redisCfg)
	if err != nil {
		log.Fatalf("failed to initialize Redis: %v", err)
	}

	// Initialize HTTP router and inject Redis Adapter (Dependency Inversion)
	r := delivery.NewRouter(redisAdapter)

	// Read HTTP port from environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	// Start HTTP server
	if err := r.Run("0.0.0.0:" + port); err != nil && err != http.ErrServerClosed {
		log.Fatalf("failed to run server: %v", err)
	}
}
