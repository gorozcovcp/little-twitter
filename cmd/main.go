package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorozcovcp/little-twitter/config"
	"github.com/gorozcovcp/little-twitter/internal/domain/service"
	"github.com/gorozcovcp/little-twitter/internal/handler"
	mongorepo "github.com/gorozcovcp/little-twitter/internal/repository/mongo"
	redisrepo "github.com/gorozcovcp/little-twitter/internal/repository/redis"
)

func main() {
	config := config.LoadConfig()

	// MongoDB
	db := mongorepo.NewMongoDatabase(config.MongoURI, config.DBName)

	// Redis
	redisClient := redisrepo.NewRedisClient(config.RedisAddr)

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

	r := setupRouter(tweetHandler, userHandler, timelineHandler)

	srv := &http.Server{
		Addr:    config.ServerAddr,
		Handler: r,
	}

	// Graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Printf("Server running on %s", config.ServerAddr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %s", err)
	}

	log.Println("Server exiting")
}

func setupRouter(tweetHandler *handler.TweetHandler, userHandler *handler.UserHandler, timelineHandler *handler.TimelineHandler) *gin.Engine {
	r := gin.Default()
	r.POST("/tweet", tweetHandler.PostTweet)
	r.POST("/follow", userHandler.Follow)
	r.GET("/timeline/:userID", timelineHandler.GetTimeline)

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	return r
}
