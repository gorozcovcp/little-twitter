package mongo

import (
	"context"
	"fmt"
	"time"

	"github.com/gorozcovcp/little-twitter/internal/domain/model"
	"github.com/gorozcovcp/little-twitter/internal/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoTweetRepository struct {
	collection *mongo.Collection
}

func NewMongoTweetRepository(db *mongo.Database) ports.TweetRepository {
	return &MongoTweetRepository{
		collection: db.Collection("tweets"),
	}
}

func (r *MongoTweetRepository) Save(ctx context.Context, tweet *model.Tweet) error {
	_, err := r.collection.InsertOne(ctx, tweet)
	return err
}

func (r *MongoTweetRepository) GetByUsersBefore(ctx context.Context, userIDs []string, before time.Time, limit int) ([]model.Tweet, error) {
	fmt.Printf("GetByUsersBefore------userIDs:%v\n", userIDs)
	filter := bson.M{"user_id": bson.M{"$in": userIDs}}
	if !before.IsZero() {
		filter["created"] = bson.M{"$lt": before}
	}

	opts := options.Find().SetSort(bson.D{{Key: "created", Value: -1}}).SetLimit(int64(20))
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tweets []model.Tweet
	if err := cursor.All(ctx, &tweets); err != nil {
		return nil, err
	}
	return tweets, nil
}
