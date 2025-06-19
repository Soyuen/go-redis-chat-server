package realtimeiface

type Connection interface {
	HandleConnection(conn WSConn, channel string, onMessage func(raw []byte) *Message, onClose func())
}
