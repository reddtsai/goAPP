package chatbot

import (
	"github.com/gin-gonic/gin"
)

type HttpHandler struct {
	// service *Service
}

func (h *HttpHandler) RegisterRoutes(router *gin.Engine) {
	v1 := router.Group("api/v1")

	chatbotGroup := v1.Group("/chatbot")
	userGroup := chatbotGroup.Group("/chats")
	{
		userGroup.GET("", h.ListChats)
		userGroup.POST("", h.CreateChat)
		userGroup.DELETE("/:chat_id", h.DeleteChat)
		userGroup.POST("/:chat_id/messages", h.SendMessage)
		userGroup.GET("/:chat_id/messages", h.ListMessages)
	}
}

func (h *HttpHandler) CreateChat(c *gin.Context) {

}

func (h *HttpHandler) ListChats(c *gin.Context) {

}

func (h *HttpHandler) DeleteChat(c *gin.Context) {

}

func (h *HttpHandler) SendMessage(c *gin.Context) {

}

func (h *HttpHandler) ListMessages(c *gin.Context) {

}
