package city

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
	return nil
}

func (s *Repo) Create(ctx context.Context, createModal CreateModal) (string, error) {
	timestamp := store.GetCurrentTimestamps()

	objectId, err := bson.ObjectIDFromHex(createModal.StateId)
	if err != nil {
		return "", err
	}

	dataToInsert := City{
		Name:       createModal.Name,
		StateId:    objectId,
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

func (s *Repo) FindByFilter(ctx context.Context, filterListModel FilterListModel, project []string, inclusive bool) ([]TCityRes, error) {

	var list []TCityRes

	name := filterListModel.Name
	stateId := filterListModel.StateId
	page := filterListModel.Page
	limit := filterListModel.Limit

	matchQuery := bson.D{}

	if len(name) != 0 {
		matchQuery = append(matchQuery, bson.E{Key: "name", Value: bson.M{"$regex": name, "$options": "i"}})
	}

	if len(stateId) != 0 {

		stateObjectId, err := bson.ObjectIDFromHex(stateId)
		if err != nil {
			return list, err
		}

		matchQuery = append(matchQuery, bson.E{Key: "stateId", Value: stateObjectId})
	}

	pipeline := mongo.Pipeline{
		{
			{Key: "$match", Value: matchQuery},
		},
		{{
			Key: "$lookup", Value: bson.D{
				{Key: "from", Value: "state"},         // The collection to join with
				{Key: "localField", Value: "stateId"}, // Field in the current collection
				{Key: "foreignField", Value: "_id"},   // Field in the "state" collection
				{Key: "as", Value: "stateDetails"},    // Output array field
				{Key: "pipeline", Value: bson.A{
					bson.D{
						{Key: "$project", Value: bson.D{
							{Key: "name", Value: 1},
						}},
					},
				}},
			},
		}},
		{{
			Key: "$unwind", Value: bson.D{
				{Key: "path", Value: "$stateDetails"},
				{Key: "preserveNullAndEmptyArrays", Value: true},
			},
		}},
		{
			{Key: "$project", Value: store.GenerateProjection(project, inclusive)},
		},
		{
			{Key: "$skip", Value: limit * (page - 1)},
		},
		{
			{Key: "$limit", Value: limit},
		},
	}

	cursor, err := s.getColl().Aggregate(ctx, pipeline)

	if err != nil {
		return nil, err
	}

	if err := cursor.All(context.TODO(), &list); err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	return list, nil
}

func (s *Repo) FindById(ctx context.Context, id string, project []string, inclusive bool) (City, error) {

	var model City

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
