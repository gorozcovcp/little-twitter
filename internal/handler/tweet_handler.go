package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorozcovcp/little-twitter/internal/domain/service"
)

type TweetHandler struct {
	tweetService *service.TweetService
}

func NewTweetHandler(tweetService *service.TweetService) *TweetHandler {
	return &TweetHandler{tweetService: tweetService}
}

type postTweetRequest struct {
	UserID  string `json:"user_id" binding:"required"`
	Content string `json:"content" binding:"required,max=280"`
}

func (h *TweetHandler) PostTweet(c *gin.Context) {
	var req postTweetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if err := h.tweetService.PostTweet(c.Request.Context(), req.UserID, req.Content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to post tweet"})
		return
	}

	c.Status(http.StatusCreated)
}
