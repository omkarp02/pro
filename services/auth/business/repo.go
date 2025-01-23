package bussiness

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

type Repo struct {
	*db.Database
	collName string
}

func NewRepo(curDb *db.Database, collName string) *Repo {
	store := &Repo{
		Database: curDb,
		collName: collName,
	}

	store.createIndexes()

	return store
}

func (s *Repo) createIndexes() error {
	collection := s.getColl()

	// Define the unique index for the "email" field
	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "user_id", Value: 1}},
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	return err
}

func (s *Repo) getColl() *mongo.Collection {
	return s.DB.Database(s.DBName).Collection(s.collName)
}

func (s *Repo) Create(ctx context.Context, createPayload CreateModal) (string, error) {
	ownerObjectId, err := bson.ObjectIDFromHex(createPayload.OwnerID)
	if err != nil {
		return "", err
	}

	auditFields, err := store.GenerateCreateAuditFields(createPayload.CreatorId)
	if err != nil {
		return "", err
	}

	newAddress := Business{
		Name:        createPayload.Name,
		OwnerID:     ownerObjectId,
		Category:    createPayload.Category,
		Description: createPayload.Description,
		Address:     store.Address(createPayload.Address),
		Contacts:    createPayload.Contacts,
		Website:     createPayload.Website,
		LogoUrl:     createPayload.LogoUrl,
		Active:      createPayload.Active,
		AuditFields: auditFields,
	}

	result, err := s.getColl().InsertOne(ctx, newAddress)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return "", errutil.ErrDocumentAlreadyExist
		}
		return "", err
	}

	if id, ok := result.InsertedID.(bson.ObjectID); ok {
		return id.Hex(), nil
	}

	return "", errutil.ErrDatabase
}

func (s *Repo) FindByFilter(ctx context.Context, filterListModel FilterListModel, project []string, inclusive bool) ([]Business, error) {

	var list []Business

	query := bson.M{}

	name := filterListModel.Name
	page := filterListModel.Page
	limit := filterListModel.Limit

	if len(name) != 0 {
		query["name"] = bson.M{"$regex": name, "$options": "i"}
	}

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

func (s *Repo) FindById(ctx context.Context, id string, project []string, inclusive bool) (Business, error) {

	var business Business

	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return business, fmt.Errorf("invalid id format: %v", err)
	}

	filter := bson.M{"_id": objectId}
	findOneOptions := options.FindOne()

	if len(project) != 0 {
		projection := store.GenerateProjection(project, inclusive)
		findOneOptions.SetProjection(projection)
	}

	err = s.getColl().FindOne(ctx, filter, findOneOptions).Decode(&business)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return business, errutil.NotFound("Product")
		}
		return business, err
	}

	return business, nil

}
