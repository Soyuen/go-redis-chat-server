package realtimeiface

type ChatChannelManager interface {
	GetOrCreateChannel(channel string) *Broadcaster
	Broadcast(msg Message)
}
