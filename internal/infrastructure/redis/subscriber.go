package redis

import (
	"context"

	"github.com/Soyuen/go-redis-chat-server/pkg/loggeriface"
	"github.com/Soyuen/go-redis-chat-server/pkg/pubsub"
	"github.com/Soyuen/go-redis-chat-server/pkg/realtimeiface"
)

var _ realtimeiface.ChannelEventSubscriber = (*RedisSubscriber)(nil)

type RedisSubscriber struct {
	pubsub         pubsub.PubSub
	channelManager realtimeiface.ChannelManager
	cancelFuncs    map[string]context.CancelFunc
	logger         loggeriface.Logger
}

func NewRedisSubscriber(pub pubsub.PubSub, manager realtimeiface.ChannelManager,
	logger loggeriface.Logger,
) *RedisSubscriber {
	return &RedisSubscriber{
		pubsub:         pub,
		channelManager: manager,
		cancelFuncs:    make(map[string]context.CancelFunc),
		logger:         logger,
	}
}

func (s *RedisSubscriber) Start(channel string) {
	if _, exists := s.cancelFuncs[channel]; exists {
		return // already listening
	}

	ctx, cancel := context.WithCancel(context.Background())
	s.cancelFuncs[channel] = cancel

	go func() {
		sub, err := s.pubsub.Subscribe(ctx, channel)
		if err != nil {
			s.logger.Fatalw("[RedisSubscriber] subscribe error", "channel", channel, "error", err)
			return
		}
		defer sub.Close()

		for {
			msg, err := sub.Receive(ctx)
			if err != nil {
				s.logger.Fatalw("[RedisSubscriber] receive error", "channel", channel, "error", err)
				break
			}
			s.channelManager.Broadcast(realtimeiface.Message{
				Channel: msg.Channel,
				Data:    msg.Payload,
			})
		}
	}()
}

func (s *RedisSubscriber) Stop() {
	for _, cancel := range s.cancelFuncs {
		cancel()
	}
	s.cancelFuncs = make(map[string]context.CancelFunc)
}
