package realtimeiface

import "github.com/gorilla/websocket"

type Connection interface {
	HandleConnection(conn *websocket.Conn, channel string)
}
