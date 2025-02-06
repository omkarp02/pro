package address

import (
	"context"

	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/services/utils/store"
	"github.com/omkarp02/pro/utils/errutil"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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

func (s *Repo) Create(ctx context.Context, createPayload CreateAddressModel) (string, error) {

	userObjectId, err := bson.ObjectIDFromHex(createPayload.UserID)
	if err != nil {
		return "", err
	}

	newAddress := Address{
		IsPrimary:  createPayload.IsPrimary,
		Type:       createPayload.Type,
		UserID:     userObjectId,
		Timestamps: store.GetCurrentTimestamps(),
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

func (s *Repo) UpdateById(ctx context.Context, updatePayload UpdateAddressModel) error {

	objectIds, err := store.SliceOfHexToObjectID(updatePayload.UserID, updatePayload.Id)

	if err != nil {
		return err
	}

	userId := objectIds[0]
	addressId := objectIds[1]

	newAddress := Address{
		IsPrimary: updatePayload.IsPrimary,
	}

	newAddress.Timestamps.UpdatedAt = store.GetCurrentTimestamps().UpdatedAt

	query := bson.M{"userId": userId, "_id": addressId}
	update := bson.M{"$set": newAddress}

	result, err := s.getColl().UpdateOne(ctx, query, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return errutil.InternalServerError()
	}

	return nil
}

func (s *Repo) Delete(ctx context.Context, addressId string, userId string) error {
	objectIds, err := store.SliceOfHexToObjectID(addressId, userId)
	if err != nil {
		return err
	}

	userObjectId := objectIds[0]
	addressObjectId := objectIds[1]

	query := bson.M{"userId": userObjectId, "_id": addressObjectId}

	res, err := s.getColl().DeleteOne(ctx, query)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		errutil.InternalServerError()
	}

	return err
}

func (s *Repo) GetAddressByUserId(ctx context.Context, userId string) ([]Address, error) {

	var addressList []Address

	userObjectId, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"userId": userObjectId}

	cursor, err := s.getColl().Find(ctx, filter)
	if err != nil {
		return addressList, nil
	}

	if err := cursor.All(context.TODO(), &addressList); err != nil {
		return nil, err
	}

	return addressList, nil
}
