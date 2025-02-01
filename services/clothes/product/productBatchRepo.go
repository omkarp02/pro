package product

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/services/utils/store"
	"github.com/omkarp02/pro/utils/errutil"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ProductBatchRepo struct {
	*db.Database
	collName string
}

func NewProductBatchRepo(curDb *db.Database, collName string) *ProductBatchRepo {
	store := &ProductBatchRepo{
		Database: curDb,
		collName: collName,
	}

	return store
}

func (s *ProductBatchRepo) createIndexes() error {
	collection := s.getColl()

	//This is code for create single index
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "code", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)

	return err
}

func (s *ProductBatchRepo) getColl() *mongo.Collection {
	return s.DB.Database(s.DBName).Collection(s.collName)
}

func (s *ProductBatchRepo) Create(ctx context.Context, createPayload CreateProductBatchModel) (string, error) {

	auditFields, err := store.GenerateCreateAuditFields(createPayload.CreatorId)
	if err != nil {
		return "", err
	}

	batchDetails := ProductBatch{
		Code:        createPayload.Code,
		Name:        createPayload.Name,
		ProductList: createPayload.ProductList,
		AuditFields: &auditFields,
	}

	result, err := s.getColl().InsertOne(ctx, batchDetails)

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

func (s *ProductBatchRepo) UpdateBatchImg(ctx context.Context, batchId string, productDetail TBatchProductDetails) error {

	objectIds, err := store.SliceOfHexToObjectID(productDetail.ProductDetailId, productDetail.ProductListId)
	if err != nil {
		return err
	}

	updatePayload := BatchProductDetails{
		ProductDetailId: objectIds[0],
		ProductListId:   objectIds[1],
		ImgLink:         productDetail.ImgLink,
	}

	fmt.Println(objectIds[0])

	query := bson.M{"code": batchId}
	update := bson.M{
		"$push": bson.M{"batchProductDetails": updatePayload},
		"$set":  bson.M{"timestamp.updatedAt": time.Now()},
	}

	result, err := s.getColl().UpdateOne(ctx, query, update)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if result.MatchedCount == 0 {
		return errutil.NotFound("product batch")
	}

	return nil
}

func (s *ProductBatchRepo) FindById(ctx context.Context, id string, project []string, inclusive bool) (ProductBatch, error) {

	var batchDetail ProductBatch

	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return batchDetail, fmt.Errorf("invalid id format: %v", err)
	}

	filter := bson.M{"_id": objectId}
	findOneOptions := options.FindOne()

	if len(project) != 0 {
		projection := store.GenerateProjection(project, inclusive)
		findOneOptions.SetProjection(projection)
	}

	err = s.getColl().FindOne(ctx, filter, findOneOptions).Decode(&batchDetail)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return batchDetail, errutil.NotFound("Product batch")
		}
		return batchDetail, errutil.NotFound("Product batch")
	}

	return batchDetail, nil

}

func (s *ProductBatchRepo) FindByCode(ctx context.Context, code string, project []string, inclusive bool) (ProductBatch, error) {

	var batchDetail ProductBatch

	filter := bson.M{"code": code}
	findOneOptions := options.FindOne()

	if len(project) != 0 {
		projection := store.GenerateProjection(project, inclusive)
		findOneOptions.SetProjection(projection)
	}

	err := s.getColl().FindOne(ctx, filter, findOneOptions).Decode(&batchDetail)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return batchDetail, errutil.NotFound("Product batch")
		}
		return batchDetail, errutil.NotFound("Product batch")
	}

	return batchDetail, nil

}

func (s *ProductBatchRepo) FindByFilter(ctx context.Context, filterListModel FilterProductBatchListModel, project []string, inclusive bool) ([]ProductBatch, error) {

	var list []ProductBatch

	query := bson.M{}

	page := filterListModel.Page
	limit := filterListModel.Limit
	creatorId := filterListModel.CreatorId

	creatorObjectId, err := bson.ObjectIDFromHex(creatorId)
	if err != nil {
		return list, err
	}

	query["createdBy"] = creatorObjectId

	findOptions := options.Find().SetSkip(int64(limit * (page - 1))).SetLimit(int64(limit))

	if len(project) > 0 {
		projection := store.GenerateProjection(project, inclusive)
		findOptions.SetProjection(projection)
	}

	cursor, err := s.getColl().Find(ctx, query, findOptions)
	if err != nil {
		return nil, err
	}
	if err := cursor.All(context.TODO(), &list); err != nil {
		return nil, err
	}

	return list, nil
}

// func (s *Repo) createIndexes() error {
// 	collection := s.getColl()

// 	// Define the unique index for the "email" field
// 	catIdIndexModel := mongo.IndexModel{
// 		Keys:    bson.D{{Key: "catId", Value: 1}},
// 		Options: options.Index().SetUnique(true),
// 	}

// 	slugIndexModal := mongo.IndexModel{
// 		Keys:    bson.D{{Key: "slug", Value: 1}},
// 		Options: options.Index().SetUnique(true),
// 	}

// 	_, err := collection.Indexes().CreateMany(context.Background(), []mongo.IndexModel{catIdIndexModel, slugIndexModal})
// 	return err
// }
