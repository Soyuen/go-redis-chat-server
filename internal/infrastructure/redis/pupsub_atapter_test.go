package redis

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisPubSubAdapter_PublishAndSubscribe(t *testing.T) {
	testStr := "I'm a test"

	// Start miniredis mock Redis server
	s := miniredis.RunT(t)
	defer s.Close()

	// Create a redis.Client connected to the mock Redis server
	rdb := redis.NewClient(&redis.Options{
		Addr: s.Addr(),
	})

	adapter := NewRedisPubSubAdapter(rdb)

	ch := make(chan string, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 用 channel 確保訂閱已建立
	ready := make(chan struct{})

	// Use channel to ensure subscription is established
	go func() {
		sub, err := adapter.Subscribe(ctx, "test-channel")
		assert.NoError(t, err)
		if err != nil {
			return
		}
		close(ready) // Notify subscription is ready
		for {
			msg, err := sub.Receive(ctx)
			if err != nil {
				break
			}
			ch <- string(msg.Payload)
		}
	}()

	// Wait for the subscription to be established
	<-ready

	// Publish message
	err := adapter.Publish(ctx, "test-channel", []byte(testStr))
	assert.NoError(t, err)

	// Wait to receive message, up to 1 second
	select {
	case msg := <-ch:
		assert.Equal(t, testStr, msg)
	case <-time.After(time.Second):
		t.Fatal("timeout waiting for message")
	}
}
