package filter

import (
	"context"

	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/services/utils/store"
	"github.com/omkarp02/pro/utils/errutil"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type FilterRepo struct {
	*db.Database
	collName string
}

func NewRepoFilter(curDb *db.Database, collName string) *FilterRepo {
	store := &FilterRepo{
		Database: curDb,
		collName: collName,
	}

	return store
}

func (s *FilterRepo) getColl() *mongo.Collection {
	return s.DB.Database(s.DBName).Collection(s.collName)
}

func (s *FilterRepo) Create(ctx context.Context, createFilterModal CreateFilterModal) (string, error) {

	typeId, err := bson.ObjectIDFromHex(createFilterModal.Type)
	if err != nil {
		return "", err
	}
	categoryId, err := bson.ObjectIDFromHex(createFilterModal.Category)
	if err != nil {
		return "", err
	}

	filter := Filter{
		Name:       createFilterModal.Name,
		Type:       typeId,
		Status:     createFilterModal.Status,
		Category:   categoryId,
		Timestamps: store.GetCurrentTimestamps(),
	}

	result, err := s.getColl().InsertOne(ctx, filter)

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

func (s *FilterRepo) FindByFilter(ctx context.Context, filterListModel FilterListModel, project []string, inclusive bool) ([]Filter, error) {

	var list []Filter

	query := bson.M{}

	page := filterListModel.Page
	limit := filterListModel.Limit
	category := filterListModel.Category
	filterType := filterListModel.Type

	if len(category) != 0 {
		categoryObjectId, err := bson.ObjectIDFromHex(category)
		if err != nil {
			return list, err
		}
		query["category"] = categoryObjectId
	}
	if len(filterType) != 0 {
		typeObjectId, err := bson.ObjectIDFromHex(filterType)
		if err != nil {
			return list, err
		}
		query["type"] = typeObjectId
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
