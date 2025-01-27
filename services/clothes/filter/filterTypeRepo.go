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

type FilterTypeRepo struct {
	*db.Database
	collName string
}

func NewRepoFilterType(curDb *db.Database, collName string) *FilterTypeRepo {
	store := &FilterTypeRepo{
		Database: curDb,
		collName: collName,
	}

	return store
}

func (s *FilterTypeRepo) getColl() *mongo.Collection {
	return s.DB.Database(s.DBName).Collection(s.collName)
}

func (s *FilterTypeRepo) Create(ctx context.Context, createFilterTypeModal CreateFilterTypeModal) (string, error) {

	filterType := FilterType{
		Name:       createFilterTypeModal.Name,
		Timestamps: store.GetCurrentTimestamps(),
	}

	result, err := s.getColl().InsertOne(ctx, filterType)
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

func (s *FilterTypeRepo) FindByFilter(ctx context.Context, filterListModel FilterListModel, project []string, inclusive bool) ([]FilterType, error) {

	var list []FilterType

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
