package realtime

import (
	"sync"

	"github.com/Soyuen/go-redis-chat-server/pkg/realtimeiface"
)

type BroadcasterImpl struct {
	mu      sync.RWMutex
	clients map[realtimeiface.Client]bool
}

func NewBroadcaster() realtimeiface.Broadcaster {
	return &BroadcasterImpl{
		clients: make(map[realtimeiface.Client]bool),
	}
}

func (b *BroadcasterImpl) Register(client realtimeiface.Client) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.clients[client] = true
}

func (b *BroadcasterImpl) Unregister(client realtimeiface.Client) {
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
