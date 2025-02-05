package cart

import (
	"github.com/omkarp02/pro/services/utils/store"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type CartItem struct {
	CartId      string `bson:"cartId,omitempty" json:"cartId,omitempty"`
	ProductCode string `bson:"productCode,omitempty" json:"productCode,omitempty"`
	Size        string `bson:"size,omitempty" json:"size,omitempty"`
	Quantity    int    `bson:"quantity,omitempty" json:"quantity,omitempty"`
}

type Cart struct {
	ID     bson.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	UserId bson.ObjectID `bson:"userId,omitempty" json:"userId,omitempty"`
	Items  []CartItem    `bson:"items,omitempty" json:"items,omitempty"`
	// TotalItems int              `bson:"totalItems,omitempty" json:"totalItems,omitempty"`
	// TotalPrice float64          `bson:"totalPrice,omitempty" json:"totalPrice,omitempty"`
	Timestamps store.Timestamps `bson:"timestamp,inline" json:"timestamp"`
}
