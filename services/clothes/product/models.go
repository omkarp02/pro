package product

import (
	"github.com/omkarp02/pro/services/utils/store"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Variation struct {
	Size  string  `json:"size,omitempty" bson:"size,omitempty"`
	Price float64 `json:"price,omitempty" bson:"price,omitempty"`
}

type Price struct {
	BasePrice  float64     `json:"basePrice,omitempty" bson:"basePrice,omitempty"`
	Variations []Variation `json:"variations,omitempty" bson:"variations,omitempty"`
}

type ProductList struct {
	ID         bson.ObjectID    `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string           `json:"name,omitempty" bson:"name,omitempty"`
	Sizes      []string         `json:"sizes,omitempty" bson:"sizes,omitempty"`
	Color      string           `json:"color,omitempty" bson:"color,omitempty"`
	Price      Price            `json:"price,omitempty" bson:"price,omitempty"`
	ImgLink    string           `json:"imgLink,omitempty" bson:"imgLink,omitempty"`
	Stock      int              `json:"stock,omitempty" bson:"stock,omitempty"`
	Discount   int              `json:"discount,omitempty" bson:"discount,omitempty"`
	Detail     bson.ObjectID    `json:"detail,omitempty" bson:"detail,omitempty"`
	Timestamps store.Timestamps `bson:",inline"`
}
