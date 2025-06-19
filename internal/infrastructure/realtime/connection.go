package realtime

import (
	"github.com/Soyuen/go-redis-chat-server/pkg/loggeriface"
	"github.com/Soyuen/go-redis-chat-server/pkg/realtimeiface"
)

type Connection struct {
	manager       *ChannelManager
	logger        loggeriface.Logger
	clientFactory *ClientFactory
}

func NewConnection(m *ChannelManager, logger loggeriface.Logger, clientFactory *ClientFactory) *Connection {
	return &Connection{
		manager:       m,
		logger:        logger,
		clientFactory: clientFactory,
	}
}

// 確保 Connection 有實作 interface
var _ realtimeiface.Connection = (*Connection)(nil)

func (h *Connection) HandleConnection(
	conn realtimeiface.WSConn,
	channel string,
	onMessage func(raw []byte) *realtimeiface.Message,
	onClose func(),
) {
	client := h.clientFactory.New(conn)

	b := h.manager.GetOrCreateChannel(channel)
	b.Register(client)

	go h.handleWrite(client)
	h.handleRead(client, b, onMessage, onClose)
}

func (h *Connection) handleWrite(client realtimeiface.Client) {
	defer func() {
		if r := recover(); r != nil {
			h.logger.Errorw("Recovered from panic in WritePump", "error", r)
		}
	}()
	client.WritePump()
}

func (h *Connection) handleRead(
	client realtimeiface.Client,
	b realtimeiface.Broadcaster,
	onMessage func(raw []byte) *realtimeiface.Message,
	onClose func(),
) {
	defer func() {
		b.Unregister(client)
		client.Close()
		if onClose != nil {
			onClose()
		}
	}()

	err := client.ReadPump(func(message []byte) {
		processed := onMessage(message)
		if processed != nil {
			h.manager.Broadcast(*processed)
		}
	})
	if err != nil {
		return
	}
}
