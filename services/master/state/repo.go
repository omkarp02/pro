package state

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/services/utils/store"
	"github.com/omkarp02/pro/utils/errutil"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

	if err := store.createIndexes(); err != nil {
		log.Fatalf("Error while creating index of %s collection", collName)
	}

	return store
}

func (s *Repo) getColl() *mongo.Collection {
	return s.DB.Database(s.DBName).Collection(s.collName)
}

func (s *Repo) createIndexes() error {
	collection := s.getColl()

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)

	return err
}

func (s *Repo) Create(ctx context.Context, createModal CreateModal) (string, error) {

	timestamp := store.GetCurrentTimestamps()

	dataToInsert := State{
		Name:       createModal.Name,
		Timestamps: &timestamp,
	}

	result, err := s.getColl().InsertOne(ctx, dataToInsert)

	if mongo.IsDuplicateKeyError(err) {
		return "", errutil.ErrDocumentAlreadyExist
	} else if err != nil {
		return "", err
	}

	if id, ok := result.InsertedID.(bson.ObjectID); ok {
		return id.Hex(), nil
	}

	return "", errutil.ErrDatabase

}

func (s *Repo) FindByFilter(ctx context.Context, filterListModel FilterListModel, project []string, inclusive bool) ([]State, error) {

	var list []State

	query := bson.M{}

	name := filterListModel.Name
	page := filterListModel.Page
	limit := filterListModel.Limit

	if len(name) != 0 {
		query["name"] = bson.M{"$regex": name, "$options": "i"}
	}

	findOptions := options.Find().SetSkip(int64(limit * (page - 1))).SetLimit(int64(limit))

	if len(project) > 0 {
		projection := store.GenerateProjection(project, inclusive)
		findOptions.SetProjection(projection)
	}

	cursor, err := s.getColl().Find(ctx, query, findOptions)
	if err != nil {
		return nil, err
	}
	if err := cursor.All(context.TODO(), &list); err != nil {
		return nil, err
	}

	return list, nil
}

func (s *Repo) FindById(ctx context.Context, id string, project []string, inclusive bool) (State, error) {

	var model State

	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return model, fmt.Errorf("invalid id format: %v", err)
	}

	filter := bson.M{"_id": objectId}
	findOneOptions := options.FindOne()

	if len(project) != 0 {
		projection := store.GenerateProjection(project, inclusive)
		findOneOptions.SetProjection(projection)
	}

	err = s.getColl().FindOne(ctx, filter, findOneOptions).Decode(&model)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model, errutil.NotFound("Product")
		}
		return model, err
	}

	return model, nil

}
