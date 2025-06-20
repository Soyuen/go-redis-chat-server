package realtime

import (
	"testing"

	"github.com/Soyuen/go-redis-chat-server/internal/infrastructure/realtime/mocks"
	"github.com/Soyuen/go-redis-chat-server/pkg/realtimeiface"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestChannelManager_GetOrCreateChannel_ReturnsExisting(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	b := mocks.NewMockBroadcaster(ctrl)
	cm := NewChannelManager()
	cm.channels["test"] = b

	got := cm.GetOrCreateChannel("test")
	assert.Equal(t, b, got)
}

func TestChannelManager_GetOrCreateChannel_CreatesNew(t *testing.T) {
	cm := NewChannelManager()

	got := cm.GetOrCreateChannel("new")
	assert.NotNil(t, got)
	assert.Contains(t, cm.channels, "new")
}

func TestChannelManager_Broadcast(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	b := mocks.NewMockBroadcaster(ctrl)
	cm := NewChannelManager()
	cm.channels["room"] = b

	msg := realtimeiface.Message{Channel: "room", Data: []byte("hi")}

	b.EXPECT().Broadcast(gomock.Any()).Times(1)

	cm.Broadcast(msg)
}

func TestChannelManager_CloseChannel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	b := mocks.NewMockBroadcaster(ctrl)
	cm := NewChannelManager()
	cm.channels["test"] = b

	b.EXPECT().CloseAllClients().Times(1)

	cm.CloseChannel("test")
	assert.NotContains(t, cm.channels, "test")
}

func TestChannelManager_CloseAllChannels(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	b1 := mocks.NewMockBroadcaster(ctrl)
	b2 := mocks.NewMockBroadcaster(ctrl)

	cm := NewChannelManager()
	cm.channels["a"] = b1
	cm.channels["b"] = b2

	b1.EXPECT().CloseAllClients().Times(1)
	b2.EXPECT().CloseAllClients().Times(1)

	cm.CloseAllChannels()
	assert.Empty(t, cm.channels)
}
