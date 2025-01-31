package categories

import (
	"context"
	"fmt"

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
		Status:      createCategoryModal.Status,
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

func (s *Repo) Find(ctx context.Context, filterCategoryModal FilterCategoryModal, project []string, inclusive bool) ([]Category, error) {

	fmt.Println(filterCategoryModal)

	var categoryList []Category

	query := bson.M{"status": filterCategoryModal.Status}

	projection := bson.M{}
	for _, field := range project {
		if inclusive {
			projection[field] = 1
		} else {
			projection[field] = 0
		}
	}

	findOptions := options.Find().SetProjection(projection)

	cursor, err := s.getColl().Find(ctx, query, findOptions)
	if err != nil {
		return nil, err
	}
	if err := cursor.All(context.TODO(), &categoryList); err != nil {
		return nil, err
	}

	return categoryList, nil

}

func (s *Repo) Count(ctx context.Context, filterCategoryModal FilterCategoryModal) (int, error) {

	query := bson.M{"status": filterCategoryModal.Status}

	count, err := s.getColl().CountDocuments(ctx, query)
	if err != nil {
		return 0, err
	}

	return int(count), nil

}
