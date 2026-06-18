package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type CFC struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	Subject   string        `bson:"subject"`
	Category  string        `bson:"category"`
	UserID    bson.ObjectID `bson:"userId"`
	Details   string        `bson:"details"`
	CreatedAt time.Time     `bson:"createdAt"`
	UpdatedAt time.Time     `bson:"updatedAt"`
}

type CFCAnswer struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	CFCID     bson.ObjectID `bson:"cfcId"`
	UserID    bson.ObjectID `bson:"userId"`
	Answer    string        `bson:"answer"`
	CreatedAt time.Time     `bson:"createdAt"`
	UpdatedAt time.Time     `bson:"updatedAt"`
}
