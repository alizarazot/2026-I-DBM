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

func (us *UserStore) GetUsers(ctx context.Context, page, limit uint) ([]*model.User, uint, error) {
	cursor, err := us.c.Aggregate(
		ctx,
		mongo.Pipeline{
			bson.D{{Key: "$skip", Value: page * limit}},
			bson.D{{Key: "$limit", Value: limit}},
		})
	if err != nil {
		return nil, 0, err
	}

	var usersdb []modeldb.User
	if err := cursor.All(ctx, &usersdb); err != nil {
		return nil, 0, err
	}

	users := make([]*model.User, len(usersdb))
	for i, userdb := range usersdb {
		users[i] = dbUserToModelUser(userdb)
	}

	count, err := us.c.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, 0, err
	}

	return users, uint(count), nil
}

func (us *UserStore) GetUsersByRole(ctx context.Context, role model.UserRole, page, limit uint) ([]*model.User, uint, error) {
	cursor, err := us.c.Aggregate(
		ctx,
		mongo.Pipeline{
			bson.D{{Key: "$match", Value: bson.D{{Key: "role", Value: role.CanonicalString()}}}},
			bson.D{{Key: "$skip", Value: page * limit}},
			bson.D{{Key: "$limit", Value: limit}},
		})
	if err != nil {
		return nil, 0, err
	}

	var usersdb []modeldb.User
	if err := cursor.All(ctx, &usersdb); err != nil {
		return nil, 0, err
	}

	users := make([]*model.User, len(usersdb))
	for i, userdb := range usersdb {
		users[i] = dbUserToModelUser(userdb)
	}

	count, err := us.c.CountDocuments(ctx, bson.D{{Key: "role", Value: role.CanonicalString()}})
	if err != nil {
		return nil, 0, err
	}

	return users, uint(count), nil
}
