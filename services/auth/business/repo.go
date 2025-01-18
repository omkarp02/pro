package bussiness

import (
	"context"

	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/services/utils/store"
	"github.com/omkarp02/pro/utils/errutil"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Repo struct {
	*db.Database
	collName string
}

func NewRepo(curDb *db.Database, collName string) *Repo {
	store := &Repo{
		Database: curDb,
		collName: collName,
	}

	store.createIndexes()

	return store
}

func (s *Repo) createIndexes() error {
	collection := s.getColl()

	// Define the unique index for the "email" field
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "user_id", Value: 1}},
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	return err
}

func (s *Repo) getColl() *mongo.Collection {
	return s.DB.Database(s.DBName).Collection(s.collName)
}

func (s *Repo) Create(ctx context.Context, createPayload CreateBusinessModel) (string, error) {

	ownerObjectId, err := bson.ObjectIDFromHex(createPayload.OwnerID)
	if err != nil {
		return "", err
	}

	newAddress := Business{
		Name:        createPayload.Name,
		OwnerID:     ownerObjectId,
		Category:    createPayload.Category,
		Description: createPayload.Description,
		Address:     store.Address(createPayload.Address),
		Contacts:    createPayload.Contacts,
		Website:     createPayload.Website,
		LogoUrl:     createPayload.LogoUrl,
		Active:      createPayload.Active,
		Timestamps:  store.GetCurrentTimestamps(),
	}

	result, err := s.getColl().InsertOne(ctx, newAddress)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return "", errutil.ErrDocumentAlreadyExist
		}
		return "", err
	}

	if id, ok := result.InsertedID.(bson.ObjectID); ok {
		return id.Hex(), nil
	}

	return "", errutil.ErrDatabase
}
