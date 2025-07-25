package realtime

import (
	"encoding/json"
	"sync"

	"github.com/Soyuen/go-redis-chat-server/internal/application/realtime"
)

type ChannelManager struct {
	channels map[string]realtime.Broadcaster
	mu       sync.RWMutex
}

func NewChannelManager() *ChannelManager {
	return &ChannelManager{
		channels: make(map[string]realtime.Broadcaster),
	}
}

func (cm *ChannelManager) GetOrCreateChannel(channel string) realtime.Broadcaster {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if b, ok := cm.channels[channel]; ok {
		return b
	}

	b := NewBroadcaster()
	cm.channels[channel] = b
	return b
}

func (cm *ChannelManager) Broadcast(msg realtime.Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		// TODO: log warning
		return
	}

	broadcaster := cm.GetOrCreateChannel(msg.Channel)
	broadcaster.Broadcast(data)
}

func (cm *ChannelManager) CloseChannel(channel string) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if b, ok := cm.channels[channel]; ok {
		b.CloseAllClients()
		delete(cm.channels, channel)
	}
}

func (cm *ChannelManager) CloseAllChannels() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	for name, b := range cm.channels {
		b.CloseAllClients()
		delete(cm.channels, name)
	}
}
