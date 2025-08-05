package realtime

//go:generate mockgen -destination=mocks/mock_wsconn.go -package=mocks github.com/Soyuen/go-redis-chat-server/internal/application/realtime WSConn
//go:generate mockgen -destination=mocks/mock_client.go -package=mocks github.com/Soyuen/go-redis-chat-server/internal/application/realtime Client
//go:generate mockgen -destination=mocks/mock_client_factory.go -package=mocks github.com/Soyuen/go-redis-chat-server/internal/application/realtime ClientFactory
//go:generate mockgen -destination=mocks/mock_broadcaster.go -package=mocks github.com/Soyuen/go-redis-chat-server/internal/application/realtime Broadcaster
//go:generate mockgen -destination=mocks/mock_channel_manager.go -package=mocks github.com/Soyuen/go-redis-chat-server/internal/application/realtime ChannelManager
//go:generate mockgen -destination=mocks/mock_channel_event_subscriber.go -package=mocks github.com/Soyuen/go-redis-chat-server/internal/application/realtime ChannelEventSubscriber
//go:generate mockgen -destination=mocks/mock_connection.go -package=mocks github.com/Soyuen/go-redis-chat-server/internal/application/realtime Connection

// Message DTO
type Message struct {
	Channel string `json:"channel"`
	Data    []byte `json:"data"`
}

// WebSocket connection abstraction
type WSConn interface {
	WriteMessage(messageType int, data []byte) error
	ReadMessage() (messageType int, data []byte, err error)
	Close() error
}

// Client abstraction
type Client interface {
	Send(message []byte)
	ReadPump(onMessage func([]byte)) error
	WritePump()
	Close()
}

// Client factory
type ClientFactory interface {
	New(conn WSConn) Client
}

// Channel broadcaster
type Broadcaster interface {
	Register(client Client)
	Unregister(client Client)
	Broadcast(message []byte)
	CloseAllClients()
}

// Channel manager
type ChannelManager interface {
	GetOrCreateChannel(channel string) Broadcaster
	Broadcast(msg Message)
	CloseChannel(channel string)
	CloseAllChannels()
}

// Channel event subscriber
type ChannelEventSubscriber interface {
	Start(channel string)
	Stop()
}

// Connection handler
type Connection interface {
	HandleConnection(conn WSConn, channel string, onMessage func(raw []byte) *Message, onClose func())
}
