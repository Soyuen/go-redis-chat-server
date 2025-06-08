package realtime

import (
	"github.com/Soyuen/go-redis-chat-server/pkg/realtimeiface"
	"github.com/gorilla/websocket"
)

type Connection struct {
	manager *ChannelManager
}

func NewConnection(m *ChannelManager) *Connection {
	return &Connection{manager: m}
}

// Ensure the interface is implemented
var _ realtimeiface.Connection = (*Connection)(nil)

func (h *Connection) HandleConnection(conn *websocket.Conn, channel string) {
	client := realtimeiface.NewClient(conn)

	b := h.manager.GetOrCreateChannel(channel)
	b.Register(client)
	defer b.Unregister(client)

	// Start the write loop (asynchronous)
	go client.WritePump()

	// Start the read loop (executed synchronously)
	client.ReadPump(func(message []byte) {
		// Wrap the message with the channel name and broadcast it
		msg := realtimeiface.Message{
			Channel: channel,
			Data:    message,
		}
		h.manager.Broadcast(msg)
	})
}
