package chatbot

import (
	"github.com/gin-gonic/gin"
)

type HttpHandler struct {
	service *Service
}

func (h *HttpHandler) RegisterRoutes(router *gin.Engine) {
	v1 := router.Group("api/v1")

	chatbotGroup := v1.Group("/chatbot")
	userGroup := chatbotGroup.Group("/chat")
	{
		userGroup.GET("", h.ListChats)
		userGroup.POST("", h.CreateChat)
		userGroup.DELETE("/:chat_id", h.DeleteChat)
		userGroup.POST("/:chat_id/message", h.SendMessage)
		userGroup.GET("/:chat_id/message", h.ListMessages)
	}
}

func (h *HttpHandler) CreateChat(c *gin.Context) {

}

func (h *HttpHandler) ListChats(c *gin.Context) {

}

func (h *HttpHandler) DeleteChat(c *gin.Context) {

}

func (h *HttpHandler) SendMessage(c *gin.Context) {
	params := SendMessageParams{
		Message: "你好",
	}
	result, err := h.service.SendMessage(params)
	if err != nil {
		c.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": result.Answer,
	})
}

func (h *HttpHandler) ListMessages(c *gin.Context) {

}
