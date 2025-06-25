package bootstrap

import (
	"github.com/Soyuen/go-redis-chat-server/internal/config"
)

func Initialize(cfg config.AppConfig) (*AppDependencies, error) {
	appDependencies, err := InitAppDependencies(Logger)
	if err != nil {
		return nil, err
	}
	return appDependencies, err
}
