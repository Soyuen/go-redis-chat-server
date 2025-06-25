package handler

import (
	"net/http"

	appchat "github.com/Soyuen/go-redis-chat-server/internal/application/chat"
	apperr "github.com/Soyuen/go-redis-chat-server/internal/errors"
	"github.com/Soyuen/go-redis-chat-server/internal/infrastructure/realtime"
	"github.com/Soyuen/go-redis-chat-server/internal/presenter"
	"github.com/Soyuen/go-redis-chat-server/pkg/loggeriface"
	"github.com/Soyuen/go-redis-chat-server/pkg/realtimeiface"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ChatHandler struct {
	manager      realtimeiface.ChannelManager
	connection   realtimeiface.Connection
	chatService  appchat.ChatService
	upgrader     websocket.Upgrader
	upgraderFunc func(w http.ResponseWriter, r *http.Request) (realtimeiface.WSConn, error)
	presenter    presenter.MessagePresenterInterface
	logger       loggeriface.Logger
}

func NewChatHandler(manager realtimeiface.ChannelManager, connection realtimeiface.Connection,
	chatService appchat.ChatService, presenter presenter.MessagePresenterInterface, logger loggeriface.Logger) *ChatHandler {
	h := &ChatHandler{
		manager:     manager,
		connection:  connection,
		chatService: chatService,
		logger:      logger,
		presenter:   presenter,
		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
	h.upgraderFunc = func(w http.ResponseWriter, r *http.Request) (realtimeiface.WSConn, error) {
		conn, err := h.upgrader.Upgrade(w, r, nil)
		if err != nil {
			return nil, err
		}
		return realtime.NewWSConnWrapper(conn), nil
	}
	return h
}

func (h *ChatHandler) JoinChannel(c *gin.Context) {
	channel := c.Query("channel")
	if channel == "" {
		c.JSON(http.StatusBadRequest, apperr.ErrorResponse{
			Code: apperr.ErrCodeInvalidRequestBody,
		})
		return
	}

	nickname := c.Query("nickname")
	if nickname == "" || nickname == "System" {
		c.JSON(http.StatusBadRequest, apperr.ErrorResponse{
			Code: apperr.ErrCodeInvalidRequestBody,
		})
		return
	}

	conn, err := h.upgraderFunc(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, apperr.ErrorResponse{
			Code: apperr.ErrCodeWebSocketUpgradeFailed,
		})
		return
	}

	if err := h.chatService.CreateRoom(channel); err != nil {
		conn.Close()
		c.JSON(http.StatusInternalServerError, apperr.ErrorResponse{
			Code: apperr.ErrCodeChannelCreationFailed,
		})
		return
	}

	if err := h.chatService.BroadcastSystemMessage(channel, nickname, "joined"); err != nil {
		h.logger.Warnw("failed to announce join", "err", err)
	}

	h.connection.HandleConnection(
		conn,
		channel,
		h.messageHandler(channel, nickname),
		h.leaveHandler(channel, nickname),
	)

}

func (h *ChatHandler) messageHandler(channel, nickname string) func(raw []byte) *realtimeiface.Message {
	return func(raw []byte) *realtimeiface.Message {
		dmsg, err := h.chatService.ProcessIncoming(raw, nickname, channel)
		if err != nil {
			h.logger.Warnw("failed to parse message", "err", err)
			return nil
		}
		return h.presenter.Format(dmsg)
	}
}

func (h *ChatHandler) leaveHandler(channel, nickname string) func() {
	return func() {
		if err := h.chatService.BroadcastSystemMessage(channel, nickname, "left"); err != nil {
			h.logger.Warnw("failed to announce leave", "err", err)
		}
	}
}

// for testing
func (h *ChatHandler) SetUpgraderFunc(f func(w http.ResponseWriter, r *http.Request) (realtimeiface.WSConn, error)) {
	h.upgraderFunc = f
}
