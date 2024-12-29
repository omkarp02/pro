package useraccount

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/services/utils/helper"
	"github.com/omkarp02/pro/services/utils/store"
	"github.com/omkarp02/pro/utils/errutil"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Store struct {
	*db.Database
	collName string
}

func NewStore(curDb *db.Database, collName string) *Store {
	store := &Store{
		Database: curDb,
		collName: collName,
	}

	if err := store.createIndexes(); err != nil {
		log.Fatalf("Error while creating index of %s collection", collName)
	}

	return store
}

func (s *Store) getColl() *mongo.Collection {
	return s.DB.Database(s.DBName).Collection(s.collName)
}

func (s *Store) createIndexes() error {
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

func (s *Store) CreateUserAccount(user CreateUserAccountModal) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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
		fmt.Println("final", id.Hex(), id)
		return id.Hex(), nil
	}

	return "", errutil.ErrDatabase
}

func (s *Store) GetUserAccountByEmail(email string) (UserAccount, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var userAccount UserAccount

	filter := bson.D{{Key: "email", Value: email}}
	err := s.getColl().FindOne(ctx, filter).Decode(&userAccount)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return UserAccount{}, errutil.ErrDocumentNotFound
	} else if err != nil {
		slog.Error("error will getting user account", "err", err)
		return UserAccount{}, err
	}

	return userAccount, nil
}

func (s *Store) GetUserFromRefreshToken(refreshToken string) (UserAccount, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var userAccount UserAccount

	filter := bson.D{{Key: "refresh_token", Value: refreshToken}}
	project := bson.D{{Key: "refresh_token", Value: 1}, {Key: "_id", Value: 0}}
	findOptions := options.FindOne().SetProjection(project)
	err := s.getColl().FindOne(ctx, filter, findOptions).Decode(&userAccount)

	if errors.Is(err, mongo.ErrNoDocuments) {
		return UserAccount{}, errutil.ErrDocumentNotFound
	} else if err != nil {
		slog.Error("error will getting user account", "err", err)
		return UserAccount{}, err
	}

	return userAccount, nil
}

func (s *Store) GetUserAccount(query map[string]interface{}, project map[string]interface{}) (*UserAccount, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var userAccount UserAccount

	filterBson := bson.M(query)
	projectBson := bson.M(project)

	findOptions := options.FindOne().SetProjection(projectBson)

	err := s.getColl().FindOne(ctx, filterBson, findOptions).Decode(&userAccount)

	if err != nil {
		slog.Error("error will getting user account", "err", err)
		return nil, err
	}

	return &userAccount, nil
}

func (s *Store) UpdateUserAccountById(id string, userAccount UserAccount) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := store.CreateBsonFromStruct(userAccount)
	result, err := s.getColl().UpdateByID(ctx, id, update)
	if err != nil {
		return false, err
	}

	return result.Acknowledged, nil
}

func (s *Store) HandleRefreshTokenForLogin(userId string, refreshToken string, oldRefreshToken string) error {
	var action string

	if len(oldRefreshToken) == 0 {
		action = "push"
	} else {
		if err := s.PullUserRefreshToken(oldRefreshToken); err != nil {
			if errors.Is(err, errutil.ErrDocumentNotFound) {
				action = "reinitialize"
			} else {
				return err
			}
		} else {
			action = "push"
		}

	}

	if err := s.UpdateUserRefreshToken(userId, action, refreshToken); err != nil {
		return err
	}

	return nil
}

func (s *Store) UpdateUserRefreshToken(userId string, action string, refreshToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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

func (s *Store) PullUserRefreshToken(refreshToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

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

func (s *Store) createUserAccountModalFromData(userAccountData CreateUserAccountModal) UserAccount {

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
		AuthProvider: authProviderSlice,
	}

	return newUserAccount
}
