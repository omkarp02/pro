package product

import (
	"github.com/omkarp02/pro/services/utils/store"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ProductList struct {
	ID         bson.ObjectID    `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string           `json:"name,omitempty" bson:"name,omitempty"`
	Sizes      []string         `json:"sizes,omitempty" bson:"sizes,omitempty"`
	Color      string           `json:"color,omitempty" bson:"color,omitempty"`
	Price      float64          `json:"price,omitempty" bson:"price,omitempty"`
	ImgLink    string           `json:"imgLink,omitempty" bson:"imgLink,omitempty"`
	Stock      int              `json:"stock,omitempty" bson:"stock,omitempty"`
	Discount   int              `json:"discount,omitempty" bson:"discount,omitempty"`
	Detail     bson.ObjectID    `json:"detail,omitempty" bson:"detail,omitempty"`
	Category   bson.ObjectID    `json:"category,omitempty" bson:"category,omitempty"`
	BatchId    string           `json:"batchId,omitempty" bson:"batchId,omitempty"`
	Timestamps store.Timestamps `bson:",inline"`
}

type ProductDetail struct {
	Description Description      `json:"description,omitempty" bson:"description,omitempty"`
	Variations  []Variation      `json:"variations,omitempty" bson:"variations,omitempty"`
	Timestamps  store.Timestamps `bson:",inline"`
}

type Variation struct {
	Size    string   `json:"size,omitempty" bson:"size,omitempty"`
	Price   float64  `json:"price,omitempty" bson:"price,omitempty"`
	ImgLink []string `json:"imgLink,omitempty" bson:"imgLink,omitempty"`
}

type Description struct {
	ProductDetails  string          `json:"productDetails,omitempty" bson:"productDetails,omitempty"`
	MaterialAndCare MaterialAndCare `json:"materialAndCare,omitempty" bson:"materialAndCare,omitempty"`
	Specifications  Specifications  `json:"specifications,omitempty" bson:"specifications,omitempty"`
}

type MaterialAndCare struct {
	Material         string `json:"material,omitempty" bson:"material,omitempty"`
	CareInstructions string `json:"careInstructions,omitempty" bson:"careInstructions,omitempty"`
}

type Specifications struct {
	SleeveLength    string `json:"sleeveLength,omitempty" bson:"sleeveLength,omitempty"`
	Collar          string `json:"collar,omitempty" bson:"collar,omitempty"`
	Fit             string `json:"fit,omitempty" bson:"fit,omitempty"`
	PatternType     string `json:"patternType,omitempty" bson:"patternType,omitempty"`
	Occasion        string `json:"occasion,omitempty" bson:"occasion,omitempty"`
	Length          string `json:"length,omitempty" bson:"length,omitempty"`
	Hemline         string `json:"hemline,omitempty" bson:"hemline,omitempty"`
	Placket         string `json:"placket,omitempty" bson:"placket,omitempty"`
	PlacketLength   string `json:"placketLength,omitempty" bson:"placketLength,omitempty"`
	Cuff            string `json:"cuff,omitempty" bson:"cuff,omitempty"`
	Transparency    string `json:"transparency,omitempty" bson:"transparency,omitempty"`
	WeavePattern    string `json:"weavePattern,omitempty" bson:"weavePattern,omitempty"`
	MainTrend       string `json:"mainTrend,omitempty" bson:"mainTrend,omitempty"`
	NumberOfItems   int    `json:"numberOfItems,omitempty" bson:"numberOfItems,omitempty"`
	PackageContains string `json:"packageContains,omitempty" bson:"packageContains,omitempty"`
}
