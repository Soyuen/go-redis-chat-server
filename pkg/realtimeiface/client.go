package realtimeiface

import (
	"sync"

	"github.com/Soyuen/go-redis-chat-server/pkg/loggeriface"
	"github.com/gorilla/websocket"
)

type Client struct {
	conn      *websocket.Conn
	send      chan []byte
	logger    loggeriface.Logger
	closeOnce sync.Once
}

func NewClient(conn *websocket.Conn, logger loggeriface.Logger,
) *Client {
	return &Client{
		conn:   conn,
		logger: logger,
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
