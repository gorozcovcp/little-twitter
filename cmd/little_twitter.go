package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorozcovcp/little-twitter/internal/domain/service"
	"github.com/gorozcovcp/little-twitter/internal/handler"
	"github.com/gorozcovcp/little-twitter/internal/repository/mongo"
	"github.com/gorozcovcp/little-twitter/internal/repository/redis"
	rds "github.com/redis/go-redis/v9"
	mongocli "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()

	// MongoDB
	mongoClient, err := mongocli.Connect(ctx, options.Client().ApplyURI("mongodb://mongo:27017"))
	if err != nil {
		log.Fatal(err)
	}
	db := mongoClient.Database("uala")

	// Redis
	redisClient := rds.NewClient(&rds.Options{
		Addr: "redis:6379",
	})

	// Repositories
	tweetRepo := mongo.NewMongoTweetRepository(db)
	userRepo := mongo.NewMongoUserRepository(db)
	timelineCache := redis.NewRedisTimelineCache(redisClient, 30*time.Second)

	// Services
	tweetService := service.NewTweetService(tweetRepo, userRepo, timelineCache)
	userService := service.NewUserService(userRepo, timelineCache)

	// Handlers
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
