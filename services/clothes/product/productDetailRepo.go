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

	store.createIndexes()

	return store
}

func (s *ProductDetailRepo) createIndexes() error {
	collection := s.getColl()

	slugIndexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "slug", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(context.Background(), slugIndexModel)

	return err
}

func (s *ProductDetailRepo) getColl() *mongo.Collection {
	return s.DB.Database(s.DBName).Collection(s.collName)
}

func (s *ProductDetailRepo) Create(ctx context.Context, createProductDetailModel CreateProductDetailModel, creatorId string) (string, error) {

	auditFields, err := store.GenerateCreateAuditFields(creatorId)
	if err != nil {
		return "", err
	}

	productDetail := ProductDetail{
		Description: createProductDetailModel.Description,
		Variations:  createProductDetailModel.Variations,
		Slug:        createProductDetailModel.Slug,
		ImgLink:     createProductDetailModel.ImgLink,
		AuditFields: &auditFields,
		Name:        createProductDetailModel.Name,
		BatchId:     createProductDetailModel.BatchId,
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

func (s *ProductDetailRepo) FindBySlug(ctx context.Context, slug string, project []string, inclusive bool) (ProductDetail, error) {

	var productDetail ProductDetail

	filter := bson.M{"slug": slug}
	findOneOptions := options.FindOne()

	if len(project) != 0 {
		projection := bson.M{}
		for _, field := range project {
			if inclusive {
				projection[field] = 1
			} else {
				projection[field] = 0
			}
		}
		findOneOptions.SetProjection(projection)
	}

	err := s.getColl().FindOne(ctx, filter, findOneOptions).Decode(&productDetail)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return productDetail, errutil.NotFound("Product")
		}
		return productDetail, err
	}

	return productDetail, nil

}

func (s *ProductDetailRepo) GetProductsPriceBySizes(ctx context.Context, ids []string, sizes []string, project []string, inclusive bool) ([]ProductDetail, error) {

	var productDetailList []ProductDetail

	objectIds, err := store.SliceOfHexToObjectID(ids...)
	if err != nil {
		return nil, fmt.Errorf("invalid id format: %v", err)
	}

	projection := store.GenerateProjection(project, inclusive)
	projection = append(projection, bson.E{Key: "variations", Value: bson.D{
		{Key: "$filter", Value: bson.D{
			{Key: "input", Value: "$variations"},
			{Key: "as", Value: "variation"},
			{Key: "cond", Value: bson.D{
				{Key: "$in", Value: bson.A{"$$variation.size", sizes}},
			}},
		}},
	}})

	pipeline := mongo.Pipeline{
		{
			{Key: "$match", Value: bson.D{
				{Key: "_id", Value: bson.D{
					{Key: "$in", Value: objectIds},
				}},
			}},
		},
		{
			{Key: "$project", Value: projection},
		},
	}

	cursor, err := s.getColl().Aggregate(ctx, pipeline)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errutil.NotFound("Product")
		}
		return nil, err
	}

	if err := cursor.All(context.TODO(), &productDetailList); err != nil {
		return nil, err
	}

	return productDetailList, nil

}

func (s *ProductDetailRepo) FindByIds(ctx context.Context, ids []string, project []string, exclusive bool) ([]ProductDetail, error) {

	var productDetail []ProductDetail

	objectIds, err := store.SliceOfHexToObjectID(ids...)
	if err != nil {
		return productDetail, fmt.Errorf("invalid id format: %v", err)
	}

	filter := bson.M{"_id": bson.M{"$in": objectIds}}
	findOptions := options.Find()

	if len(project) != 0 {
		projection := bson.M{}
		for _, field := range project {
			if exclusive {
				projection[field] = 0
			} else {
				projection[field] = 1
			}
		}
		findOptions.SetProjection(projection)
	}

	cursor, err := s.getColl().Find(ctx, filter, findOptions)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return productDetail, errutil.NotFound("Product")
		}
		return productDetail, err
	}

	if err := cursor.All(context.TODO(), &productDetail); err != nil {
		return nil, err
	}

	return productDetail, nil
}
