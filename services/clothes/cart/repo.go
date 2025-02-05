package cart

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/omkarp02/pro/db"
	"github.com/omkarp02/pro/services/utils/store"
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
	productCodes := []string{}

	totalItems := 0

	for _, item := range payload.Items {
		totalItems += item.Quantity
		cartItem := CartItem{
			CartId:      "C-" + strconv.Itoa(utils.GenerateRandomNumber(5)),
			ProductCode: item.ProductCode,
			Size:        item.Size,
			Quantity:    item.Quantity,
		}

		cartItems = append(cartItems, cartItem)
		productCodes = append(productCodes, cartItem.ProductCode)
	}

	filter := bson.M{"userId": userObjectId, "items.productCode": bson.M{"$nin": productCodes}}
	// filter := bson.M{"userId": userObjectId}
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
		// "$inc": bson.M{
		// 	"totalItems": payload.CurCartTotalItems,
		// 	"totalPrice": payload.CurTotalPrice,
		// },
	}

	opts := options.Update().SetUpsert(true)
	_, err = collection.UpdateOne(ctx, filter, update, opts)
	fmt.Println("err", err)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return errutil.ErrDocumentAlreadyExist
		}
		return err
	}

	return nil
}

type CartItemFieldToUpdate string

func (s *Repo) UpdateCartItem(ctx context.Context, userId string, cartId string, field string, value interface{}) error {

	userObjectId, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return err
	}

	filter := bson.M{
		"userId":       userObjectId,
		"items.cartId": cartId,
	}

	fieldToUpdate := "items.$." + field

	update := bson.M{
		"$set": bson.M{
			fieldToUpdate: value,
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

func (s *Repo) FindById(ctx context.Context, userId string, project []string, inclusive bool) (Cart, error) {

	var cartDetails Cart

	objectId, err := bson.ObjectIDFromHex(userId)
	if err != nil {
		return cartDetails, err
	}

	filter := bson.M{"userId": objectId}
	findOneOptions := options.FindOne()

	if len(project) > 0 {
		projection := store.GenerateProjection(project, inclusive)
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
