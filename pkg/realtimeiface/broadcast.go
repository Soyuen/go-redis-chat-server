package realtimeiface

type Broadcaster interface {
	Register(client Client)
	Unregister(client Client)
	Broadcast(message []byte)
}
