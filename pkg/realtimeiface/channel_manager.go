package realtimeiface

type ChannelManager interface {
	GetOrCreateChannel(channel string) Broadcaster
	Broadcast(msg Message)
	CloseChannel(channel string)
	CloseAllChannels()
}
