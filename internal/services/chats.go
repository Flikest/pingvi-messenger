package services

import (
	"log/slog"
	"net/http"

	"github.com/Flikest/PingviMessenger/internal/entity"
	"github.com/Flikest/PingviMessenger/pkg/jwt"
	"github.com/gin-gonic/gin"
)

var (
	ChanErrorHandling = make(chan error)
	ChanChats         = make(chan []entity.Chat)
)

func (s Service) DataFromTheStartPage(ctx *gin.Context) {
	var token string = ctx.GetHeader("pinguiJWT")

	pyload, err := jwt.JwtPayloadFromRequest(token)
	if err != nil {
		ctx.Redirect(302, "https://web-pingui/login/")
	}

	go s.Storage.DataFromTheStartPage(pyload, ChanChats)
}

// @Accept
// @Router       ws://localhost:9000/chats/messenger/{id} [get]

// @Router       chats/{chat_name} [get]
func (s Service) GetChat(ctx *gin.Context) {
	chat_name := ctx.Param("chat_name")
	go s.Storage.GetChat(ChanChats, chat_name)
	ctx.JSON(http.StatusOK, <-ChanChats)

}

// @Router       chats/ [post]
func (s Service) CreateChat(ctx *gin.Context) {
	body := entity.Chat{}
	ctx.BindJSON(&body)

	go s.Storage.CreateChat(body.ID, body, ChanErrorHandling)
	if <-ChanErrorHandling != nil {
		slog.Info("error creating chat: ", <-ChanErrorHandling)
	}
	ctx.JSON(http.StatusOK, body)

}

// @Router       chats/ [put]
func (s Service) UpdateChat(ctx *gin.Context) {
	body := entity.Chat{}
	ctx.BindJSON(body)

	go s.Storage.UpdateChat(body, ChanErrorHandling)
	if <-ChanErrorHandling != nil {
		slog.Info("failed to update chat: ", <-ChanErrorHandling)
		ctx.JSON(http.StatusBadRequest, <-ChanErrorHandling)
	}
	ctx.JSON(http.StatusOK, body)
}

// @Router       chats/{chat_id} [delete]
func (s Service) DeleteChat(ctx *gin.Context) {
	chat_id := ctx.Param("chat_id")

	go s.Storage.DeleteChat(chat_id, ChanErrorHandling)
	if <-ChanErrorHandling != nil {
		ctx.JSON(http.StatusBadRequest, <-ChanErrorHandling)
	}
	ctx.JSON(http.StatusOK, chat_id)
}

// @Router       users/{chat_id} [post]
func (s Service) AddUserChat(ctx *gin.Context) {
	body := entity.Participant{}
	ctx.BindJSON(&body)

	go s.Storage.AddUser(body.Chat_ID, body.User_ID, ChanErrorHandling)
	if <-ChanErrorHandling != nil {
		ctx.JSON(http.StatusBadRequest, <-ChanErrorHandling)
	} else {
		ctx.JSON(http.StatusOK, body.User_ID)
	}
}

// @Router       users/{chat_id} [delete]
func (s Service) DropUserFromChat(ctx *gin.Context) {
	body := entity.Participant{}
	ctx.BindJSON(&body)

	s.Storage.DropUserFromChat(body.User_ID, body.Chat_ID, ChanErrorHandling)
	if <-ChanErrorHandling != nil {
		ctx.JSON(http.StatusBadRequest, <-ChanErrorHandling)
	} else {
		ctx.JSON(http.StatusOK, body.User_ID)
	}
}
