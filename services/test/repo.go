package test

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/services/utils"
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

	//This is code for create single index

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)

	//here is dome index

	// Define the unique index for the "email" field
	// indexModel := mongo.IndexModel{
	// 	Keys:    bson.D{{Key: "email", Value: 1}},
	// 	Options: options.Index().SetUnique(true).SetPartialFilterExpression(bson.D{{Key: "email", Value: bson.D{{Key: "$exists", Value: true}, {Key: "$ne", Value: nil}}}}),
	// }

	//This is code for create many indexes

	catIdIndexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "catId", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	slugIndexModal := mongo.IndexModel{
		Keys:    bson.D{{Key: "slug", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err = collection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{catIdIndexModel, slugIndexModal})

	return err
}

func (s *Repo) Create(ctx context.Context, createModal CreateModal) (string, error) {

	// auditFields, err := store.GenerateCreateAuditFields(createModal.CreatorId)
	// if err != nil {
	// 	return "", err
	// }
	// timestamp := store.GetCurrentTimestamps()

	dataToInsert := Model{
		// AuditFields: auditFields,
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

func (s *Repo) FindByFilter(ctx context.Context, filterListModel FilterListModel, project []string, inclusive bool) ([]Model, error) {

	var list []Model

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

func (s *Repo) FindById(ctx context.Context, id string, project []string, inclusive bool) (Model, error) {

	var model Model

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

func (s *Repo) FindByFilterAggrgate(ctx context.Context, filterListModel FilterListModel, project []string, inclusive bool) ([]Model, error) {

	var list []Model

	name := filterListModel.Name
	page := filterListModel.Page
	limit := filterListModel.Limit

	matchQuery := bson.D{}

	if len(name) != 0 {
		matchQuery = append(matchQuery, bson.E{Key: "name", Value: bson.M{"$regex": name, "$options": "i"}})
	}

	// if len(stateId) != 0 {
	// 	matchQuery = append(matchQuery, bson.E{Key: "stateId", Value: stateId})
	// }

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

	fmt.Printf("%+v\n", list)

	return list, nil
}

func (s *Repo) UpdateById(ctx context.Context, field string, value string, id string) (utils.UPDATE_RESULT, error) {

	var updateResult utils.UPDATE_RESULT

	query := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{field: value}}

	result, err := s.getColl().UpdateOne(ctx, query, update)
	updateResult.MatchedCount = int(result.MatchedCount)
	updateResult.ModifiedCount = int(result.ModifiedCount)

	return updateResult, err
}

func (s *Repo) IncReviewCount(ctx context.Context, reviewId string) error {
	objID, err := bson.ObjectIDFromHex(reviewId)
	if err != nil {
		return errors.New("invalid review ID format")
	}

	// Define filter to find the review
	filter := bson.M{"_id": objID}

	// Define update to increment HelpfulVotes
	update := bson.M{"$inc": bson.M{"helpfulVotes": 1}}

	// Perform the update operation
	result, err := s.getColl().UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errutil.ErrDocumentNotFound
	}

	return nil
}
