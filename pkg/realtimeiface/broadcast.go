package realtimeiface

import (
	"sync"
)

type Broadcaster struct {
	clients map[*Client]bool
	mu      sync.RWMutex
}

func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		clients: make(map[*Client]bool),
	}
}

func (b *Broadcaster) Register(client *Client) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.clients[client] = true
}

func (b *Broadcaster) Unregister(client *Client) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.clients, client)
	client.Close()
}

func (b *Broadcaster) Broadcast(message []byte) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	for client := range b.clients {
		client.Send(message)
	}
}
