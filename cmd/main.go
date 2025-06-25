package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	r := router.NewRouter(appDependencies.ChannelManager, appDependencies.Connection, appDependencies.ChatSvc, appDependencies.Presenter, bootstrap.Logger)
	port := envCfg.Port

	server := &http.Server{
		Addr:    "0.0.0.0:" + port,
		Handler: r,
	}

	go func() {
		bootstrap.Logger.Infow("[main] starting server...", "port", port, "debug", envCfg.IsDebug)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			bootstrap.Logger.Fatalw("[main] failed to run server", "error", err)
		}
	}()

	gracefulShutdown(server, appDependencies)
}

func gracefulShutdown(server *http.Server, deps *bootstrap.AppDependencies) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	bootstrap.Logger.Infow("[main] shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// close HTTP server
	if err := server.Shutdown(ctx); err != nil {
		bootstrap.Logger.Errorw("[main] server forced to shutdown", "error", err)
	}

	// close others resources
	shutdownResources(deps)

	bootstrap.Logger.Infow("[main] server gracefully stopped")
}

func shutdownResources(deps *bootstrap.AppDependencies) {
	deps.ChannelManager.CloseAllChannels()
	if err := deps.RedisClient.Close(); err != nil {
		bootstrap.Logger.Errorw("failed to close Redis", "error", err)
	}
}
