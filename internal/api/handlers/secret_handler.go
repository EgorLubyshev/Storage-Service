package handlers

import (
	"github.com/gin-gonic/gin"
)

type SecretHandler struct {
	service SecretService
}

func (h *SecretHandler) Create(c *gin.Context) {
	userID := c.GetString("user_id")

	var req struct {
		Data string `json:"data" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	id, err := h.service.Create(c.Request.Context(), userID, req.Data)
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to store secret"})
		return
	}

	c.JSON(201, gin.H{"id": id})
}
