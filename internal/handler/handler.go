package handler

import (
	services "github.com/Flikest/PingviMessenger/internal/controller"
	"github.com/Flikest/PingviMessenger/pkg/middleware"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	Service *services.Service
}

func NewHandler(s *services.Service) *Handler {
	return &Handler{Service: s}
}

func (h Handler) InitRouter() *gin.Engine {
	router := gin.Default()

	v1 := router.Group("/chats")

	v1.Use(middleware.IsAuthorized)
	{
		v1.GET("/messenger/", h.Service.Сorrespondence)
		v1.GET("/:chat_name", h.Service.GetChat)
		v1.POST("/", h.Service.CreateChat)
		v1.PUT("/", h.Service.UpdateChat)
		v1.DELETE("/:chat_id", h.Service.DeleteChat)
	}

	v2 := router.Group("/message")

	v2.Use(middleware.IsAuthorized)
	{
		v2.GET("/:message_id", h.Service.GetMessage)
		v2.PUT("/", h.Service.UpdateMessage)
		v2.DELETE("/:message_id", h.Service.DelelteMessage)
	}

	v3 := router.Group("/users")

	v3.Use(middleware.IsAuthorized)
	{
		v3.POST("/:chat_id", h.Service.AddUserChat)
		v3.DELETE("/chat_id", h.Service.DropUserFromChat)
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
