package review

import (
	"context"
	"log"

	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/services/utils/store"
	"github.com/omkarp02/pro/utils/errutil"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ReviewVoteRepo struct {
	*db.Database
	collName string
}

func NewReviewVoteRepo(curDb *db.Database, collName string) *ReviewVoteRepo {
	store := &ReviewVoteRepo{
		Database: curDb,
		collName: collName,
	}

	if err := store.createIndexes(); err != nil {
		log.Fatalf("Error while creating index of %s collection", collName)
	}

	return store
}

func (s *ReviewVoteRepo) getColl() *mongo.Collection {
	return s.DB.Database(s.DBName).Collection(s.collName)
}

func (s *ReviewVoteRepo) createIndexes() error {

	return nil
}

func (s *ReviewVoteRepo) Create(ctx context.Context, createModal CreateReviewVoteModal) (string, error) {

	timestamp := store.GetCurrentTimestamps()

	objectIds, err := store.SliceOfHexToObjectID(createModal.UserId, createModal.ReviewId)
	if err != nil {
		return "", err
	}

	dataToInsert := ReviewVote{
		UserId:    objectIds[0],
		ReviewId:  objectIds[1],
		IsHelpful: createModal.IsHelpful,
		CreatedAt: timestamp.CreatedAt,
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

func (s *ReviewVoteRepo) FindByFilter(ctx context.Context, filterListModel FilterListReviewVoteModel, project []string, inclusive bool) ([]ReviewVote, error) {

	var list []ReviewVote

	query := bson.M{}

	reviewId := filterListModel.ReviewId

	if len(reviewId) != 0 {
		query["reviewId"] = reviewId
	}

	// page := filterListModel.Page
	// limit := filterListModel.Limit

	// findOptions := options.Find().SetSkip(int64(limit * (page - 1))).SetLimit(int64(limit))

	// if len(project) > 0 {
	// 	projection := store.GenerateProjection(project, inclusive)
	// 	findOptions.SetProjection(projection)
	// }

	cursor, err := s.getColl().Find(ctx, query)
	if err != nil {
		return nil, err
	}
	if err := cursor.All(context.TODO(), &list); err != nil {
		return nil, err
	}

	return list, nil
}
