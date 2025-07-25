// internal/infrastructure/redis/pubsub_adapter.go
package redis

import (
	"context"
	"fmt"

	"github.com/Soyuen/go-redis-chat-server/internal/application/pubsub"
	"github.com/redis/go-redis/v9"
)

var _ pubsub.PubSub = (*RedisPubSubAdapter)(nil)

type RedisPubSubAdapter struct {
	client *redis.Client
}

func NewRedisPubSubAdapter(client *redis.Client) pubsub.PubSub {
	return &RedisPubSubAdapter{client: client}
}

func (r *RedisPubSubAdapter) Publish(ctx context.Context, channel string, message []byte) error {
	return r.client.Publish(ctx, channel, message).Err()
}

func (r *RedisPubSubAdapter) Subscribe(ctx context.Context, channel string) (pubsub.Subscription, error) {
	redisPubSub := r.client.Subscribe(ctx, channel)
	return &redisSubscription{redisPubSub: redisPubSub}, nil
}

type redisSubscription struct {
	redisPubSub *redis.PubSub
}

func (s *redisSubscription) Receive(ctx context.Context) (*pubsub.Message, error) {
	msg, err := s.redisPubSub.ReceiveMessage(ctx)
	if err != nil {
		if err == context.Canceled || err == context.DeadlineExceeded {
			return nil, err
		}
		return nil, fmt.Errorf("error receiving message from Redis Pub/Sub: %w", err)
	}

	return &pubsub.Message{
		Channel: msg.Channel,
		Payload: []byte(msg.Payload),
	}, nil
}

func (s *redisSubscription) Unsubscribe(channels ...string) error {
	return s.redisPubSub.Unsubscribe(context.Background(), channels...)
}

func (s *redisSubscription) Close() error {
	return s.redisPubSub.Close()
}

var _ pubsub.Subscription = (*redisSubscription)(nil)
