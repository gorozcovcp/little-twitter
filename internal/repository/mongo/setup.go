package mongo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDatabase(uri, dbName string) *mongo.Database {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatalf("Failed to create Mongo client: %v", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("Failed to connect to Mongo: %v", err)
	}

	db := client.Database(dbName)
	ensureIndexes(ctx, db)
	return db
}

func ensureIndexes(ctx context.Context, db *mongo.Database) {
	tweetCollection := db.Collection("tweets")
	userCollection := db.Collection("users")

	_, err := tweetCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{
			{Key: "user_id", Value: 1},
			{Key: "created", Value: -1},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create tweet index: %v", err)
	}

	_, err = userCollection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "_id", Value: 1}},
	})
	if err != nil {
		log.Fatalf("Failed to create user index: %v", err)
	}
}
