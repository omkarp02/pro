package categories

import (
	"context"

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

	store.createIndexes()
	return store
}

func (s *Repo) createIndexes() error {
	collection := s.getColl()

	// Define the unique index for the "email" field
	catIdIndexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "catId", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	slugIndexModal := mongo.IndexModel{
		Keys:    bson.D{{Key: "slug", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{catIdIndexModel, slugIndexModal})
	return err
}

func (s *Repo) getColl() *mongo.Collection {
	return s.DB.Database(s.DBName).Collection(s.collName)
}

func (s *Repo) Create(ctx context.Context, createCategoryModal CreateCategoryModal) (string, error) {

	cat := Category{
		CatId:       createCategoryModal.CatId,
		Name:        createCategoryModal.Name,
		Description: createCategoryModal.Description,
		ImgLink:     createCategoryModal.ImgLink,
		Icon:        createCategoryModal.Icon,
		IsActive:    createCategoryModal.IsActive,
		Slug:        createCategoryModal.Slug,
		Timestamps:  store.GetCurrentTimestamps(),
	}

	result, err := s.getColl().InsertOne(ctx, cat)

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
