package realtime

import (
	"testing"

	"github.com/Soyuen/go-redis-chat-server/internal/application/realtime"
	"github.com/Soyuen/go-redis-chat-server/internal/infrastructure/realtime/mocks"
	"github.com/Soyuen/go-redis-chat-server/internal/testhelper"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestChannelManager_GetOrCreateChannel_ReturnsExisting(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	b := mocks.NewMockBroadcaster(ctrl)
	cm := NewChannelManager()
	cm.channels[testhelper.ChannelTest] = b

	got := cm.GetOrCreateChannel(testhelper.ChannelTest)
	assert.Equal(t, b, got)
}

func TestChannelManager_GetOrCreateChannel_CreatesNew(t *testing.T) {
	cm := NewChannelManager()

	got := cm.GetOrCreateChannel(testhelper.ChannelTest)
	assert.NotNil(t, got)
	assert.Contains(t, cm.channels, testhelper.ChannelTest)
}

func TestChannelManager_Broadcast(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	b := mocks.NewMockBroadcaster(ctrl)
	cm := NewChannelManager()
	cm.channels[testhelper.ChannelTest] = b

	msg := realtime.Message{
		Channel: testhelper.ChannelTest,
		Data:    []byte(testhelper.MessageTest),
	}

	b.EXPECT().Broadcast(gomock.Any()).Times(1)

	cm.Broadcast(msg)
}

func TestChannelManager_CloseChannel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	b := mocks.NewMockBroadcaster(ctrl)
	cm := NewChannelManager()
	cm.channels[testhelper.ChannelTest] = b

	b.EXPECT().CloseAllClients().Times(1)

	cm.CloseChannel(testhelper.ChannelTest)
	assert.NotContains(t, cm.channels, testhelper.ChannelTest)
}

func TestChannelManager_CloseAllChannels(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	b1 := mocks.NewMockBroadcaster(ctrl)
	b2 := mocks.NewMockBroadcaster(ctrl)

	cm := NewChannelManager()

	cm.channels[testhelper.ChannelTest] = b1
	cm.channels[testhelper.ChannelGeneral] = b2

	b1.EXPECT().CloseAllClients().Times(1)
	b2.EXPECT().CloseAllClients().Times(1)

	cm.CloseAllChannels()

	assert.Empty(t, cm.channels)
}
