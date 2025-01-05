package cart

import (
	"github.com/omkarp02/pro/services/utils/store"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type CartItem struct {
	CartId    string        `bson:"cartId,omitempty" json:"cartId,omitempty"`
	ProductId bson.ObjectID `bson:"productId,omitempty" json:"productId,omitempty"`
	Size      string        `bson:"size,omitempty" json:"size,omitempty"`
	Quantity  int           `bson:"quantity,omitempty" json:"quantity,omitempty"`
}

type Cart struct {
	UserId     string           `bson:"userId,omitempty" json:"userId,omitempty"`
	Items      []CartItem       `bson:"items,omitempty" json:"items,omitempty"`
	TotalItems int              `bson:"totalItems,omitempty" json:"totalItems,omitempty"`
	TotalPrice float64          `bson:"totalPrice,omitempty" json:"totalPrice,omitempty"`
	Timestamps store.Timestamps `bson:",inline"`
}
