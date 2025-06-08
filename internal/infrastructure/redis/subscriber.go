package redis

import (
	"context"
	"log"

	"github.com/Soyuen/go-redis-chat-server/pkg/pubsub"
	"github.com/Soyuen/go-redis-chat-server/pkg/realtimeiface"
)

type RedisSubscriber struct {
	pubsub         pubsub.PubSub
	channelManager realtimeiface.ChatChannelManager
	cancelFuncs    map[string]context.CancelFunc
}

func NewRedisSubscriber(pub pubsub.PubSub, manager realtimeiface.ChatChannelManager) *RedisSubscriber {
	return &RedisSubscriber{
		pubsub:         pub,
		channelManager: manager,
		cancelFuncs:    make(map[string]context.CancelFunc),
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
			log.Printf("subscribe error on channel %s: %v", channel, err)
			return
		}
		defer sub.Close()

		for {
			msg, err := sub.Receive(ctx)
			if err != nil {
				log.Printf("receive error on channel %s: %v", channel, err)
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
