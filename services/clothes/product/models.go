package product

import (
	"github.com/omkarp02/pro/services/utils/store"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ProductList struct {
	ID          bson.ObjectID      `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Sizes       []string           `json:"sizes,omitempty" bson:"sizes,omitempty"`
	Color       string             `json:"color,omitempty" bson:"color,omitempty"`
	Price       float64            `json:"price,omitempty" bson:"price,omitempty"`
	ImgLink     string             `json:"imgLink,omitempty" bson:"imgLink,omitempty"`
	Stock       int                `json:"stock,omitempty" bson:"stock,omitempty"`
	Discount    int                `json:"discount,omitempty" bson:"discount,omitempty"`
	Slug        string             `json:"slug,omitempty" bson:"slug,omitempty"`
	Detail      bson.ObjectID      `json:"detail,omitempty" bson:"detail,omitempty"`
	Category    bson.ObjectID      `json:"category,omitempty" bson:"category,omitempty"`
	BatchId     string             `json:"batchId,omitempty" bson:"batchId,omitempty"`
	Gender      string             `json:"gender,omitempty" bson:"gender,omitempty"`
	Collection  []string           `json:"collection,omitempty" bson:"collection,omitempty"`
	Tags        []string           `json:"tags,omitempty" bson:"tags,omitempty"`
	AuditFields *store.AuditFields `json:",inline,omitempty" bson:",inline,omitempty"`
}

type ProductDetail struct {
	ID          bson.ObjectID      `json:"id,omitempty" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Description Description        `json:"description,omitempty" bson:"description,omitempty"`
	Variations  []Variation        `json:"variations,omitempty" bson:"variations,omitempty"`
	ImgLink     []string           `json:"imgLink,omitempty" bson:"imgLink,omitempty"`
	BatchId     string             `json:"batchId,omitempty" bson:"batchId,omitempty"`
	Slug        string             `json:"slug,omitempty" bson:"slug,omitempty"`
	AuditFields *store.AuditFields `json:",inline,omitempty" bson:",inline,omitempty"`
}

type ProductTemplate struct {
	Name   string        `json:"name,omitempty" bson:"name,omitempty"`
	List   ProductList   `json:"size,omitempty" bson:"size,omitempty"`
	Detail ProductDetail `json:"detail,omitempty" bson:"detail,omitempty"`
}

type Variation struct {
	Size     string  `json:"size,omitempty" bson:"size,omitempty" validate:"required"`
	Price    float64 `json:"price,omitempty" bson:"price,omitempty" validate:"required"`
	Discount int     `json:"discount,omitempty" bson:"discount,omitempty" validate:"required,min=0,max=100"`
}

type Description struct {
	ProductDetails string         `json:"productDetails,omitempty" bson:"productDetails,omitempty"`
	Specifications Specifications `json:"specifications,omitempty" bson:"specifications,omitempty"`
}

type Specifications struct {
	SleeveLength    string `json:"sleeveLength,omitempty" bson:"sleeveLength,omitempty"`
	Collar          string `json:"collar,omitempty" bson:"collar,omitempty"`
	Fit             string `json:"fit,omitempty" bson:"fit,omitempty"`
	Fabric          string `json:"fabric,omitempty" bson:"fabric,omitempty"`
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

type BatchProductDetails struct {
	ProductListId   bson.ObjectID `json:"productListId,omitempty" bson:"productListId,omitempty"`
	ProductDetailId bson.ObjectID `json:"productDetailId,omitempty" bson:"productDetailId,omitempty"`
	ImgLink         string        `json:"imgLink,omitempty" bson:"imgLink,omitempty"`
}

type ProductBatch struct {
	ID          bson.ObjectID         `bson:"_id,omitempty" json:"id,omitempty"`
	Code        string                `bson:"code,omitempty" json:"code,omitempty"`
	Name        string                `bson:"name,omitempty" json:"name,omitempty"`
	ProductList []BatchProductDetails `bson:"batchProductDetails,omitempty" json:"batchProductDetails,omitempty"`
	AuditFields *store.AuditFields    `json:"timestamp,inline,omitempty" bson:",inline,omitempty"`
}
