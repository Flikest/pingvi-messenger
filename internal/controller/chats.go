package services

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Flikest/PingviMessenger/internal/entity"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

// @Accept
// @Router       ws://localhost:9000/chats/messenger/{id} [get]
func (s Service) Сorrespondence(ctx *gin.Context) {
	w, r := ctx.Writer, ctx.Request
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Info("failed connect to ws", err)
		return
	}
	defer conn.Close()

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			slog.Debug("ugh I don't want to accept this message", err)
			return
		}

		if err := conn.WriteMessage(messageType, p); err != nil {
			slog.Info("message not sent", err)
			return
		}

		chat_id := r.PathValue("id")
		messege := entity.Messege{
			Content:     p,
			SendingTime: time.Now(),
		}
		s.Storage.AddMesage(chat_id, messege)
	}
}

// @Router       chats/{chat_name} [get]
func (s Service) GetChat(ctx *gin.Context) {
	chat_name := ctx.Param("chat_name")
	result, err := s.Storage.GetChat(chat_name)
	if err != nil {
		ctx.JSON(http.StatusOK, err)
	} else {
		ctx.JSON(http.StatusOK, result)
	}
}

// @Router       chats/ [post]
func (s Service) CreateChat(ctx *gin.Context) {
	body := entity.Chat{}
	ctx.BindJSON(&body)

	status, err := s.Storage.CreateChat(body.Creator_id)
	if err != nil {
		ctx.JSON(status, "failed to create chat")
	} else {
		ctx.JSON(status, "chat created")
	}

}

// @Router       chats/ [put]
func (s Service) UpdateChat(ctx *gin.Context) {
	body := entity.Chat{}
	ctx.BindJSON(body)

	status, err := s.Storage.UpdateChat(body)
	if err != nil {
		ctx.JSON(status, "failed to update chat")
	} else {
		ctx.JSON(status, "chat updated")
	}
}

// @Router       chats/{chat_id} [delete]
func (s Service) DeleteChat(ctx *gin.Context) {
	chat_id := ctx.Param("chat_id")

	status, err := s.Storage.DeleteChat(chat_id)
	if err != nil {
		ctx.JSON(status, "failed to delete chat")
	} else {
		ctx.JSON(status, "chat deleted!")
	}
}

// @Router       message/{messege_Id} [get]
func (s Service) GetMessage(ctx *gin.Context) {
	message_id := ctx.Param("message_id")

	result, err := s.Storage.GetMessage(message_id)
	if err != nil {
		ctx.JSON(404, "message not found!")
	} else {
		ctx.JSON(http.StatusOK, result)
	}
}

// @Router       message/ [put]
func (s Service) UpdateMessage(ctx *gin.Context) {
	body := entity.Messege{}
	ctx.BindJSON(&body)

	status, err := s.Storage.UpdateMessage(body.Message_id, body)
	if err != nil {
		ctx.JSON(status, "the message was not updated!")
	} else {
		ctx.JSON(status, "message updated!")
	}
}

// @Router       message/{messege_Id} [delete]
func (s Service) DelelteMessage(ctx *gin.Context) {
	message_id := ctx.Param("message_id")

	status, err := s.Storage.DelelteMessage(message_id)
	if err != nil {
		ctx.JSON(status, "the message was not updated!")
	} else {
		ctx.JSON(status, "message updated!")
	}
}

// @Router       users/{chat_id} [post]
func (s Service) AddUserChat(ctx *gin.Context) {
	chat_id := ctx.Param("chat_id")

	body := entity.User{}
	ctx.BindJSON(&body)

	status, err := s.Storage.AddUser(chat_id, body.ID)
	if err != nil {
		ctx.JSON(status, "Failed to add user")
	} else {
		reply := fmt.Sprintf("user %s added to chat %s", body.ID, chat_id)
		ctx.JSON(status, reply)
	}
}

// @Router       users/{chat_id} [delete]
func (s Service) DropUserFromChat(ctx *gin.Context) {
	body := entity.Participant{}
	ctx.BindJSON(&body)

	chat_id := ctx.Param("chat_id")

	status, err := s.Storage.DropUserFromChat(body, chat_id)
	if err != nil {
		ctx.JSON(status, "failed to kick out user")
	} else {
		ctx.JSON(status, "the user will be kicked")
	}
}
