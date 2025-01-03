package product

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

type ProductListRepo struct {
	*db.Database
	collName string
}

func NewProductListRepo(curDb *db.Database, collName string) *ProductListRepo {
	store := &ProductListRepo{
		Database: curDb,
		collName: collName,
	}

	store.createIndexes()

	return store
}

func (s *ProductListRepo) createIndexes() error {
	collection := s.getColl()

	// Define the unique index for the "email" field
	collectionIndexModal := mongo.IndexModel{
		Keys: bson.D{{Key: "collection", Value: 1}},
		Options: options.Index().SetPartialFilterExpression(bson.M{
			"tags": bson.M{
				"$exists": true,
				"$ne":     nil,
				"$not":    bson.M{"$size": 0},
			},
		}),
	}

	priceIndexModal := mongo.IndexModel{
		Keys: bson.D{{Key: "price", Value: 1}},
	}

	_, err := collection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{collectionIndexModal, priceIndexModal})

	return err
}

func (s *ProductListRepo) getColl() *mongo.Collection {
	return s.DB.Database(s.DBName).Collection(s.collName)
}

func (s *ProductListRepo) Create(ctx context.Context, createProductListModel CreateProductListModel) (string, error) {

	detailId, err := bson.ObjectIDFromHex(createProductListModel.Detail)
	if err != nil {
		return "", err
	}

	categoryId, err := bson.ObjectIDFromHex(createProductListModel.Category)
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
		Category:   categoryId,
		BatchId:    createProductListModel.BatchId,
		Gender:     createProductListModel.Gender,
		Collection: createProductListModel.Collection,
		Tags:       createProductListModel.Tags,
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

func (s *ProductListRepo) AddProductsToCollection(ctx context.Context, addProductToCollectionModel AddProductToCollectionModel) error {

	listOfObjectIds, err := store.SliceOfHexToObjectID(addProductToCollectionModel.ProductId)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": bson.M{"$in": listOfObjectIds}}
	update := bson.M{"$set": bson.M{"collection": addProductToCollectionModel.CollectionName}}

	result, err := s.getColl().UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}

	fmt.Println(result)

	return nil

}
