package useraccount

import (
	"context"
	"errors"
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
		RefreshToken: "",
		Timestamps:   services.GetCurrentTimestamps(),
	}

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

	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, nil
	} else if err != nil {
		slog.Error("error will getting user account", "err", err)
		return nil, err
	}

	return &userAccount, nil
}

func (s *Store) UpdateUserRefreshToken(userId bson.ObjectID, refreshToken string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{"refresh_token": refreshToken}
	filter := bson.D{{Key: "_id", Value: userId}}
	_, err := s.getColl().UpdateOne(ctx, filter, bson.M{"$set": update})

	return err
}
