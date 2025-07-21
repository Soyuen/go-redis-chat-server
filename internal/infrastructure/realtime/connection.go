package realtime

import (
	"github.com/Soyuen/go-redis-chat-server/internal/application/realtime"
	"github.com/Soyuen/go-redis-chat-server/pkg/loggeriface"
)

type Connection struct {
	manager       realtime.ChannelManager
	logger        loggeriface.Logger
	clientFactory realtime.ClientFactory
}

func NewConnection(m realtime.ChannelManager, logger loggeriface.Logger, clientFactory realtime.ClientFactory) *Connection {
	return &Connection{
		manager:       m,
		logger:        logger,
		clientFactory: clientFactory,
	}
}

// 確保 Connection 有實作 interface
var _ realtime.Connection = (*Connection)(nil)

func (h *Connection) HandleConnection(
	conn realtime.WSConn,
	channel string,
	onMessage func(raw []byte) *realtime.Message,
	onClose func(),
) {
	client := h.clientFactory.New(conn)

	b := h.manager.GetOrCreateChannel(channel)
	b.Register(client)

	go h.handleWrite(client)
	h.handleRead(client, b, onMessage, onClose)
}

func (h *Connection) handleWrite(client realtime.Client) {
	defer func() {
		if r := recover(); r != nil {
			h.logger.Errorw("Recovered from panic in WritePump", "error", r)
		}
	}()
	client.WritePump()
}

func (h *Connection) handleRead(
	client realtime.Client,
	b realtime.Broadcaster,
	onMessage func(raw []byte) *realtime.Message,
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
