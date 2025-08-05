package errors

import "errors"

type ErrorResponse struct {
	Code int `json:"code"`
}

var (
	ErrInvalidRequestBody     = errors.New("invalid request body")
	ErrWebSocketUpgradeFailed = errors.New("websocket upgrade failed")
	ErrChannelCreationFailed  = errors.New("create channel failed")
	ErrCreateRoom             = errors.New("create room failed")
	ErrAddUserToRoom          = errors.New("add user to room failed")
	ErrBroadcastSystemMessage = errors.New("broadcast system message failed")
)
