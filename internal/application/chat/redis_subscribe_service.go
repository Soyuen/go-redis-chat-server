package chat

import (
	"context"
	"log"

	"github.com/Soyuen/go-redis-chat-server/pkg/pubsub"
	"github.com/Soyuen/go-redis-chat-server/pkg/realtimeiface"
)

type RedisSubscriberService struct {
	pubsub         pubsub.PubSub
	channelManager realtimeiface.ChatChannelManager
	cancelFuncs    map[string]context.CancelFunc
}

func NewRedisSubscriberService(pub pubsub.PubSub, manager realtimeiface.ChatChannelManager) *RedisSubscriberService {
	return &RedisSubscriberService{
		pubsub:         pub,
		channelManager: manager,
		cancelFuncs:    make(map[string]context.CancelFunc),
	}
}

func (s *RedisSubscriberService) Start(channel string) {
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

func (s *RedisSubscriberService) Stop() {
	for _, cancel := range s.cancelFuncs {
		cancel()
	}
	s.cancelFuncs = make(map[string]context.CancelFunc)
}
