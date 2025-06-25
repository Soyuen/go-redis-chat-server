package bootstrap

import (
	"fmt"

	"github.com/Soyuen/go-redis-chat-server/internal/application/chat"
	"github.com/Soyuen/go-redis-chat-server/internal/config"
	"github.com/Soyuen/go-redis-chat-server/internal/infrastructure/realtime"
	appredis "github.com/Soyuen/go-redis-chat-server/internal/infrastructure/redis"
	"github.com/Soyuen/go-redis-chat-server/internal/presenter"
	"github.com/Soyuen/go-redis-chat-server/pkg/loggeriface"
	"github.com/Soyuen/go-redis-chat-server/pkg/pubsub"
	"github.com/redis/go-redis/v9"
)

type AppDependencies struct {
	RedisClient    *redis.Client
	RedisCache     *appredis.RedisAdapter
	RedisPubSub    pubsub.PubSub
	ChannelManager *realtime.ChannelManager
	Subscriber     *appredis.RedisSubscriber
	Connection     *realtime.Connection
	ChatSvc        chat.ChatService
	Presenter      presenter.MessagePresenterInterface
}

func LoadRedisConfig() config.RedisConfig {
	return config.LoadRedisConfigFromEnv()
}

func InitRedisAdapter(cfg config.RedisConfig) (*appredis.RedisAdapter, error) {
	return appredis.NewRedisAdapter(cfg)
}

func InitRedisPubSub(client *redis.Client) pubsub.PubSub {
	return appredis.NewRedisPubSubAdapter(client)
}

func InitChannelManager() *realtime.ChannelManager {
	return realtime.NewChannelManager()
}

func InitSubscriber(pubsub pubsub.PubSub, manager *realtime.ChannelManager, logger loggeriface.Logger) *appredis.RedisSubscriber {
	return appredis.NewRedisSubscriber(pubsub, manager, logger)
}

func InitConnectionHandler(manager *realtime.ChannelManager, logger loggeriface.Logger) *realtime.Connection {
	clientFactory := realtime.NewClientFactory(logger)
	return realtime.NewConnection(manager, logger, clientFactory)
}

func InitPresenter(logger loggeriface.Logger) presenter.MessagePresenterInterface {
	return presenter.NewMessagePresenter(logger)
}

func InitChatService(manager *realtime.ChannelManager, subscriber *appredis.RedisSubscriber, presenter presenter.MessagePresenterInterface) chat.ChatService {
	return chat.NewChatService(manager, subscriber, presenter)
}

func InitAppDependencies(logger loggeriface.Logger) (*AppDependencies, error) {
	redisCfg := LoadRedisConfig()

	redisAdapter, err := InitRedisAdapter(redisCfg)
	if err != nil {
		return nil, fmt.Errorf("redis adapter init failed: %w", err)
	}

	client := redisAdapter.RawClient()

	pubsub := InitRedisPubSub(client)
	manager := InitChannelManager()
	subscriber := InitSubscriber(pubsub, manager, logger)
	conn := InitConnectionHandler(manager, logger)
	presenter := InitPresenter(logger)
	chatSvc := InitChatService(manager, subscriber, presenter)

	return &AppDependencies{
		RedisClient:    client,
		RedisCache:     redisAdapter,
		RedisPubSub:    pubsub,
		ChannelManager: manager,
		Subscriber:     subscriber,
		Connection:     conn,
		ChatSvc:        chatSvc,
		Presenter:      presenter,
	}, nil
}
