package filter

import (
	"context"

	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/services/utils/store"
	"github.com/omkarp02/pro/utils/errutil"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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
