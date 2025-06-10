package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Soyuen/go-redis-chat-server/internal/bootstrap"
	"github.com/Soyuen/go-redis-chat-server/internal/config"
	"github.com/Soyuen/go-redis-chat-server/internal/delivery/router"
)

func main() {
	envCfg := config.LoadEnvConfig()

	err := bootstrap.InitLogger(envCfg)
	if err != nil {
		log.Fatalf("failed to init logger: %v", err)
		os.Exit(1)
	}
	appDependencies, err := bootstrap.Initialize(envCfg)

	if err != nil {
		bootstrap.Logger.Fatalw("[main] failed to initailize", "error", err)
	}

	r := router.NewRouter(appDependencies.ChannelManager, appDependencies.Connection, appDependencies.ChatSvc, *appDependencies.Presenter, bootstrap.Logger)
	port := envCfg.Port
	bootstrap.Logger.Infow("[main] starting server...",
		"port", port,
		"debug", envCfg.IsDebug,
	)

	if err := r.Run("0.0.0.0:" + port); err != nil && err != http.ErrServerClosed {
		bootstrap.Logger.Fatalw("[main] failed to run server", "error", err)
	}
}
