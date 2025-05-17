package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorozcovcp/little-twitter/internal/domain/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

type followRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	FollowID string `json:"follow_id" binding:"required"`
}

func (h *UserHandler) Follow(c *gin.Context) {
	var req followRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("Follow binding------error: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	err := h.userService.Follow(c.Request.Context(), req.UserID, req.FollowID)
	if err != nil {
		fmt.Printf("Follow------error: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to follow user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Following user"})
}
