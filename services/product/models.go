package product

import (
	"github.com/omkarp02/pro/services/utils/store"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ProductList struct {
	ID         bson.ObjectID    `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string           `json:"name,omitempty" bson:"name,omitempty"`
	Sizes      []string         `json:"sizes,omitempty" bson:"sizes,omitempty"`
	Colors     []string         `json:"colors,omitempty" bson:"colors,omitempty"`
	ImgLink    string           `json:"imgLink,omitempty" bson:"imgLink,omitempty"`
	Stock      int              `json:"stock,omitempty" bson:"stock,omitempty"`
	Price      float64          `json:"price,omitempty" bson:"price,omitempty"`
	Timestamps store.Timestamps `bson:",inline"`
}
