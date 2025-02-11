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

	codeIndexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "code", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	priceIndexModal := mongo.IndexModel{
		Keys: bson.D{{Key: "price", Value: 1}},
	}

	_, err := collection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{collectionIndexModal, priceIndexModal, codeIndexModel})

	return err
}

func (s *ProductListRepo) getColl() *mongo.Collection {
	return s.DB.Database(s.DBName).Collection(s.collName)
}

func (s *ProductListRepo) Create(ctx context.Context, createProductListModel CreateProductListModel, creatorId string) (string, error) {

	ids, err := store.SliceOfHexToObjectID(createProductListModel.Detail)
	if err != nil {
		return "", err
	}

	auditFields, err := store.GenerateCreateAuditFields(creatorId)
	if err != nil {
		return "", err
	}

	owner := ProductList{
		Detail:      ids[0],
		Name:        createProductListModel.Name,
		Code:        createProductListModel.Code,
		Sizes:       createProductListModel.Sizes,
		Color:       createProductListModel.Color,
		ImgLink:     createProductListModel.ImgLink,
		Price:       createProductListModel.Price,
		Slug:        createProductListModel.Slug,
		Stock:       createProductListModel.Stock,
		Discount:    createProductListModel.Discount,
		AuditFields: &auditFields,
		Category:    createProductListModel.Category,
		BatchId:     createProductListModel.BatchId,
		Gender:      createProductListModel.Gender,
		Collection:  createProductListModel.Collection,
		Tags:        createProductListModel.Tags,
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

func (s *ProductListRepo) FindByFilterBackup(ctx context.Context, filterProductListModel FilterProductListModel, project []string, inclusive bool) ([]ProductList, int, error) {

	var productList []ProductList

	query := bson.M{}

	sizes := filterProductListModel.Sizes
	name := filterProductListModel.Name
	color := filterProductListModel.Color
	maxPrice := filterProductListModel.MaxPrice
	minPrice := filterProductListModel.MinPrice
	collection := filterProductListModel.Collection
	category := filterProductListModel.Category
	page := filterProductListModel.Page
	limit := filterProductListModel.Limit
	getCountFlag := filterProductListModel.Count

	if len(sizes) != 0 {
		query["sizes"] = bson.M{"$in": sizes}
	}
	if len(color) != 0 {
		query["color"] = color
	}
	if len(name) != 0 {
		query["name"] = bson.M{"$regex": name, "$options": "i"}
	}
	if len(collection) != 0 {
		query["collection"] = collection
	}
	if len(category) != 0 {
		query["category"] = category
	}
	if maxPrice != 0 && minPrice != 0 {
		query["price"] = bson.M{"$gte": minPrice, "$lte": maxPrice}
	} else if maxPrice != 0 {
		query["price"] = bson.M{"$lte": maxPrice}
	} else if minPrice != 0 {
		query["price"] = bson.M{"$gte": minPrice}
	}

	fmt.Println(query)

	findOptions := options.Find().SetSkip(int64(limit * (page - 1))).SetLimit(int64(limit))

	if len(project) > 0 {
		projection := store.GenerateProjection(project, inclusive)
		findOptions.SetProjection(projection)
	}

	respChan := make(chan bool)
	countChan := make(chan int64)
	errChan := make(chan error, 2)

	go func() {
		defer close(respChan)
		cursor, err := s.getColl().Find(ctx, query, findOptions)
		if err != nil {
			respChan <- false
			errChan <- err
			return
		}
		if err := cursor.All(context.TODO(), &productList); err != nil {
			respChan <- false
			errChan <- err
			return
		}

		errChan <- nil
		respChan <- true
	}()

	go func() {
		var count int64
		var err error
		if getCountFlag {
			count, err = s.getColl().CountDocuments(ctx, query)
			if err != nil {
				countChan <- 0
				errChan <- err
				return
			}
		}

		errChan <- nil
		countChan <- count
	}()

	count := <-countChan
	err := <-errChan
	if err != nil {
		return productList, 0, err
	}
	err = <-errChan
	<-respChan

	return productList, int(count), err
}

func (s *ProductListRepo) FindByFilter(ctx context.Context, filterProductListModel FilterProductListModel, project []string, inclusive bool) ([]ProductList, error) {

	var productList []ProductList

	query := bson.M{}

	sizes := filterProductListModel.Sizes
	name := filterProductListModel.Name
	color := filterProductListModel.Color
	maxPrice := filterProductListModel.MaxPrice
	minPrice := filterProductListModel.MinPrice
	collection := filterProductListModel.Collection
	category := filterProductListModel.Category
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
	if len(collection) != 0 {
		query["collection"] = collection
	}
	if len(category) != 0 {
		query["category"] = category
	}
	if maxPrice != 0 && minPrice != 0 {
		query["price"] = bson.M{"$gte": minPrice, "$lte": maxPrice}
	} else if maxPrice != 0 {
		query["price"] = bson.M{"$lte": maxPrice}
	} else if minPrice != 0 {
		query["price"] = bson.M{"$gte": minPrice}
	}

	fmt.Println(query)

	findOptions := options.Find().SetSkip(int64(limit * (page - 1))).SetLimit(int64(limit))

	if len(project) > 0 {
		projection := store.GenerateProjection(project, inclusive)
		findOptions.SetProjection(projection)
	}

	cursor, err := s.getColl().Find(ctx, query, findOptions)
	if err != nil {
		return productList, err
	}
	err = cursor.All(context.TODO(), &productList)
	return productList, err
}

func (s *ProductListRepo) AddProductsToCollection(ctx context.Context, addProductToCollectionModel AddProductToCollectionModel) error {

	listOfObjectIds, err := store.SliceOfHexToObjectID(addProductToCollectionModel.ProductId...)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": bson.M{"$in": listOfObjectIds}}
	update := bson.M{"$push": bson.M{"collection": addProductToCollectionModel.CollectionName}}

	result, err := s.getColl().UpdateMany(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errutil.NotFound("Product")
	}

	return nil

}
