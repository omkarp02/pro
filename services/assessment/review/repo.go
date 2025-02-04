package review

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

	userObjectId, err := bson.ObjectIDFromHex(createModal.UserId)
	if err != nil {
		return "", err
	}

	dataToInsert := Review{
		UserId:           userObjectId,
		Timestamps:       &timestamp,
		ProductCode:      createModal.ProductCode,
		Rating:           createModal.Rating,
		ReviewText:       createModal.ReviewText,
		VerifiedPurchase: createModal.VerifiedPurchase,
		Status:           createModal.Status,
	}

	result, err := s.getColl().InsertOne(ctx, dataToInsert)

	if err != nil {
		return "", err
	}

	if id, ok := result.InsertedID.(bson.ObjectID); ok {
		return id.Hex(), nil
	}

	return "", errutil.ErrDatabase

}

func (s *Repo) FindByFilter(ctx context.Context, filterListModel FilterListModel, project []string, inclusive bool) ([]Review, error) {

	var list []Review

	query := bson.M{}

	page := filterListModel.Page
	limit := filterListModel.Limit

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

func (s *Repo) FindById(ctx context.Context, id string, project []string, inclusive bool) (Review, error) {

	var model Review

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

func (s *Repo) UpdateVotes(ctx context.Context, reviewId string, helpful bool) error {
	objID, err := bson.ObjectIDFromHex(reviewId)
	if err != nil {
		return errors.New("invalid reviewId format")
	}

	// Define filter to find the review
	filter := bson.M{"_id": objID}

	// Define update to increment HelpfulVotes
	update := bson.M{"$inc": bson.M{"helpfulVotes": 1}}
	if !helpful {
		update = bson.M{"$inc": bson.M{"notHelpfulVotes": 1}}
	}

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
