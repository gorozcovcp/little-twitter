package mongo

import (
	"context"

	"github.com/gorozcovcp/little-twitter/internal/domain/model"
	"github.com/gorozcovcp/little-twitter/internal/ports"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoUserRepository struct {
	collection *mongo.Collection
}

func NewMongoUserRepository(db *mongo.Database) ports.UserRepository {
	return &MongoUserRepository{
		collection: db.Collection("users"),
	}
}

func (r *MongoUserRepository) Follow(ctx context.Context, userID, followID string) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": userID},
		bson.M{"$addToSet": bson.M{"follows": followID}},
		options.Update().SetUpsert(true),
	)
	return err
}

func (r *MongoUserRepository) GetByID(ctx context.Context, userID string) (*model.User, error) {
	var user model.User
	err := r.collection.FindOne(ctx, bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *MongoUserRepository) GetFollowersOf(ctx context.Context, followeeID string) ([]model.User, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"follows": followeeID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []model.User
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}
