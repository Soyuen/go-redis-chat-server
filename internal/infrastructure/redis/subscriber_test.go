package redis_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	loggermock "github.com/Soyuen/go-redis-chat-server/internal/infrastructure/logger/mocks"
	realtimemock "github.com/Soyuen/go-redis-chat-server/internal/infrastructure/realtime/mocks"
	"github.com/Soyuen/go-redis-chat-server/internal/infrastructure/redis"
	pubsubmock "github.com/Soyuen/go-redis-chat-server/internal/infrastructure/redis/mocks"
	"github.com/Soyuen/go-redis-chat-server/pkg/pubsub"
	"github.com/Soyuen/go-redis-chat-server/pkg/realtimeiface"

	"github.com/golang/mock/gomock"
)

func TestRedisSubscriber_Start_BroadcastsReceivedMessages(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Arrange
	mockPubSub := pubsubmock.NewMockPubSub(ctrl)
	mockSubscription := pubsubmock.NewMockSubscription(ctrl)
	mockChannelManager := realtimemock.NewMockChannelManager(ctrl)
	mockLogger := loggermock.NewMockLogger(ctrl)

	channel := "test-room"
	msgData := []byte("hello")

	// Simulate Subscription behavior
	mockPubSub.EXPECT().
		Subscribe(gomock.Any(), channel).
		Return(mockSubscription, nil)
		// First message is received successfully; second triggers context.Canceled to exit the loop.
	gomock.InOrder(
		mockSubscription.EXPECT().
			Receive(gomock.Any()).
			Return(&pubsub.Message{
				Channel: channel,
				Payload: msgData,
			}, nil),
		mockSubscription.EXPECT().
			Receive(gomock.Any()).
			Return(nil, context.Canceled),
	)

	mockSubscription.EXPECT().Close()

	mockChannelManager.EXPECT().
		Broadcast(realtimeiface.Message{
			Channel: channel,
			Data:    msgData,
		})

	// No need to mock Fatalw because context.Canceled won't trigger Fatalw
	// (assuming you've updated the production code accordingly).
	subscriber := redis.NewRedisSubscriber(mockPubSub, mockChannelManager, mockLogger)
	subscriber.Start(channel)

	// Wait for the goroutine to finish execution
	time.Sleep(10 * time.Millisecond)

	// Stop to cleanup
	subscriber.Stop()
}

func TestRedisSubscriber_Start_DoesNotResubscribeIfAlreadyStarted(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPubSub := pubsubmock.NewMockPubSub(ctrl)
	mockChannelManager := realtimemock.NewMockChannelManager(ctrl)
	mockLogger := loggermock.NewMockLogger(ctrl)

	subscriber := redis.NewRedisSubscriber(mockPubSub, mockChannelManager, mockLogger)

	mockSub := pubsubmock.NewMockSubscription(ctrl)
	mockSub.EXPECT().Receive(gomock.Any()).Return(nil, context.Canceled)
	mockSub.EXPECT().Close()

	mockPubSub.EXPECT().
		Subscribe(gomock.Any(), "room").
		Return(mockSub, nil).Times(1)

	subscriber.Start("room")
	subscriber.Start("room")

	time.Sleep(20 * time.Millisecond)
	subscriber.Stop()
}

func TestRedisSubscriber_Start_LogsFatalOnReceiveError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPubSub := pubsubmock.NewMockPubSub(ctrl)
	mockSubscription := pubsubmock.NewMockSubscription(ctrl)
	mockChannelManager := realtimemock.NewMockChannelManager(ctrl)
	mockLogger := loggermock.NewMockLogger(ctrl)

	mockPubSub.EXPECT().
		Subscribe(gomock.Any(), "room").
		Return(mockSubscription, nil)

	mockSubscription.EXPECT().
		Receive(gomock.Any()).
		Return(nil, fmt.Errorf("receive failed"))

	mockSubscription.EXPECT().Close()

	mockLogger.EXPECT().
		Fatalw("[RedisSubscriber] receive error", "channel", "room", "error", gomock.Any())

	subscriber := redis.NewRedisSubscriber(mockPubSub, mockChannelManager, mockLogger)
	subscriber.Start("room")

	time.Sleep(10 * time.Millisecond)
	subscriber.Stop()
}
