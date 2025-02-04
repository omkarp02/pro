package product

import (
	"context"
	"errors"
	"log"

	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/services/utils"
	"github.com/omkarp02/pro/services/utils/store"
	"github.com/omkarp02/pro/utils/errutil"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ProductStatusRepo struct {
	*db.Database
	collName string
}

func NewProductStatusRepo(curDb *db.Database, collName string) *ProductStatusRepo {
	store := &ProductStatusRepo{
		Database: curDb,
		collName: collName,
	}

	if err := store.createIndexes(); err != nil {
		log.Fatalf("Error while creating index of %s collection", collName)
	}

	return store
}

func (s *ProductStatusRepo) getColl() *mongo.Collection {
	return s.DB.Database(s.DBName).Collection(s.collName)
}

func (s *ProductStatusRepo) createIndexes() error {
	collection := s.getColl()

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "productCode", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)

	return err
}

func (s *ProductStatusRepo) Create(ctx context.Context, createModal CreateProductStatusModal) (string, error) {

	timestamp := store.GetCurrentTimestamps()

	dataToInsert := ProductStatus{
		Timestamps:  &timestamp,
		ProductCode: createModal.ProductCode,
		Status:      createModal.Status,
	}

	result, err := s.getColl().InsertOne(ctx, dataToInsert)

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

func (s *ProductStatusRepo) UpdateByProductCode(ctx context.Context, id string, field string, value any) (utils.UPDATE_RESULT, error) {

	var updatedResult utils.UPDATE_RESULT

	query := bson.M{"productCode": id}
	update := bson.M{"$set": bson.M{field: value}}

	result, err := s.getColl().UpdateOne(ctx, query, update)
	if err != nil {
		return updatedResult, err
	}

	updatedResult.MatchedCount = int(result.MatchedCount)
	updatedResult.ModifiedCount = int(result.ModifiedCount)

	return updatedResult, nil
}

func (s *ProductStatusRepo) FindById(ctx context.Context, id string, project []string, inclusive bool) (ProductStatus, error) {

	var model ProductStatus

	filter := bson.M{"productCode": id}
	findOneOptions := options.FindOne()

	if len(project) != 0 {
		projection := store.GenerateProjection(project, inclusive)
		findOneOptions.SetProjection(projection)
	}

	err := s.getColl().FindOne(ctx, filter, findOneOptions).Decode(&model)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model, errutil.NotFound("Product Status")
		}
		return model, err
	}

	return model, nil

}
