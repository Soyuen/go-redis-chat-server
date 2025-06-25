package realtimeiface

type Client interface {
	Send(message []byte)
	ReadPump(onMessage func([]byte)) error
	WritePump()
	Close()
}

type ClientFactory interface {
	New(conn WSConn) Client
}
