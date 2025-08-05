package pubsub

//go:generate mockgen -destination=mocks/mock_subscription.go -package=mocks github.com/Soyuen/go-redis-chat-server/internal/application/pubsub Subscription
//go:generate mockgen -destination=mocks/mock_pubsub.go -package=mocks github.com/Soyuen/go-redis-chat-server/internal/application/pubsub PubSub

import "context"

type Message struct {
	Channel string
	Payload []byte
}

type Subscription interface {
	Receive(ctx context.Context) (*Message, error)
	Unsubscribe(channels ...string) error
	Close() error
}

type PubSub interface {
	Publish(ctx context.Context, channel string, message []byte) error
	Subscribe(ctx context.Context, channel string) (Subscription, error)
}
