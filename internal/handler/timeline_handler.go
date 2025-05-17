package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorozcovcp/little-twitter/internal/domain/service"
)

type TimelineHandler struct {
	tweetService *service.TweetService
}

func NewTimelineHandler(tweetService *service.TweetService) *TimelineHandler {
	return &TimelineHandler{tweetService: tweetService}
}

func (h *TimelineHandler) GetTimeline(c *gin.Context) {
	userID := c.Param("userID")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing user ID"})
		return
	}

	// Pagination params
	limit := 20
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	since := time.Now()
	if s := c.Query("since"); s != "" {
		if t, err := time.Parse(time.RFC3339, s); err == nil {
			since = t
		}
	}

	timeline, err := h.tweetService.GetTimeline(c.Request.Context(), userID, since, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get timeline"})
		return
	}

	c.JSON(http.StatusOK, timeline)
}
