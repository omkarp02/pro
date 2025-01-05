package cart

import (
	"context"
	"strconv"
	"time"

	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/utils"
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
	collection := s.getColl()

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "userId", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	return err
}

func (s *Repo) Create(ctx context.Context, payload CreateCartModel) error {
	collection := s.getColl()

	payload.Items.CartId = "C-" + strconv.Itoa(utils.GenerateRandomNumber(5))

	filter := bson.M{"userId": payload.UserId}
	update := bson.M{
		"$setOnInsert": bson.M{
			"userId":    payload.UserId,
			"createdAt": time.Now(),
		},
		"$set": bson.M{
			"updatedAt": time.Now(),
		},
		"$push": bson.M{
			"items": payload.Items,
		},
		"$inc": bson.M{
			"totalItems": 1,
			"totalPrice": payload.Price,
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	return nil

}
