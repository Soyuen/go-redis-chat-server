package bootstrap

import (
	"fmt"

	"github.com/Soyuen/go-redis-chat-server/internal/application/cache"
	"github.com/Soyuen/go-redis-chat-server/internal/application/chat"
	"github.com/Soyuen/go-redis-chat-server/internal/application/pubsub"
	"github.com/Soyuen/go-redis-chat-server/internal/config"
	domainchat "github.com/Soyuen/go-redis-chat-server/internal/domain/chat"
	"github.com/Soyuen/go-redis-chat-server/internal/infrastructure/realtime"
	infraredis "github.com/Soyuen/go-redis-chat-server/internal/infrastructure/redis"
	"github.com/Soyuen/go-redis-chat-server/internal/presenter"
	"github.com/Soyuen/go-redis-chat-server/pkg/loggeriface"
	"github.com/redis/go-redis/v9"
)

type AppDependencies struct {
	RedisClient    *redis.Client
	RedisCache     cache.RedisCache
	RedisPubSub    pubsub.PubSub
	ChannelManager *realtime.ChannelManager
	Subscriber     *infraredis.RedisSubscriber
	Connection     *realtime.Connection
	ChatSvc        chat.ChatService
	Presenter      presenter.MessagePresenterInterface
}

func LoadRedisConfig() config.RedisConfig {
	return config.LoadRedisConfigFromEnv()
}

func InitRedisAdapter(cfg config.RedisConfig) (*infraredis.RedisAdapter, error) {
	return infraredis.NewRedisAdapter(cfg)
}

func InitRedisPubSub(client *redis.Client) pubsub.PubSub {
	return infraredis.NewRedisPubSubAdapter(client)
}

func InitChannelManager() *realtime.ChannelManager {
	return realtime.NewChannelManager()
}

func InitSubscriber(pubsub pubsub.PubSub, manager *realtime.ChannelManager, logger loggeriface.Logger) *infraredis.RedisSubscriber {
	return infraredis.NewRedisSubscriber(pubsub, manager, logger)
}

func InitConnectionHandler(manager *realtime.ChannelManager, logger loggeriface.Logger) *realtime.Connection {
	clientFactory := realtime.NewClientFactory(logger)
	return realtime.NewConnection(manager, logger, clientFactory)
}

func InitPresenter(logger loggeriface.Logger) presenter.MessagePresenterInterface {
	return presenter.NewMessagePresenter(logger)
}

func InitChatService(manager *realtime.ChannelManager, subscriber *infraredis.RedisSubscriber, presenter presenter.MessagePresenterInterface, memberRepo domainchat.ChatMemberRepository) chat.ChatService {
	return chat.NewChatService(manager, subscriber, presenter, memberRepo, Logger)
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
	memberRepo := infraredis.NewChatMemberRepository(redisAdapter)
	chatSvc := InitChatService(manager, subscriber, presenter, memberRepo)

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
