package realtime

import (
	"testing"

	"github.com/Soyuen/go-redis-chat-server/internal/infrastructure/realtime/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestBroadcaster_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockClient(ctrl)
	b := NewBroadcaster().(*BroadcasterImpl)

	b.Register(mockClient)
	assert.Contains(t, b.clients, mockClient)
}

func TestBroadcaster_Unregister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := mocks.NewMockClient(ctrl)
	b := NewBroadcaster().(*BroadcasterImpl)
	b.Register(mockClient)

	mockClient.EXPECT().Close()
	b.Unregister(mockClient)
	assert.NotContains(t, b.clients, mockClient)
}

func TestBroadcaster_Broadcast(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// creat two mock clients
	mockClient1 := mocks.NewMockClient(ctrl)
	mockClient2 := mocks.NewMockClient(ctrl)

	b := NewBroadcaster().(*BroadcasterImpl)

	// Register mock clients
	b.Register(mockClient1)
	b.Register(mockClient2)

	// Expect both clients to receive the same message
	message := []byte("test message")
	mockClient1.EXPECT().Send(message).Times(1)
	mockClient2.EXPECT().Send(message).Times(1)

	b.Broadcast(message)
}

func TestBroadcaster_CloseAllClients(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient1 := mocks.NewMockClient(ctrl)
	mockClient2 := mocks.NewMockClient(ctrl)

	b := NewBroadcaster().(*BroadcasterImpl)
	b.Register(mockClient1)
	b.Register(mockClient2)

	mockClient1.EXPECT().Close()
	mockClient2.EXPECT().Close()

	b.CloseAllClients()
	assert.Empty(t, b.clients)
}

func TestBroadcaster_Broadcast_NoClients(t *testing.T) {
	b := NewBroadcaster().(*BroadcasterImpl)
	//Calling with no clients should not panic.
	b.Broadcast([]byte("hello"))
}

func TestBroadcaster_Unregister_NonExistentClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockClient := mocks.NewMockClient(ctrl)
	b := NewBroadcaster().(*BroadcasterImpl)

	//Unregistering a non-existent client should be safe;
	// Close should not be called, so no expectation is set.
	b.Unregister(mockClient)
}
