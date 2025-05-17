package main

import (
	"fmt"
	"log"
	"net/http"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/gorozcovcp/little-twitter/internal/domain/service"
	"github.com/gorozcovcp/little-twitter/internal/handler"
	mongorepo "github.com/gorozcovcp/little-twitter/internal/repository/mongo"
	redisrepo "github.com/gorozcovcp/little-twitter/internal/repository/redis"
)

func main() {
	// MongoDB
	mongoURI := getEnv("MONGO_URI", "mongodb://mongo:27017/?directConnection=true")
	fmt.Printf("mongoURI: %s", mongoURI)
	db := mongorepo.NewMongoDatabase(mongoURI, "uala")
	fmt.Printf("db: %s", db.Name())
	// Redis
	redisAddr := getEnv("REDIS_ADDR", "redis:6379")
	redisClient := redisrepo.NewRedisClient(redisAddr)

	// Repositories
	tweetRepo := mongorepo.NewMongoTweetRepository(db)
	userRepo := mongorepo.NewMongoUserRepository(db)
	timelineCache := redisrepo.NewRedisTimelineCache(redisClient, redisrepo.DefaultTTL())

	// Services
	tweetService := service.NewTweetService(tweetRepo, userRepo, timelineCache)
	userService := service.NewUserService(userRepo, timelineCache)

	// HTTP Handlers
	tweetHandler := handler.NewTweetHandler(tweetService)
	userHandler := handler.NewUserHandler(userService)
	timelineHandler := handler.NewTimelineHandler(tweetService)

	r := gin.Default()
	r.POST("/tweet", tweetHandler.PostTweet)
	r.POST("/follow", userHandler.Follow)
	r.GET("/timeline/:userID", timelineHandler.GetTimeline)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	log.Println("Server running on :8080")
	srv.ListenAndServe()
}

func getEnv(key, fallback string) string {
	if value, ok := syscall.Getenv(key); ok {
		return value
	}
	return fallback
}
