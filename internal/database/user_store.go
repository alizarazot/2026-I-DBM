package database

import (
	"context"

	modeldb "github.com/alizarazot/2026-i-dbm/internal/database/model"
	"github.com/alizarazot/2026-i-dbm/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserStore struct {
	c *mongo.Collection
}

func NewUserStore(client *mongo.Client, database string, collection string) *UserStore {
	return &UserStore{c: client.Database(database).Collection(collection)}
}

func (us *UserStore) GetUsers(ctx context.Context, page, limit uint) ([]*model.User, error) {
	cursor, err := us.c.Aggregate(
		ctx,
		mongo.Pipeline{
			bson.D{{Key: "$skip", Value: page * limit}},
			bson.D{{Key: "$limit", Value: limit}},
		})
	if err != nil {
		return nil, err
	}

	var usersdb []modeldb.User
	if err := cursor.All(ctx, &usersdb); err != nil {
		return nil, err
	}

	users := make([]*model.User, len(usersdb))
	for i, userdb := range usersdb {
		users[i] = dbUserToModelUser(userdb)
	}

	return users, nil
}

func (us *UserStore) GetTotalUsers(ctx context.Context) (uint, error) {
	count, err := us.c.CountDocuments(ctx, bson.D{})
	return uint(count), err
}
