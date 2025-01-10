package cart

import (
	"github.com/omkarp02/pro/services/clothes/product"
	"github.com/omkarp02/pro/services/utils/store"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func GetCartByUserIdAndPopulateUserQuery(userId bson.ObjectID, sizes []string) []bson.D {
	return mongo.Pipeline{
		// Match the documents by userId
		{{Key: "$match", Value: bson.D{
			{Key: "userId", Value: userId},
		}}},
		// Unwind the items array
		{{Key: "$unwind", Value: bson.D{
			{Key: "path", Value: "$items"},
			{Key: "preserveNullAndEmptyArrays", Value: true},
		}}},
		// Lookup the product details
		{{Key: "$lookup", Value: bson.D{
			{Key: "from", Value: "product_detail"},
			{Key: "let", Value: bson.D{
				{Key: "productId", Value: "$items.productId"},
			}},
			{Key: "pipeline", Value: mongo.Pipeline{
				{{Key: "$match", Value: bson.D{
					{Key: "$expr", Value: bson.D{
						{Key: "$eq", Value: bson.A{"$_id", "$$productId"}},
					}},
				}}},
				{{Key: "$project", Value: bson.D{
					{Key: "name", Value: 1},
					{Key: "previewImg", Value: 1},
					{Key: "variations", Value: bson.D{
						{Key: "$filter", Value: bson.D{
							{Key: "input", Value: "$variations"},
							{Key: "as", Value: "variation"},
							{Key: "cond", Value: bson.D{
								{Key: "$in", Value: bson.A{"$$variation.size", sizes}},
							}},
						}},
					}},
				}}},
			}},
			{Key: "as", Value: "items.product"},
		}}},
		// Set the first product in the items array (no longer an array)
		{{Key: "$set", Value: bson.D{
			{Key: "items.product", Value: bson.D{
				{Key: "$arrayElemAt", Value: bson.A{"$items.product", 0}},
			}},
		}}},
		// Group the documents back to their original structure
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: "$_id"},
			{Key: "userId", Value: bson.D{{Key: "$first", Value: "$userId"}}},
			{Key: "items", Value: bson.D{{Key: "$push", Value: "$items"}}},
			{Key: "totalItems", Value: bson.D{{Key: "$first", Value: "$totalItems"}}},
			{Key: "totalPrice", Value: bson.D{{Key: "$first", Value: "$totalPrice"}}},
			{Key: "createdAt", Value: bson.D{{Key: "$first", Value: "$createdAt"}}},
			{Key: "updatedAt", Value: bson.D{{Key: "$first", Value: "$updatedAt"}}},
		}}},
	}
}

type TCartItemWithProductPopulated struct {
	CartId    string                `bson:"cartId,omitempty" json:"cartId,omitempty"`
	ProductId bson.ObjectID         `bson:"productId,omitempty" json:"productId,omitempty"`
	Size      string                `bson:"size,omitempty" json:"size,omitempty"`
	Quantity  int                   `bson:"quantity,omitempty" json:"quantity,omitempty"`
	Product   product.ProductDetail `bson:"product,omitempty" json:"product,omitempty"`
}

type TGetCartByUserIdAndPopulateUserQuery struct {
	ID         bson.ObjectID                   `json:"_id,omitempty" bson:"_id,omitempty"`
	UserId     bson.ObjectID                   `bson:"userId,omitempty" json:"userId,omitempty"`
	Items      []TCartItemWithProductPopulated `bson:"items,omitempty" json:"items,omitempty"`
	TotalItems int                             `bson:"totalItems,omitempty" json:"totalItems,omitempty"`
	TotalPrice float64                         `bson:"totalPrice,omitempty" json:"totalPrice,omitempty"`
	Timestamps store.Timestamps                `bson:",inline" json:"timestamp"`
}
