package realtime

import (
	"sync"

	"github.com/Soyuen/go-redis-chat-server/internal/application/realtime"
	"github.com/Soyuen/go-redis-chat-server/pkg/loggeriface"
	"github.com/gorilla/websocket"
)

var _ realtime.Client = (*Client)(nil)

type Client struct {
	conn      realtime.WSConn
	send      chan []byte
	logger    loggeriface.Logger
	closeOnce sync.Once
}

type ClientFactory struct {
	logger loggeriface.Logger
}

func NewClientFactory(logger loggeriface.Logger) *ClientFactory {
	return &ClientFactory{
		logger: logger,
	}
}

func (f *ClientFactory) New(conn realtime.WSConn) realtime.Client {
	return &Client{
		conn:   conn,
		logger: f.logger,
		send:   make(chan []byte, 256),
	}
}

func (c *Client) Send(message []byte) {
	select {
	case c.send <- message:
	default:
		c.Close()
	}
}

func (c *Client) ReadPump(onMessage func([]byte)) error {
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Errorw("Unexpected close", "error", err)
			} else {
				c.logger.Infow("Normal disconnection", "error", err)
			}
			return err
		}
		onMessage(msg)
	}
}

func (c *Client) WritePump() {
	for msg := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			c.logger.Errorw("write message failed", "error", err)
			_ = c.conn.Close()
			break
		}
	}
}

func (c *Client) Close() {
	c.conn.Close()
	c.closeOnce.Do(func() {
		close(c.send)
	})
}
