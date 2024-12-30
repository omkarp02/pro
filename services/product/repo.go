package product

import (
	"context"

	"github.com/omkarp02/pro/db"
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

	// if err := store.createIndexes(); err != nil {
	// 	log.Fatalf("Error while creating index of %s collection", collName)
	// }

	return store
}

func (s *Repo) getColl() *mongo.Collection {
	return s.DB.Database(s.DBName).Collection(s.collName)
}

func (s *Repo) createIndexes() error {
	collection := s.getColl()

	// Define the unique index for the "email" field
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	return err
}

// func (s *Repo) Create(ctx context.Context, createProductModal CreateProductListModal) (string, error) {

// 	productList := ProductList{
// 		Name:       createProductModal.Name,
// 		Sizes:      createProductModal.Sizes,
// 		Colors:     createProductModal.Colors,
// 		ImgLink:    createProductModal.ImgLink,
// 		Stock:      createProductModal.Stock,
// 		Price:      createProductModal.Price,
// 		Timestamps: store.GetCurrentTimestamps(),
// 	}

// 	result, err := s.getColl().InsertOne(ctx, productList)

// 	if mongo.IsDuplicateKeyError(err) {
// 		return "", errutil.ErrDocumentAlreadyExist
// 	} else if err != nil {
// 		return "", err
// 	}

// 	if id, ok := result.InsertedID.(bson.ObjectID); ok {
// 		fmt.Println("final", id.Hex(), id)
// 		return id.Hex(), nil
// 	}

// 	return "", errutil.ErrDatabase
// }

// func (s *Repo) Create(ctx context.Context, createProductModal CreateProductListModal) (string, error) {

// 	productList := ProductList{
// 		Name:       createProductModal.Name,
// 		Sizes:      createProductModal.Sizes,
// 		Colors:     createProductModal.Colors,
// 		ImgLink:    createProductModal.ImgLink,
// 		Stock:      createProductModal.Stock,
// 		Price:      createProductModal.Price,
// 		Timestamps: store.GetCurrentTimestamps(),
// 	}

// 	result, err := s.getColl().InsertOne(ctx, productList)

// 	if mongo.IsDuplicateKeyError(err) {
// 		return "", errutil.ErrDocumentAlreadyExist
// 	} else if err != nil {
// 		return "", err
// 	}

// 	if id, ok := result.InsertedID.(bson.ObjectID); ok {
// 		fmt.Println("final", id.Hex(), id)
// 		return id.Hex(), nil
// 	}

// 	return "", errutil.ErrDatabase
// }
