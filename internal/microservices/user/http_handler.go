package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpHandler struct {
	service *Service
}

func (h *HttpHandler) RegisterRoutes(router *gin.Engine) {
	v1 := router.Group("api/v1")
	userGroup := v1.Group("/users")
	{
		userGroup.GET("/:id", h.GetUser)
		userGroup.POST("/", h.CreateUser)
	}
}

func (h *HttpHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	user := h.service.GetUser(id)
	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *HttpHandler) CreateUser(c *gin.Context) {
	var req struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := h.service.CreateUser(req.Name)
	c.JSON(http.StatusCreated, gin.H{"user": user})
}
