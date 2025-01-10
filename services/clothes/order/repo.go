package order

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

func (s *Repo) getColl() *mongo.Collection {
	return s.DB.Database(s.DBName).Collection(s.collName)
}

func (s *Repo) createIndexes() error {
	collection := s.getColl()

	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "userId", Value: 1}},
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	return err
}

func (s *Repo) Create(ctx context.Context, payload CreateOrderModal) (string, error) {
	collection := s.getColl()

	userObjectId, err := bson.ObjectIDFromHex(payload.UserId)
	if err != nil {
		return "", err
	}

	itemsObjectId, err := store.SliceOfHexToObjectID(payload.Items)
	if err != nil {
		return "", err
	}

	formattedOrder := Order{
		UserId: userObjectId,
		ShippingInfo: ShippingInfo{
			Address: store.Address(payload.ShippingInfo.Address),
		},
		Items:         itemsObjectId,
		TotalPrice:    payload.TotalPrice,
		ShippingPrice: payload.ShippingPrice,
	}

	result, err := collection.InsertOne(ctx, formattedOrder)
	if err != nil {
		return "", err
	}

	if id, ok := result.InsertedID.(bson.ObjectID); ok {
		return id.Hex(), nil
	}

	return "", errutil.ErrDatabase
}
