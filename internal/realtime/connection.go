package realtime

import (
	"github.com/Soyuen/go-redis-chat-server/pkg/loggeriface"
	"github.com/Soyuen/go-redis-chat-server/pkg/realtimeiface"
	"github.com/gorilla/websocket"
)

type Connection struct {
	manager *ChannelManager
	logger  loggeriface.Logger
}

func NewConnection(m *ChannelManager, logger loggeriface.Logger) *Connection {
	return &Connection{
		manager: m,
		logger:  logger}
}

// Ensure the interface is implemented
var _ realtimeiface.Connection = (*Connection)(nil)

func (h *Connection) HandleConnection(conn *websocket.Conn, channel string) {
	client := realtimeiface.NewClient(conn, h.logger)

	b := h.manager.GetOrCreateChannel(channel)
	b.Register(client)

	go h.handleWrite(client)
	h.handleRead(client, channel, b)
}

func (h *Connection) handleWrite(client *realtimeiface.Client) {
	defer func() {
		if r := recover(); r != nil {
			h.logger.Errorw("Recovered from panic in WritePump", "error", r)
		}
	}()
	client.WritePump()
}

func (h *Connection) handleRead(client *realtimeiface.Client, channel string, b *realtimeiface.Broadcaster) {
	defer func() {
		b.Unregister(client)
		client.Close()
	}()

	client.ReadPump(func(message []byte) {
		msg := realtimeiface.Message{
			Channel: channel,
			Data:    message,
		}
		h.manager.Broadcast(msg)
	})
}
