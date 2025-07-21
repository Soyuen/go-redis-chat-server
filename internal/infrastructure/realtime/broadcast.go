package realtime

import (
	"sync"

	"github.com/Soyuen/go-redis-chat-server/internal/application/realtime"
)

type BroadcasterImpl struct {
	mu      sync.RWMutex
	clients map[realtime.Client]bool
}

func NewBroadcaster() realtime.Broadcaster {
	return &BroadcasterImpl{
		clients: make(map[realtime.Client]bool),
	}
}

func (b *BroadcasterImpl) Register(client realtime.Client) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.clients[client] = true
}

func (b *BroadcasterImpl) Unregister(client realtime.Client) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if _, exists := b.clients[client]; exists {
		delete(b.clients, client)
		client.Close()
	}
}

func (b *BroadcasterImpl) Broadcast(message []byte) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for client := range b.clients {
		client.Send(message)
	}
}

func (b *BroadcasterImpl) CloseAllClients() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for client := range b.clients {
		client.Close()
		delete(b.clients, client)
	}
}
