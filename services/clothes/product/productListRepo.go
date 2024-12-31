package product

import (
	"context"

	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/services/utils/store"
	"github.com/omkarp02/pro/utils/errutil"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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

func (s *ProductListRepo) FindByFilter(ctx context.Context, filterProductListModel FilterProductListModel) ([]ProductList, error) {

	var productList []ProductList

	query := bson.M{}

	sizes := filterProductListModel.Sizes
	name := filterProductListModel.Name
	color := filterProductListModel.Color
	maxPrice := filterProductListModel.MaxPrice
	minPrice := filterProductListModel.MinPrice
	page := filterProductListModel.Page
	limit := filterProductListModel.Limit

	if len(sizes) != 0 {
		query["sizes"] = bson.M{"$in": sizes}
	}
	if len(color) != 0 {
		query["color"] = color
	}
	if len(name) != 0 {
		query["name"] = bson.M{"$regex": name, "$options": "i"}
	}
	if maxPrice != 0 && minPrice != 0 {
		query["price"] = bson.M{"$gte": minPrice, "$lte": maxPrice}
	} else if maxPrice != 0 {
		query["price"] = bson.M{"$lte": maxPrice}
	} else if minPrice != 0 {
		query["price"] = bson.M{"$gte": minPrice}
	}

	findOptions := options.Find().SetSkip(int64(limit * (page - 1))).SetLimit(int64(limit))

	cursor, err := s.getColl().Find(ctx, query, findOptions)
	if err != nil {
		return nil, err
	}
	if err := cursor.All(context.TODO(), &productList); err != nil {
		return nil, err
	}

	return productList, nil
}
