package product

import (
	"context"

	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/services/utils/store"
	"github.com/omkarp02/pro/utils/errutil"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ProductDetailRepo struct {
	*db.Database
	collName string
}

func NewProductDetailRepo(curDb *db.Database, collName string) *ProductDetailRepo {
	store := &ProductDetailRepo{
		Database: curDb,
		collName: collName,
	}

	return store
}

func (s *ProductDetailRepo) getColl() *mongo.Collection {
	return s.DB.Database(s.DBName).Collection(s.collName)
}

func (s *ProductDetailRepo) Create(ctx context.Context, createProductDetailModel CreateProductDetailModel) (string, error) {

	productDetail := ProductDetail{
		Description: createProductDetailModel.Description,
		Variations:  createProductDetailModel.Variations,
		Timestamps:  store.GetCurrentTimestamps(),
	}

	result, err := s.getColl().InsertOne(ctx, productDetail)

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
