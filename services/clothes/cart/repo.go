package cart

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/utils"
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

	userObjectId, err := bson.ObjectIDFromHex(payload.UserId)
	if err != nil {
		return err
	}

	var cartItems []CartItem

	totalItems := 0

	for _, item := range payload.Items {
		productId, err := bson.ObjectIDFromHex(item.ProductId)
		if err != nil {
			return err
		}

		totalItems += item.Quantity

		cartItems = append(cartItems, CartItem{
			CartId:    "C-" + strconv.Itoa(utils.GenerateRandomNumber(5)),
			ProductId: productId,
			Size:      item.Size,
			Quantity:  item.Quantity,
		})
	}

	filter := bson.M{"userId": userObjectId}
	update := bson.M{
		"$setOnInsert": bson.M{
			"userId":    userObjectId,
			"createdAt": time.Now(),
		},
		"$set": bson.M{
			"updatedAt": time.Now(),
		},
		"$push": bson.M{
			"items": bson.M{
				"$each": cartItems,
			},
		},
		"$inc": bson.M{
			"totalItems": payload.CurCartTotalItems,
			"totalPrice": payload.CurTotalPrice,
		},
	}

	opts := options.Update().SetUpsert(true)
	_, err = collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}

func (s *Repo) UpdateCartItemQuantity(ctx context.Context, userId string, cartId string, quantity int) error {

	userObjectId, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}

	filter := bson.M{
		"userId":       userObjectId,
		"items.cartId": cartId,
	}

	update := bson.M{
		"$set": bson.M{
			"items.$.quantity": quantity,
		},
	}
	result, err := s.getColl().UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errutil.NotFound("Cart Item")
	}

	return nil
}

func (s *Repo) FindById(ctx context.Context, userId string, project []string, exclusive bool) (Cart, error) {

	var cartDetails Cart

	objectId, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return cartDetails, err
	}

	filter := bson.M{"userId": objectId}
	findOneOptions := options.FindOne()

	if len(project) != 0 {
		projection := bson.M{}
		for _, field := range project {
			if exclusive {
				projection[field] = 0
			} else {
				projection[field] = 1
			}
		}
		findOneOptions.SetProjection(projection)
	}

	err = s.getColl().FindOne(ctx, filter, findOneOptions).Decode(&cartDetails)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return cartDetails, errutil.NotFound("Product")
		}
		return cartDetails, err
	}

	return cartDetails, nil
}
