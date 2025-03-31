package services

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Flikest/PingviMessenger/internal/entity"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/websocket"
	"github.com/savsgio/gotils/uuid"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  2048,
	WriteBufferSize: 2048,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type RequestMessage struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	Operation string `json:"operation"`
}

var chanDataFromTheStartPage chan []entity.Chat
var chanCounterMessage chan int
var chanAllMessageFromChat chan []entity.Message

func jwtPayloadFromRequest(tokenString string) (string, error) {

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return "", err
	}

	result, err := json.Marshal(claims["sub"])
	if err != nil {
		slog.Info("Error: %s", err)
		return "", err
	}
	return string(result), nil
}

func sendMessageInchat(mesage_id uuid.UUID)

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

	var token string = ctx.GetHeader("pinguiJWT")

	pyload, err := jwtPayloadFromRequest(token)
	if err != nil {
		ctx.Redirect(302, "https://web-pingui/login/")
	}

	chat_ID := ctx.Query("chat_id")

	go s.Storage.CountMessages(chat_ID, chanCounterMessage)

	for {
		var msg RequestMessage
		err := conn.ReadJSON(msg)

		if err != nil {
			slog.Debug("ugh i don't want to accept this message", err)
			return
		}

		if err := conn.WriteJSON(RequestMessage{}); err != nil {
			slog.Info("message not sent", err)
			return
		}

		switch msg.Operation {
		case "delete":
			status, err := s.Storage.DelelteMessage(msg.ID)
			if err != nil {
				ctx.JSON(status, "the message was not updated!")
			} else {
				ctx.JSON(status, "message updated!")
			}
		case "update":

			body := entity.Message{
				Chat_ID:     chat_ID,
				Message_ID:  <-chanCounterMessage,
				Sender_ID:   pyload,
				Content:     []byte(msg.Content),
				SendingTime: time.Now(),
			}
			status, err := s.Storage.UpdateMessage(body)
			if err != nil {
				ctx.JSON(status, "the message was not updated!")
			} else {
				ctx.JSON(status, "message updated!")
			}

		default:
			messege := entity.Message{
				Chat_ID:     chat_ID,
				Message_ID:  <-chanCounterMessage + 1,
				Sender_ID:   pyload,
				Content:     []byte(msg.Content),
				SendingTime: time.Now(),
			}

			s.Storage.AddMesage(messege)

			err := conn.WriteJSON(messege)
			if err != nil {
				ctx.JSON(http.StatusOK, "the message was not sent")
			}

		}

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

	_, err := s.Storage.CreateChat(body.ID, body)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, "failed to create chat")
	} else {
		ctx.JSON(http.StatusOK, "chat created")
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
	chat_id := ctx.Param("chat_id")

	result, err := s.Storage.GetMessage(chat_id, message_id)
	if err != nil {
		ctx.JSON(404, "message not found!")
	} else {
		ctx.JSON(http.StatusOK, result)
	}
}

// @Router       users/{chat_id} [post]
func (s Service) AddUserChat(ctx *gin.Context) {

	body := entity.Participant{}
	ctx.BindJSON(&body)

	status, err := s.Storage.AddUser(body.Chat_ID, body.User_ID)
	if err != nil {
		ctx.JSON(status, "Failed to add user")
	} else {
		reply := fmt.Sprintf("user %s added to chat %s", body.Chat_ID, body.User_ID)
		ctx.JSON(status, reply)
	}
}

// @Router       users/{chat_id} [delete]
func (s Service) DropUserFromChat(ctx *gin.Context) {
	body := entity.Participant{}
	ctx.BindJSON(&body)

	status, err := s.Storage.DropUserFromChat(body.User_ID, body.Chat_ID)
	if err != nil {
		ctx.JSON(status, "failed to kick out user")
	} else {
		ctx.JSON(status, "the user will be kicked")
	}
}

func (s Service) DataFromTheStartPage(ctx *gin.Context) {
	var token string = ctx.GetHeader("pinguiJWT")

	pyload, err := jwtPayloadFromRequest(token)
	if err != nil {
		ctx.Redirect(302, "https://web-pingui/login/")
	}

	go s.Storage.DataFromTheStartPage(pyload, chanDataFromTheStartPage)
}
