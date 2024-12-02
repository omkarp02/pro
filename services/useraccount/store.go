package useraccount

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/omkarp02/pro/db"
	services "github.com/omkarp02/pro/services/utils"
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

func (s *Store) CreateUserAccount(user CreateUserAccountType) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	newUserAccount := UserAccount{
		Email:        user.Email,
		Password:     user.Password,
		RefreshToken: []string{},
		Timestamps:   services.GetCurrentTimestamps(),
	}

	fmt.Println(newUserAccount)

	result, err := s.getColl().InsertOne(ctx, newUserAccount)
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

func (s *Store) GetUserAccountByEmail(email string) (*UserAccount, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var userAccount UserAccount

	filter := bson.D{{Key: "email", Value: email}}
	err := s.getColl().FindOne(ctx, filter).Decode(&userAccount)

	if err != nil {
		slog.Error("error will getting user account", "err", err)
		return nil, err
	}

	return &userAccount, nil
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
		return UserAccount{}, services.ErrDocumentNotFound
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

	update := services.CreateBsonFromStruct(userAccount)
	result, err := s.getColl().UpdateByID(ctx, id, update)
	if err != nil {
		return false, err
	}

	return result.Acknowledged, nil
}

// func (s *Store) UpdateUserByRefreshToken(refreshToken string) (bool, error) {

// }

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
		return errors.New("invalid action")
	}

	objectId, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}

	result, err := s.getColl().UpdateByID(ctx, objectId, update)
	if result.MatchedCount == 0 {
		return services.ErrDocumentNotFound
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
		return services.ErrDocumentNotFound
	}
	if err != nil {
		return err
	}

	return nil

}
