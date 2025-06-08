package realtimeiface

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		conn: conn,
		send: make(chan []byte, 256),
	}
}

func (c *Client) Send(message []byte) {
	select {
	case c.send <- message:
	default:
		// buffer 滿了，考慮斷開連線
		c.Close()
	}
}

func (c *Client) ReadPump(onMessage func([]byte)) {
	defer c.Close()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
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
