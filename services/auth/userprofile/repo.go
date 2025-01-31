package userprofile

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
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true).SetPartialFilterExpression(bson.D{{Key: "email", Value: bson.D{{Key: "$exists", Value: true}, {Key: "$ne", Value: nil}}}}),
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	return err
}

func (s *Repo) getColl() *mongo.Collection {
	return s.DB.Database(s.DBName).Collection(s.collName)
}

func (s *Repo) Create(ctx context.Context, user CreateUserModel) (string, error) {

	var updatedUser User

	newUser := User{
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		DateOfBirth: user.DateOfBirth,
		Email:       user.Email,
		Gender:      user.Gender,
		Timestamps:  store.GetCurrentTimestamps(),
	}

	filter := bson.M{"email": user.Email}
	update := bson.M{"$set": newUser}

	opts := options.FindOneAndUpdate().SetUpsert(true)

	err := s.getColl().FindOneAndUpdate(ctx, filter, update, opts).Decode(&updatedUser)
	if err != nil {
		fmt.Println(err)
		if mongo.IsDuplicateKeyError(err) {
			return "", errutil.ErrDocumentAlreadyExist
		}
		return "", err
	}

	return updatedUser.ID.Hex(), nil
}

func (s *Repo) FindById(ctx context.Context, id string, project []string, inclusive bool) (User, error) {

	var model User

	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return model, fmt.Errorf("invalid id format: %v", err)
	}

	filter := bson.M{"_id": objectId}
	findOneOptions := options.FindOne()

	if len(project) != 0 {
		projection := store.GenerateProjection(project, inclusive)
		findOneOptions.SetProjection(projection)
	}

	err = s.getColl().FindOne(ctx, filter, findOneOptions).Decode(&model)
	if err != nil {
		return model, err
	}

	return model, nil

}
