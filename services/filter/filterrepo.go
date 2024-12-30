package filter

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

func NewRepoFilter(curDb *db.Database, collName string) *Repo {
	store := &Repo{
		Database: curDb,
		collName: collName,
	}

	return store
}

func (s *Repo) getColl() *mongo.Collection {
	return s.DB.Database(s.DBName).Collection(s.collName)
}

func (s *Repo) Create(ctx context.Context, createFilterModal CreateFilterModal) (string, error) {

	typeId, err := bson.ObjectIDFromHex(createFilterModal.Type)
	if err != nil {
		return "", err
	}

	filter := Filter{
		Name:       createFilterModal.Name,
		Type:       typeId,
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
