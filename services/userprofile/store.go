package userprofile

import (
	"context"
	"time"

	"github.com/omkarp02/pro/db"
	services "github.com/omkarp02/pro/services/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Store struct {
	*db.Database
	collName string
}

func NewStore(curDb *db.Database, collName string) *Store {
	return &Store{
		Database: curDb,
		collName: collName,
	}
}

func (s *Store) getColl() *mongo.Collection {
	return s.DB.Database(s.DBName).Collection(s.collName)
}

func (s *Store) CreateUser(user CreateUser) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	newUser := User{
		FullName:   user.FullName,
		Gender:     user.Gender,
		Age:        user.Age,
		Timestamps: services.GetCurrentTimestamps(),
	}

	result, err := s.getColl().InsertOne(ctx, newUser)
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}

func (s *Store) GetUser(userId string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user User
	defer cancel()

	objId, _ := bson.ObjectIDFromHex(userId)

	err := s.getColl().FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
