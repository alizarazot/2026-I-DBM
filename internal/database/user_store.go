package database

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"errors"
	"runtime"
	"time"

	modeldb "github.com/alizarazot/2026-i-dbm/internal/database/model"
	"github.com/alizarazot/2026-i-dbm/internal/model"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/argon2"
)

var (
	ErrUserAlreadyExists      = errors.New("user already exists")
	ErrUserNotExists          = errors.New("user doesn't exists")
	ErrUserInvalidCredentials = errors.New("invalid credentials")
)

// TODO: Should we store it in an environment variable?
const (
	PasswordTime   = 1
	PasswordMemory = 64 * 1024
	PasswordKeyLen = 128
)

var (
	PasswordSaltLen = min(PasswordKeyLen/2, 16)
	PasswordThreads = uint8(min(runtime.NumCPU(), 4))
)

type UserStore struct {
	c *mongo.Collection
}

func NewUserStore(client *mongo.Client, database string, collection string) *UserStore {
	return &UserStore{c: client.Database(database).Collection(collection)}
}

func (us *UserStore) AddUser(ctx context.Context, user *model.User, password string) error {
	// TODO: This should only check for email existence, currently it overwrites emails.
	if _, err := us.VerifyCredentials(ctx, user.Email, password); err == nil {
		return ErrUserAlreadyExists
	}

	var userdb modeldb.User

	userdb.Email = user.Email

	userdb.PasswordTime = PasswordTime
	userdb.PasswordMemory = PasswordMemory
	userdb.PasswordThreads = PasswordThreads
	userdb.PasswordKeyLen = PasswordKeyLen

	userdb.PasswordSalt = make([]byte, PasswordSaltLen)
	rand.Read(userdb.PasswordSalt)

	userdb.PasswordHash = argon2.IDKey([]byte(password), userdb.PasswordSalt, userdb.PasswordTime, userdb.PasswordMemory, userdb.PasswordThreads, userdb.PasswordKeyLen)

	userdb.Role = user.Role.CanonicalString()
	userdb.UserInfo = modeldb.UserInfo{
		FirstName:     user.Info.FirstName,
		MiddleName:    user.Info.MiddleName,
		FirstSurname:  user.Info.FirstSurname,
		SecondSurname: user.Info.SecondSurname,
		Birthdate:     user.Info.Birthdate,
		Genre:         user.Info.Genre.CanonicalString(),
	}

	userdb.CreatedAt = time.Now()
	userdb.UpdatedAt = userdb.CreatedAt

	_, err := us.c.InsertOne(ctx, userdb)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserStore) EditUser(ctx context.Context, user *model.User) error {
	var userdb modeldb.User
	if err := us.c.FindOne(ctx, bson.D{{Key: "email", Value: user.Email}}).Decode(&userdb); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrUserNotExists
		}

		return err
	}

	userdb.Email = user.Email
	userdb.Role = user.Role.CanonicalString()
	userdb.UserInfo = modeldb.UserInfo{
		FirstName:     user.Info.FirstName,
		MiddleName:    user.Info.MiddleName,
		FirstSurname:  user.Info.FirstSurname,
		SecondSurname: user.Info.SecondSurname,
		Birthdate:     user.Info.Birthdate,
		Genre:         user.Info.Genre.CanonicalString(),
	}
	userdb.UpdatedAt = userdb.CreatedAt

	_, err := us.c.ReplaceOne(ctx, bson.D{{Key: "email", Value: user.Email}}, userdb)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserStore) DeleteUser(ctx context.Context, email string) error {
	if _, err := us.c.DeleteOne(ctx, bson.D{{Key: "email", Value: email}}); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return ErrUserNotExists
		}
		return err
	}

	return nil
}

func (us *UserStore) VerifyCredentials(ctx context.Context, email string, password string) (*model.User, error) {
	var userdb modeldb.User
	if err := us.c.FindOne(ctx, bson.D{{Key: "email", Value: email}}).Decode(&userdb); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserInvalidCredentials
		}

		return nil, err
	}

	hashedPassword := argon2.IDKey([]byte(password), userdb.PasswordSalt, userdb.PasswordTime, userdb.PasswordMemory, userdb.PasswordThreads, userdb.PasswordKeyLen)

	if subtle.ConstantTimeCompare(hashedPassword, userdb.PasswordHash) == 0 {
		return nil, ErrUserInvalidCredentials
	}

	return dbUserToModelUser(userdb), nil
}

func (us *UserStore) GetUsers(ctx context.Context, page, limit uint) ([]*model.User, uint, error) {
	cursor, err := us.c.Aggregate(
		ctx,
		mongo.Pipeline{
			bson.D{{Key: "$skip", Value: page * limit}},
			bson.D{{Key: "$limit", Value: limit}},
		},
	)
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
		},
	)
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

func (us *UserStore) GetUserID(ctx context.Context, email string) (string, error) {
	var userdb modeldb.User
	if err := us.c.FindOne(ctx, bson.D{{Key: "email", Value: email}}).Decode(&userdb); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", ErrUserNotExists
		}

		return "", err
	}

	return userdb.ID.Hex(), nil
}

func (us *UserStore) GetUserEmail(ctx context.Context, id string) (string, error) {
	userid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return "", err
	}

	var userdb modeldb.User
	if err := us.c.FindOne(ctx, bson.D{{Key: "_id", Value: userid}}).Decode(&userdb); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return "", ErrUserNotExists
		}

		return "", err
	}

	return userdb.Email, nil
}

func dbUserToModelUser(userdb modeldb.User) *model.User {
	return &model.User{
		Email: userdb.Email,
		Role:  model.NewUserRole(userdb.Role),
		Info: model.UserInfo{
			FirstName:     userdb.UserInfo.FirstName,
			MiddleName:    userdb.UserInfo.MiddleName,
			FirstSurname:  userdb.UserInfo.FirstSurname,
			SecondSurname: userdb.UserInfo.SecondSurname,
			Birthdate:     userdb.UserInfo.Birthdate,
			Genre:         model.NewUserGenre(userdb.UserInfo.Genre),
		},
	}
}
