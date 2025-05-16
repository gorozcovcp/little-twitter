package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

func init() {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}

	ensureIndexes()
}

func ensureIndexes() {
	tweetCollection := client.Database("uala").Collection("tweets")
	userCollection := client.Database("uala").Collection("users")

	_, err := tweetCollection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_id", Value: 1},
			{Key: "created", Value: -1},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create tweet index: %v", err)
	}

	_, err = userCollection.Indexes().CreateOne(context.TODO(), mongo.IndexModel{
		Keys: bson.D{{Key: "_id", Value: 1}},
	})
	if err != nil {
		log.Fatalf("Failed to create user index: %v", err)
	}
}

type Tweet struct {
	UserID  string    `json:"user_id" bson:"user_id"`
	Content string    `json:"content" bson:"content"`
	Created time.Time `json:"created" bson:"created"`
}

type User struct {
	ID      string   `json:"id" bson:"_id"`
	Follows []string `json:"follows" bson:"follows"`
}

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

	tweet.Created = time.Now()
	collection := client.Database("uala").Collection("tweets")
	_, err := collection.InsertOne(context.TODO(), tweet)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to post tweet"})
		return
	}

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

	collection := client.Database("uala").Collection("users")
	_, err := collection.UpdateOne(
		context.TODO(),
		bson.M{"_id": req.UserID},
		bson.M{"$addToSet": bson.M{"follows": req.FollowID}},
		options.Update().SetUpsert(true),
	)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to follow user"})
		return
	}

	c.JSON(200, gin.H{"message": "User followed"})
}

func getTimeline(c *gin.Context) {
	userID := c.Param("userID")
	collection := client.Database("uala").Collection("users")

	var user User
	err := collection.FindOne(context.TODO(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		c.JSON(404, gin.H{"error": "User not found"})
		return
	}

	limitParam := c.DefaultQuery("limit", "20")
	sinceParam := c.DefaultQuery("since", "")

	limit, err := strconv.Atoi(limitParam)
	if err != nil || limit <= 0 {
		limit = 20
	}

	filter := bson.M{"user_id": bson.M{"$in": user.Follows}}
	if sinceParam != "" {
		if sinceTime, err := time.Parse(time.RFC3339, sinceParam); err == nil {
			filter["created"] = bson.M{"$lt": sinceTime}
		}
	}

	tweetCollection := client.Database("uala").Collection("tweets")
	findOptions := options.Find().SetSort(bson.D{{Key: "created", Value: -1}}).SetLimit(int64(limit))
	cursor, err := tweetCollection.Find(context.TODO(), filter, findOptions)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get timeline"})
		return
	}
	defer cursor.Close(context.TODO())

	var timeline []Tweet
	if err := cursor.All(context.TODO(), &timeline); err != nil {
		c.JSON(500, gin.H{"error": "Failed to parse tweets"})
		return
	}

	c.JSON(200, timeline)
}
