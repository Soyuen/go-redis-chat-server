package handler

import (
	"context"
	"errors"
	"net/http"

	appchat "github.com/Soyuen/go-redis-chat-server/internal/application/chat"
	"github.com/Soyuen/go-redis-chat-server/internal/application/realtime"
	apperr "github.com/Soyuen/go-redis-chat-server/internal/errors"
	infrarealtime "github.com/Soyuen/go-redis-chat-server/internal/infrastructure/realtime"
	"github.com/Soyuen/go-redis-chat-server/internal/presenter"
	"github.com/Soyuen/go-redis-chat-server/pkg/loggeriface"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ChatHandler struct {
	manager      realtime.ChannelManager
	connection   realtime.Connection
	chatService  appchat.ChatService
	upgrader     websocket.Upgrader
	upgraderFunc func(w http.ResponseWriter, r *http.Request) (realtime.WSConn, error)
	presenter    presenter.MessagePresenterInterface
	logger       loggeriface.Logger
}

func NewChatHandler(manager realtime.ChannelManager, connection realtime.Connection,
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
	h.upgraderFunc = func(w http.ResponseWriter, r *http.Request) (realtime.WSConn, error) {
		conn, err := h.upgrader.Upgrade(w, r, nil)
		if err != nil {
			return nil, err
		}
		return infrarealtime.NewWSConnWrapper(conn), nil
	}
	return h
}

func (h *ChatHandler) JoinChannel(c *gin.Context) {
	ctx := context.Background()
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

	if err := h.chatService.JoinChannel(ctx, channel, nickname); err != nil {
		conn.Close()
		switch {
		case errors.Is(err, apperr.ErrCreateRoom):
			h.logger.Errorw("create room failed", "err", err)
			c.JSON(http.StatusInternalServerError, apperr.ErrorResponse{
				Code: apperr.ErrCodeChannelCreationFailed,
			})
		case errors.Is(err, apperr.ErrAddUserToRoom):
			h.logger.Errorw("add user to room failed", "err", err)
			c.JSON(http.StatusInternalServerError, apperr.ErrorResponse{
				Code: apperr.ErrCodeChannelJoinFailed,
			})
		case errors.Is(err, apperr.ErrBroadcastSystemMessage):
			h.logger.Errorw("broadcast system message failed", "err", err)
			c.JSON(http.StatusInternalServerError, apperr.ErrorResponse{
				Code: apperr.ErrCodeChannelJoinFailed,
			})
		default:
			h.logger.Errorw("unknown join channel error", "err", err)
			c.JSON(http.StatusInternalServerError, apperr.ErrorResponse{
				Code: apperr.ErrCodeChannelJoinFailed,
			})
		}
		return
	}

	h.connection.HandleConnection(
		conn,
		channel,
		h.messageHandler(channel, nickname),
		h.leaveHandler(channel, nickname),
	)

}

func (h *ChatHandler) messageHandler(channel, nickname string) func(raw []byte) *realtime.Message {
	return func(raw []byte) *realtime.Message {
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
		ctx := context.Background()
		if err := h.chatService.BroadcastSystemMessage(ctx, channel, nickname, "left"); err != nil {
			h.logger.Warnw("failed to announce leave", "err", err)
		}
	}
}

// for testing
func (h *ChatHandler) SetUpgraderFunc(f func(w http.ResponseWriter, r *http.Request) (realtime.WSConn, error)) {
	h.upgraderFunc = f
}
