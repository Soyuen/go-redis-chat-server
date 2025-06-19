package realtime

import (
	"sync"

	"github.com/Soyuen/go-redis-chat-server/pkg/loggeriface"
	"github.com/Soyuen/go-redis-chat-server/pkg/realtimeiface"
	"github.com/gorilla/websocket"
)

var _ realtimeiface.Client = (*Client)(nil)

type Client struct {
	conn      realtimeiface.WSConn
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

func (f *ClientFactory) New(conn realtimeiface.WSConn) *Client {
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
		c.conn.WriteMessage(websocket.TextMessage, msg)
	}
}

func (c *Client) Close() {
	c.conn.Close()
	c.closeOnce.Do(func() {
		close(c.send)
	})
}
