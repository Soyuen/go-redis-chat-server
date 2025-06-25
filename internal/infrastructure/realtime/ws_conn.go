package realtime

import (
	"github.com/Soyuen/go-redis-chat-server/pkg/realtimeiface"
	"github.com/gorilla/websocket"
)

type WSConnWrapper struct {
	conn *websocket.Conn
}

func NewWSConnWrapper(conn *websocket.Conn) realtimeiface.WSConn {
	return &WSConnWrapper{
		conn: conn,
	}
}

func (w *WSConnWrapper) WriteMessage(mt int, data []byte) error {
	return w.conn.WriteMessage(mt, data)
}

func (w *WSConnWrapper) ReadMessage() (int, []byte, error) {
	return w.conn.ReadMessage()
}

func (w *WSConnWrapper) Close() error {
	return w.conn.Close()
}
