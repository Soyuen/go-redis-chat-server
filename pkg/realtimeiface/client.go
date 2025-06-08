package realtimeiface

import (
	"github.com/Soyuen/go-redis-chat-server/pkg/loggeriface"
	"github.com/gorilla/websocket"
)

type Client struct {
	conn   *websocket.Conn
	send   chan []byte
	logger loggeriface.Logger
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

func (c *Client) ReadPump(onMessage func([]byte)) {
	defer func() {
		c.logger.Infow("ReadPump closing for client")
		c.Close()
	}()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Errorw("Unexpected close", "error", err)
			} else {
				c.logger.Infow("Normal disconnection", "error", err)
			}
			break
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
	close(c.send)
}
