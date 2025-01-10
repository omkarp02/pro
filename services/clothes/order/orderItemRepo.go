package order

import (
	"context"

	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/services/utils/store"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type OrderItemRepo struct {
	*db.Database
	collName string
}

func OrderItemNewRepo(curDb *db.Database, collName string) *OrderItemRepo {
	store := &OrderItemRepo{
		Database: curDb,
		collName: collName,
	}

	store.createIndexes()

	return store
}

func (s *OrderItemRepo) getColl() *mongo.Collection {
	return s.DB.Database(s.DBName).Collection(s.collName)
}

func (s *OrderItemRepo) createIndexes() error {
	collection := s.getColl()

	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "userId", Value: 1}},
	}

	_, err := collection.Indexes().CreateOne(context.Background(), indexModel)
	return err
}

func (s *OrderItemRepo) Create(ctx context.Context, payload CreateOrderItemModel) error {
	collection := s.getColl()

	formattedOrder, err := s.structureOrderItem(payload)
	if err != nil {
		return err
	}

	_, err = collection.InsertOne(ctx, formattedOrder)
	if err != nil {
		return err
	}

	return nil
}

func (s *OrderItemRepo) InsertMany(ctx context.Context, payload []CreateOrderItemModel) ([]string, error) {
	collection := s.getColl()

	orderList := []OrderItem{}

	for _, order := range payload {
		formattedOrder, err := s.structureOrderItem(order)
		if err != nil {
			return nil, err
		}

		orderList = append(orderList, formattedOrder)
	}

	result, err := collection.InsertMany(ctx, orderList)
	if err != nil {
		return nil, err
	}

	return store.SliceOfObjectIDToHex(result.InsertedIDs)
}

func (s *OrderItemRepo) structureOrderItem(payload CreateOrderItemModel) (OrderItem, error) {
	listOfIds := []string{payload.ProductId, payload.OrderId, payload.UserId}
	listOfObjectId, err := store.SliceOfHexToObjectID(listOfIds)

	orderPayload := OrderItem{
		ProductId:  listOfObjectId[0],
		Price:      payload.Price,
		Size:       payload.Size,
		Quantity:   payload.Quantity,
		ItemId:     payload.ItemId,
		OrderId:    listOfObjectId[1],
		UserId:     listOfObjectId[2],
		Status:     payload.Status,
		Timestamps: store.GetCurrentTimestamps(),
	}

	return orderPayload, err
}
