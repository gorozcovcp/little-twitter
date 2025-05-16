package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const mongoURI = "mongodb://localhost:27017"

var client *mongo.Client

func init() {
	var err error
	client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
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

	tweetCollection := client.Database("uala").Collection("tweets")
	cursor, err := tweetCollection.Find(context.TODO(), bson.M{"user_id": bson.M{"$in": user.Follows}})
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
