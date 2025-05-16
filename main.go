package main

import (
	"sync"

	"github.com/gin-gonic/gin"
)

// Tweet estructura básica de un tweet
type Tweet struct {
	UserID  string `json:"user_id"`
	Content string `json:"content"`
}

// User estructura de usuario
type User struct {
	ID      string   `json:"id"`
	Follows []string `json:"follows"`
}

// In-memory storage
var (
	users  = sync.Map{}
	tweets = sync.Map{}
)

func main() {
	r := gin.Default()

	r.POST("/tweet", postTweet)
	r.POST("/follow", followUser)
	r.GET("/timeline/:userID", getTimeline)

	r.Run(":8080")
}

func postTweet(c *gin.Context) {
	var tweet Tweet
	if err := c.ShouldBindJSON(&tweet); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	if len(tweet.Content) > 280 {
		c.JSON(400, gin.H{"error": "Tweet exceeds character limit"})
		return
	}

	tweets.Store(tweet.UserID, tweet)
	c.JSON(200, gin.H{"message": "Tweet posted"})
}

func followUser(c *gin.Context) {
	var req struct {
		UserID   string `json:"user_id"`
		FollowID string `json:"follow_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	value, exists := users.Load(req.UserID)
	if !exists {
		value = &User{ID: req.UserID, Follows: []string{}}
	}
	user := value.(*User)
	user.Follows = append(user.Follows, req.FollowID)
	users.Store(req.UserID, user)

	c.JSON(200, gin.H{"message": "User followed"})
}

func getTimeline(c *gin.Context) {
	userID := c.Param("userID")
	value, exists := users.Load(userID)
	if !exists {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	user := value.(*User)
	var timeline []Tweet

	for _, followID := range user.Follows {
		if tweet, ok := tweets.Load(followID); ok {
			timeline = append(timeline, tweet.(Tweet))
		}
	}

	c.JSON(200, timeline)
}
