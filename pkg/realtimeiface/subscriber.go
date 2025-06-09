// pkg/realtimeiface/subscriber.go
package realtimeiface

type ChannelEventSubscriber interface {
	Start(channel string)
	Stop()
}
