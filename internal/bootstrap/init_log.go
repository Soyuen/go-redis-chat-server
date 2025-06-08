package bootstrap

import (
	"github.com/Soyuen/go-redis-chat-server/internal/config"
	"github.com/Soyuen/go-redis-chat-server/internal/infrastructure/logger"
	"github.com/Soyuen/go-redis-chat-server/pkg/loggeriface"
)

var Logger loggeriface.Logger

func InitLogger(cfg config.AppConfig) error {
	zapLogger, err := logger.NewZapLogger(cfg)
	if err != nil {
		return err
	}
	Logger = zapLogger
	return nil
}
