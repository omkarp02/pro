package order

import (
	"github.com/omkarp02/pro/services/utils/store"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ShippingInfo struct {
	Address store.Address `bson:"address,omitempty" json:"address,omitempty"`
}

type OrderItem struct {
	ID         bson.ObjectID    `json:"_id,omitempty" bson:"_id,omitempty"`
	ProductId  bson.ObjectID    `bson:"productId,omitempty" json:"productId,omitempty"`
	Price      float64          `bson:"price,omitempty" json:"price,omitempty"`
	Size       string           `bson:"size,omitempty" json:"size,omitempty"`
	Quantity   int              `bson:"quantity,omitempty" json:"quantity,omitempty"`
	ItemId     string           `bson:"itemId,omitempty" json:"itemId,omitempty"`
	OrderId    bson.ObjectID    `bson:"orderId,omitempty" json:"orderId,omitempty"`
	UserId     bson.ObjectID    `bson:"userId,omitempty" json:"userId,omitempty"`
	Status     string           `bson:"status,omitempty" json:"status,omitempty"`
	Timestamps store.Timestamps `bson:"timestamp,inline" json:"timestamp"`
}

type Order struct {
	ID            bson.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	UserId        bson.ObjectID   `bson:"userId,omitempty" json:"userId,omitempty"`
	ShippingInfo  ShippingInfo    `bson:"shippingInfo,omitempty" json:"shippingInfo,omitempty"`
	Items         []bson.ObjectID `bson:"items,omitempty" json:"items,omitempty"`
	TotalPrice    float64         `bson:"totalPrice,omitempty" json:"price,omitempty"`
	ShippingPrice float64         `bson:"shippingPrice,omitempty" json:"shippingPrice,omitempty"`
}
