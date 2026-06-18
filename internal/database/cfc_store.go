package database

import (
	"context"
	"time"

	modeldb "github.com/alizarazot/2026-i-dbm/internal/database/model"
	"github.com/alizarazot/2026-i-dbm/internal/model"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type CFCStore struct {
	cCFC     *mongo.Collection
	cAnswers *mongo.Collection
}

func NewCFCStore(client *mongo.Client, database string, cfcCollection, cfcAnswersCollection string) *CFCStore {
	return &CFCStore{
		cCFC:     client.Database(database).Collection(cfcCollection),
		cAnswers: client.Database(database).Collection(cfcAnswersCollection),
	}
}

func (c *CFCStore) AddCFC(ctx context.Context, cfc model.CFC, userID string) error {
	dbcfc := modelCFCToDBCFC(&cfc, userID)

	dbcfc.CreatedAt = time.Now()
	dbcfc.UpdatedAt = dbcfc.CreatedAt

	if _, err := c.cCFC.InsertOne(ctx, dbcfc); err != nil {
		return err
	}

	return nil
}

type IncompleteCFC struct {
	CFC    *model.CFC
	UserID string
}

func (c *CFCStore) GetCFC(ctx context.Context, id string) (*IncompleteCFC, error) {
	var cfcdb modeldb.CFC

	bsonid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}

	if err := c.cCFC.FindOne(ctx, bson.D{{Key: "_id", Value: bsonid}}).Decode(&cfcdb); err != nil {
		return &IncompleteCFC{}, err
	}

	return &IncompleteCFC{dbCFCToModelCFC(&cfcdb), cfcdb.UserID.Hex()}, nil
}

func (c *CFCStore) ListCFCs(ctx context.Context, page, limit uint) ([]*IncompleteCFC, uint64, error) {
	cursor, err := c.cCFC.Aggregate(
		ctx,
		mongo.Pipeline{
			bson.D{{Key: "$skip", Value: page * limit}},
			bson.D{{Key: "$limit", Value: limit}},
		},
	)
	if err != nil {
		return nil, 0, err
	}

	var cfcsdb []modeldb.CFC
	if err := cursor.All(ctx, &cfcsdb); err != nil {
		return nil, 0, err
	}

	cfcs := make([]*IncompleteCFC, len(cfcsdb))
	for idx, cfc := range cfcsdb {
		cfcs[idx] = &IncompleteCFC{
			CFC:    dbCFCToModelCFC(&cfc),
			UserID: cfc.UserID.Hex(),
		}
	}

	count, err := c.cCFC.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, 0, err
	}

	return cfcs, uint64(count), nil
}

func (c *CFCStore) AddCFCAnswer(ctx context.Context, answer *model.CFCAnswer, userId string) error {
	// TODO: Check that there isn't a previous answer.

	id, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		panic(err)
	}

	answerdb := modelCFCAnswerToDBCFCAnswer(answer)
	answerdb.UserID = id
	answerdb.CreatedAt = time.Now()
	answerdb.UpdatedAt = answerdb.CreatedAt

	if _, err := c.cAnswers.InsertOne(ctx, answerdb); err != nil {
		return err
	}

	return nil
}

type IncompleteCFCAnswer struct {
	CFCAnswer *model.CFCAnswer
	UserID    string
}

func (c *CFCStore) GetCFCAnswer(ctx context.Context, cfcID string) (*IncompleteCFCAnswer, error) {
	var answerdb modeldb.CFCAnswer
	if err := c.cAnswers.FindOne(ctx, bson.D{{Key: "cfcId", Value: cfcID}}).Decode(&answerdb); err != nil {
		return nil, err
	}

	return &IncompleteCFCAnswer{
		CFCAnswer: dbCFCAnswerToModelCFCAnswer(&answerdb),
		UserID:    answerdb.UserID.Hex(),
	}, nil
}

func modelCFCAnswerToDBCFCAnswer(answer *model.CFCAnswer) *modeldb.CFCAnswer {
	id, err := bson.ObjectIDFromHex(answer.ID)
	if err != nil {
		panic(err)
	}
	cfcid, err := bson.ObjectIDFromHex(answer.CFCID)
	if err != nil {
		panic(err)
	}

	return &modeldb.CFCAnswer{
		ID:        id,
		CFCID:     cfcid,
		Answer:    answer.Answer,
		CreatedAt: time.Time{},
		UpdatedAt: answer.UpdatedAt,
	}
}

func dbCFCAnswerToModelCFCAnswer(answerdb *modeldb.CFCAnswer) *model.CFCAnswer {
	return &model.CFCAnswer{
		ID:        answerdb.ID.Hex(),
		CFCID:     answerdb.CFCID.Hex(),
		Answer:    answerdb.Answer,
		UpdatedAt: answerdb.UpdatedAt,
	}
}

func dbCFCToModelCFC(cfcdb *modeldb.CFC) *model.CFC {
	return &model.CFC{
		ID:        cfcdb.ID.Hex(),
		Subject:   cfcdb.Subject,
		Category:  model.NewCFCCategory(cfcdb.Category),
		Details:   cfcdb.Details,
		UpdatedAt: cfcdb.UpdatedAt,
	}
}

func modelCFCToDBCFC(cfc *model.CFC, userID string) *modeldb.CFC {
	id, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		panic(err)
	}

	return &modeldb.CFC{
		Subject:   cfc.Subject,
		UserID:    id,
		Category:  cfc.Category.CanonicalString(),
		Details:   cfc.Details,
		UpdatedAt: cfc.UpdatedAt,
	}
}
