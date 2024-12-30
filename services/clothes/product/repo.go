package product

import (
	"context"

	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/services/utils/store"
	"github.com/omkarp02/pro/utils/errutil"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ProductListRepo struct {
	*db.Database
	collName string
}

func NewProductListRepo(curDb *db.Database, collName string) *ProductListRepo {
	store := &ProductListRepo{
		Database: curDb,
		collName: collName,
	}

	return store
}

func (s *ProductListRepo) getColl() *mongo.Collection {
	return s.DB.Database(s.DBName).Collection(s.collName)
}

func (s *ProductListRepo) Create(ctx context.Context, createProductListModel CreateProductListModel) (string, error) {

	detailId, err := bson.ObjectIDFromHex(createProductListModel.Detail)
	if err != nil {
		return "", err
	}

	owner := ProductList{
		Detail:     detailId,
		Name:       createProductListModel.Name,
		Sizes:      createProductListModel.Sizes,
		Color:      createProductListModel.Color,
		ImgLink:    createProductListModel.ImgLink,
		Price:      createProductListModel.Price,
		Stock:      createProductListModel.Stock,
		Discount:   createProductListModel.Discount,
		Timestamps: store.GetCurrentTimestamps(),
	}

	result, err := s.getColl().InsertOne(ctx, owner)

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
