package chat

type MessageProcessor interface {
	ProcessIncoming(raw []byte, username string, channel string) (*Message, error)
}

type Message struct {
	Channel string
	Sender  string
	Content string
}
