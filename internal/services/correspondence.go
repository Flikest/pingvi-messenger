package services

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/Flikest/PingviMessenger/internal/entity"
	"github.com/Flikest/PingviMessenger/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func (s Service) Сorrespondence(ctx *gin.Context) {
	w, r := ctx.Writer, ctx.Request

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Info("failed connect to ws", err)
		return
	}
	defer conn.Close()

	var token string = ctx.GetHeader("pinguiJWT")

	pyload, err := jwt.JwtPayloadFromRequest(token)
	if err != nil {
		ctx.Redirect(302, "https://web-pingui/login/")
	}

	chat_ID := ctx.Query("chat_id")

	for {
		var msg RequestMessage

		err := conn.ReadJSON(msg)
		if err != nil {
			slog.Debug("ugh i don't want to accept this message", err)
			return
		}

		switch msg.Operation {
		case "delete":
			status, err := s.Storage.DelelteMessage(msg.Message.Chat_ID, msg.Message.Message_ID)
			if err != nil {
				ctx.JSON(status, "the message was not updated!")
			} else {
				ctx.JSON(status, "message updated!")
			}
		case "update":
			body := entity.Message{
				Chat_ID:     chat_ID,
				Sender_ID:   pyload,
				Content:     []byte(msg.Message.Content),
				SendingTime: time.Now(),
			}
			status, err := s.Storage.UpdateMessage(body)
			if err != nil {
				ctx.JSON(status, "the message was not updated!")
			} else {
				ctx.JSON(status, "message updated!")
			}

		default:

			s.Storage.AddMesage(msg.Message)

			err := conn.WriteJSON(msg)
			if err != nil {
				ctx.JSON(http.StatusOK, "the message was not sent")
			}

		}

	}
}
