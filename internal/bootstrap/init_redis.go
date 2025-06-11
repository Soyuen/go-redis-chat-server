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
	Presenter      *presenter.MessagePresenter
}

func InitRedisSubscriberService(logger loggeriface.Logger) (*AppDependencies, error) {
	// 1. Load configuration
	redisCfg := config.LoadRedisConfigFromEnv()

	// 2. Initialize RedisAdapter (contains *redis.Client)
	redisAdapter, err := appredis.NewRedisAdapter(redisCfg)
	if err != nil {
		return nil, fmt.Errorf("redis adapter init failed: %w", err)
	}

	// 3. Extract *redis.Client from redisAdapter
	client := redisAdapter.RawClient() // You need to implement this getter

	// 4. Create PubSub
	pub := appredis.NewRedisPubSubAdapter(client)

	// 5. Create ChannelManager
	manager := realtime.NewChannelManager()

	// 6. Create Subscriber
	subscriber := appredis.NewRedisSubscriber(pub, manager, logger)
	clientFactory := realtime.NewClientFactory(logger)
	connHandler := realtime.NewConnection(manager, logger, clientFactory)
	presenter := presenter.NewMessagePresenter(logger)
	chatSvc := chat.NewChatService(manager, subscriber, *presenter)

	// 7. Return all dependencies
	return &AppDependencies{
		RedisClient:    client,
		RedisCache:     redisAdapter,
		RedisPubSub:    pub,
		ChannelManager: manager,
		Subscriber:     subscriber,
		Connection:     connHandler,
		ChatSvc:        chatSvc,
		Presenter:      presenter,
	}, nil
}
