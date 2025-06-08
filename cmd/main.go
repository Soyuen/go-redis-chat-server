package main

import (
	"log"
	"net/http"

	"github.com/Soyuen/go-redis-chat-server/internal/bootstrap"
	"github.com/Soyuen/go-redis-chat-server/internal/config"
	"github.com/Soyuen/go-redis-chat-server/internal/delivery/router"
)

func main() {
	// 讀取 Redis 設定
	appDependencies, err := bootstrap.InitRedisSubscriberService()
	if err != nil {
		log.Fatalf("failed InitRedisSubscriberService: %v", err)
	}
	// 初始化 HTTP 路由，注入 Redis Adapter（依賴反轉）
	r := router.NewRouter(appDependencies.ChannelManager, appDependencies.Connection, appDependencies.ChatSvc)

	// 讀取 HTTP Port
	envCfg := config.LoadEnvConfig()
	port := envCfg.Port
	// 啟動 HTTP server
	if err := r.Run("0.0.0.0:" + port); err != nil && err != http.ErrServerClosed {
		log.Fatalf("failed to run server: %v", err)
	}
}
