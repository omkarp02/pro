package useraccount

import (
	"context"
	"errors"
	"time"

	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/services/utils/helper"
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

func (s *Repo) getColl() *mongo.Collection {
	return s.DB.Database(s.DBName).Collection(s.collName)
}

func (s *Repo) createIndexes() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := s.getColl()

	// Define the unique index for the "email" field
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

func (s *Repo) Create(ctx context.Context, user CreateUserAccountModal) (string, error) {
	if len(user.PasswordHash) != 0 {
		hashedPassword, err := helper.HashPassword(user.PasswordHash)
		if err != nil {
			return "", err
		}
		user.PasswordHash = hashedPassword
	}

	newUserAccount := s.createUserAccountModalFromData(user)

	result, err := s.getColl().InsertOne(ctx, newUserAccount)

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

func (s *Repo) FindById(ctx context.Context, id string, project []string, inclusive bool) (UserAccount, error) {

	var userAccount UserAccount

	objectId, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return userAccount, err
	}

	filter := bson.M{"_id": objectId}
	findOneOptions := options.FindOne()

	if len(project) != 0 {
		projection := store.GenerateProjection(project, inclusive)
		findOneOptions.SetProjection(projection)
	}

	err = s.getColl().FindOne(ctx, filter, findOneOptions).Decode(&userAccount)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return userAccount, errutil.NotFound("Product")
		}
		return userAccount, err
	}

	return userAccount, nil

}

func (s *Repo) FindOne(ctx context.Context, field string, value string, project []string, inclusive bool) (UserAccount, error) {

	var userAccount UserAccount

	filter := bson.M{field: value}
	findOneOptions := options.FindOne()

	if len(project) != 0 {
		projection := store.GenerateProjection(project, inclusive)
		findOneOptions.SetProjection(projection)
	}

	err := s.getColl().FindOne(ctx, filter, findOneOptions).Decode(&userAccount)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return userAccount, errutil.NotFound("Product")
		}
		return userAccount, err
	}

	return userAccount, nil
}

func (s *Repo) UpdateUserProfileById(ctx context.Context, id string, userProfileId string) error {

	objectId, err := bson.ObjectIDFromHex(userProfileId)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"userProfileId": objectId, // New name to update
		},
	}
	_, err = s.getColl().UpdateByID(ctx, id, update)
	if err != nil {
		return err
	}

	return nil
}

func (s *Repo) UpdateUserRefreshToken(ctx context.Context, userId string, action string, refreshToken string) error {
	var update bson.M

	switch action {
	case "push":
		// Append the new refresh token to the array
		update = bson.M{"$push": bson.M{"refresh_token": refreshToken}}
	case "reinitialize":
		update = bson.M{"$set": bson.M{"refresh_token": []string{refreshToken}}}
	case "pull":
		// Remove the specific refresh token from the array
		update = bson.M{"$pull": bson.M{"refresh_token": refreshToken}}
	case "empty":
		// Set the refresh_token array to an empty string (or clear it)
		update = bson.M{"$set": bson.M{"refresh_token": []string{}}}
	default:
		return errutil.InternalServerError("Invalid Action")
	}

	objectId, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}

	result, err := s.getColl().UpdateByID(ctx, objectId, update)
	if result.MatchedCount == 0 {
		return errutil.ErrDocumentNotFound
	}
	if err != nil {
		return err
	}

	return nil
}

func (s *Repo) PullUserRefreshToken(ctx context.Context, refreshToken string) error {

	query := bson.M{"refresh_token": refreshToken}
	update := bson.M{"$pull": bson.M{"refresh_token": refreshToken}}
	result, err := s.getColl().UpdateOne(ctx, query, update)
	if result.MatchedCount == 0 {
		return errutil.ErrDocumentNotFound
	}
	if err != nil {
		return err
	}

	return nil

}

func (s *Repo) createUserAccountModalFromData(userAccountData CreateUserAccountModal) UserAccount {

	authProviderSlice := []AuthProvider{}
	for _, auth := range userAccountData.AuthProvider {
		authProviderSlice = append(authProviderSlice, AuthProvider{
			Provider:   auth.Provider,
			ProviderID: auth.ProviderID,
		})
	}

	newUserAccount := UserAccount{
		Email:        userAccountData.Email,
		PasswordHash: userAccountData.PasswordHash,
		Timestamps:   store.GetCurrentTimestamps(),
		Role:         userAccountData.Role,
		AuthProvider: authProviderSlice,
	}

	if len(userAccountData.UserProfile) != 0 {
		userProfileObjectId, err := bson.ObjectIDFromHex(userAccountData.UserProfile)

		newUserAccount.UserProfile = userProfileObjectId

		if err != nil {
			panic(err)
		}
	}

	return newUserAccount
}
