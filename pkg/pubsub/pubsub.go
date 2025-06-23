// pkg/pubsub/pubsub.go
package pubsub

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
