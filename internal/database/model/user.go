package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID bson.ObjectID `bson:"_id,omitempty"`

	Email string

	PasswordSalt, PasswordHash                   []byte
	PasswordTime, PasswordMemory, PasswordKeyLen uint32
	PasswordThreads                              uint8

	UserInfo
	Role string

	CreatedAt, UpdatedAt time.Time
}

type UserInfo struct {
	PictureFileID bson.ObjectID `bson:"pictureFileId,omitempty"`

	FirstName, MiddleName, FirstSurname, SecondSurname string

	Birthdate time.Time
	Genre     string
}
