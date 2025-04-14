package services

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/Flikest/PingviMessenger/internal/entity"
	"github.com/Flikest/PingviMessenger/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type RequestMessage struct {
	Operation string `json:"operation"`
	Message   entity.Message
}

var (
	ChanMessageError = make(chan error)
	ChanMessage      = make(chan entity.Message)
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func sendMessage(ch chan entity.Message, conn *websocket.Conn) {
	if err := conn.WriteJSON(<-ch); err != nil {
		slog.Info("message sending error: ", err)
	}
}

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
		if err := conn.ReadJSON(msg); err != nil {
			slog.Info("message reading error: ", err)
		}

		ChanMessage <- msg.Message

		switch msg.Operation {
		case "delete":

			go s.Storage.DeleteMessage(ChanMessageError, msg.Message.Chat_ID, msg.Message.Message_ID)
			if <-ChanMessageError != nil {
				slog.Info("error deleting message: ", <-ChanMessageError)
			} else {
				go sendMessage(ChanMessage, conn)
			}

		case "update":
			body := entity.Message{
				Chat_ID:     chat_ID,
				Sender_ID:   pyload,
				Content:     []byte(msg.Message.Content),
				SendingTime: time.Now(),
			}
			go s.Storage.UpdateMessage(ChanMessageError, body)
			if <-ChanMessageError != nil {
				slog.Info("error updating message: ", <-ChanMessageError)
			} else {
				go sendMessage(ChanMessage, conn)
			}

		default:
			go s.Storage.AddMesage(ChanMessageError, msg.Message)
			go sendMessage(ChanMessage, conn)
			if <-ChanMessageError != nil {
				slog.Info("error added message: ", <-ChanMessageError)
			} else {
				go sendMessage(ChanMessage, conn)
			}

		}

	}
}
