package realtimeiface

type Message struct {
	Channel string `json:"channel"`
	Data    []byte `json:"data"`
}
