package owner

import (
	"context"
	"errors"
	"fmt"
	"log"

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

	if err := store.createIndexes(); err != nil {
		log.Fatalf("Error while creating index of %s collection", collName)
	}

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

func (s *Repo) Create(ctx context.Context, createOwnerModal CreateModal) (string, error) {

	auditFields, err := store.GenerateCreateAuditFields(createOwnerModal.CreatorId)
	if err != nil {
		return "", err
	}

	owner := Owner{
		Name:        createOwnerModal.Name,
		FirstName:   createOwnerModal.FirstName,
		LastName:    createOwnerModal.LastName,
		DateOfBirth: createOwnerModal.DateOfBirth,
		Gender:      createOwnerModal.Gender,
		Email:       createOwnerModal.Email,
		AuditFields: auditFields,
	}

	result, err := s.getColl().InsertOne(ctx, owner)

	if mongo.IsDuplicateKeyError(err) {
		return "", errutil.ErrDocumentAlreadyExist
	} else if err != nil {
		return "", err
	}

	if id, ok := result.InsertedID.(bson.ObjectID); ok {
		fmt.Println("final", id.Hex(), id)
		return id.Hex(), nil
	}

	return "", errutil.ErrDatabase

}

func (s *Repo) FindByFilter(ctx context.Context, filterOwnerListModel FilterOwnerListModel, project []string, inclusive bool) ([]Owner, error) {

	var ownerList []Owner

	query := bson.M{}

	name := filterOwnerListModel.Name
	page := filterOwnerListModel.Page
	limit := filterOwnerListModel.Limit

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
	if err := cursor.All(context.TODO(), &ownerList); err != nil {
		return nil, err
	}

	return ownerList, nil
}

func (s *Repo) FindById(ctx context.Context, id string, project []string, inclusive bool) (Owner, error) {

	var owner Owner

	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return owner, fmt.Errorf("invalid id format: %v", err)
	}

	filter := bson.M{"_id": objectId}
	findOneOptions := options.FindOne()

	if len(project) != 0 {
		projection := store.GenerateProjection(project, inclusive)
		findOneOptions.SetProjection(projection)
	}

	err = s.getColl().FindOne(ctx, filter, findOneOptions).Decode(&owner)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return owner, errutil.NotFound("Product")
		}
		return owner, err
	}

	return owner, nil

}

func (s *Repo) AddBussiness(ctx context.Context, id string, updatedBy string, bussinessId string) error {

	objectIds, err := store.SliceOfHexToObjectID(bussinessId, updatedBy)
	if err != nil {
		return err
	}

	update := bson.M{
		"$push": bson.M{
			"businesses": objectIds[0], // New name to update
		},
		"$set": store.GenerateUpdateAudit(objectIds[1]),
	}
	_, err = s.getColl().UpdateByID(ctx, id, update)
	if err != nil {
		return err
	}

	return nil
}
