package realtime

import (
	"encoding/json"
	"sync"

	"github.com/Soyuen/go-redis-chat-server/pkg/realtimeiface"
)

type ChannelManager struct {
	channels map[string]*realtimeiface.Broadcaster
	mu       sync.RWMutex
}

func NewChannelManager() *ChannelManager {
	return &ChannelManager{
		channels: make(map[string]*realtimeiface.Broadcaster),
	}
}

func (cm *ChannelManager) GetOrCreateChannel(channel string) *realtimeiface.Broadcaster {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if b, ok := cm.channels[channel]; ok {
		return b
	}

	b := realtimeiface.NewBroadcaster()
	cm.channels[channel] = b
	return b
}

func (cm *ChannelManager) Broadcast(msg realtimeiface.Message) {
	channel := msg.Channel
	data, err := json.Marshal(msg)
	if err != nil {
		return // TODO log.Warn
	}

	broadcaster := cm.GetOrCreateChannel(channel)
	broadcaster.Broadcast(data)
}
