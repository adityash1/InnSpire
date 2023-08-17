package db

import (
	"context"

	"github.com/adityash1/go-reservation-api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserStore interface {
	GetUserByID(context.Context, string) (*types.User, error)
}

type MongoUserStore struct {
	client *mongo.Client
	col    *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	return &MongoUserStore{
		client: client,
		col:    client.Database(DBNAME).Collection(userCol),
	}
}

func (s *MongoUserStore) GetUserByID(ctx context.Context, id string) (*types.User, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var user types.User
	if err := s.col.FindOne(ctx, bson.M{"_id": objId}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
