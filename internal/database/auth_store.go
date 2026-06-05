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
	ErrAuthUserAlreadyExists  = errors.New("user already exists")
	ErrAuthInvalidCredentials = errors.New("invalid credentials")
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

type AuthStore struct {
	c *mongo.Collection
}

func NewAuthStore(client *mongo.Client, database string, collection string) *AuthStore {
	return &AuthStore{c: client.Database(database).Collection(collection)}
}

func (as *AuthStore) AddUser(ctx context.Context, user *model.User, password string) error {
	if _, err := as.VerifyCredentials(ctx, user.Email, password); err == nil {
		return ErrAuthUserAlreadyExists
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

	userdb.Role = modelRoleToDBRole(user.Role)
	userdb.UserInfo = modeldb.UserInfo{
		FirstName:     user.Info.FirstName,
		MiddleName:    user.Info.MiddleName,
		FirstSurname:  user.Info.FirstSurname,
		SecondSurname: user.Info.SecondSurname,
		Birthdate:     user.Info.Birthdate,
		Genre:         modelGenreToDBGenre(user.Info.Genre),
	}

	userdb.CreatedAt = time.Now()
	userdb.UpdatedAt = userdb.CreatedAt

	_, err := as.c.InsertOne(ctx, userdb)
	if err != nil {
		return err
	}

	return nil
}

func (as *AuthStore) VerifyCredentials(ctx context.Context, email string, password string) (*model.User, error) {
	var userdb modeldb.User
	if err := as.c.FindOne(ctx, bson.D{
		{
			Key: "email",
			Value: bson.D{
				{
					Key:   "$eq",
					Value: email,
				},
			},
		},
	}).Decode(&userdb); err != nil {
		return nil, ErrAuthInvalidCredentials
	}

	hashedPassword := argon2.IDKey([]byte(password), userdb.PasswordSalt, userdb.PasswordTime, userdb.PasswordMemory, userdb.PasswordThreads, userdb.PasswordKeyLen)

	if subtle.ConstantTimeCompare(hashedPassword, userdb.PasswordHash) == 0 {
		return nil, ErrAuthInvalidCredentials
	}

	user := model.User{
		Email: email,
		Role:  dbRoleToModelRole(userdb.Role),
		Info: model.UserInfo{
			FirstName:     userdb.UserInfo.FirstName,
			MiddleName:    userdb.UserInfo.MiddleName,
			FirstSurname:  userdb.UserInfo.FirstSurname,
			SecondSurname: userdb.UserInfo.SecondSurname,
			Birthdate:     userdb.UserInfo.Birthdate,
			Genre:         dbGenreToModelGenre(userdb.UserInfo.Genre),
		},
	}

	return &user, nil
}

func modelRoleToDBRole(role model.UserRole) string {
	switch role {
	case model.UserRoleManager:
		return "manager"
	case model.UserRoleTeacher:
		return "teacher"
	case model.UserRoleStudent:
		return "student"
	default:
		panic("invalid role")
	}
}

func dbRoleToModelRole(role string) model.UserRole {
	switch role {
	case "manager":
		return model.UserRoleManager
	case "teacher":
		return model.UserRoleTeacher
	case "student":
		return model.UserRoleStudent
	default:
		panic("invalid role")
	}
}

func modelGenreToDBGenre(genre model.UserGenre) string {
	switch genre {
	case model.UserGenreMale:
		return "male"
	case model.UserGenreFemale:
		return "female"
	case model.UserGenreOther:
		return "other"
	default:
		panic("invalid genre")
	}
}

func dbGenreToModelGenre(genre string) model.UserGenre {
	switch genre {
	case "male":
		return model.UserGenreMale
	case "female":
		return model.UserGenreFemale
	case "other":
		return model.UserGenreOther
	default:
		panic("invalid genre")
	}
}
