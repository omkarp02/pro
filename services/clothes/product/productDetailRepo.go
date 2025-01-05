package product

import (
	"context"
	"errors"
	"fmt"

	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/services/utils/store"
	"github.com/omkarp02/pro/utils/errutil"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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
		ImgLink:     createProductDetailModel.ImgLink,
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

func (s *ProductDetailRepo) FindById(ctx context.Context, id string, project []string, exclusive bool) (ProductDetail, error) {

	var productDetail ProductDetail

	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return productDetail, fmt.Errorf("invalid id format: %v", err)
	}

	filter := bson.M{"_id": objectId}

	projection := bson.M{}
	for _, field := range project {
		if exclusive {
			projection[field] = 0
		} else {
			projection[field] = 1
		}
	}

	findOneOptions := options.FindOne().SetProjection(projection)

	err = s.getColl().FindOne(ctx, filter, findOneOptions).Decode(&productDetail)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return productDetail, errutil.NotFound("Product")
		}
		return productDetail, errutil.NotFound("Product")
	}

	return productDetail, nil

}

func (s *ProductDetailRepo) GetProductPriceBySize(ctx context.Context, id string, size string) (float64, error) {

	var productDetail ProductDetail

	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return 0, fmt.Errorf("invalid id format: %v", err)
	}

	filter := bson.M{"_id": objectId}

	projection := bson.M{
		"variations": bson.M{
			"$elemMatch": bson.M{
				"size": size,
			},
		},
	}

	findOneOptions := options.FindOne().SetProjection(projection)

	err = s.getColl().FindOne(ctx, filter, findOneOptions).Decode(&productDetail)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 0, errutil.NotFound("Product")
		}
		return 0, errutil.NotFound("Product")
	}

	return productDetail.Variations[0].Price, nil

}
