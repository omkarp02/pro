package categories

import (
	"github.com/omkarp02/pro/services/utils/store"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Category struct {
	ID          bson.ObjectID    `json:"_id,omitempty" bson:"_id,omitempty"`
	CatId       string           `json:"catId,omitempty" bson:"catId,omitempty"`
	Name        string           `json:"name,omitempty" bson:"name,omitempty"`
	Description string           `json:"description,omitempty" bson:"description,omitempty"`
	ImgLink     string           `json:"imgLink,omitempty" bson:"imgLink,omitempty"`
	Slug        string           `json:"slug,omitempty" bson:"slug,omitempty"`
	Icon        string           `json:"icon,omitempty" bson:"icon,omitempty"`
	IsActive    bool             `json:"isActive,omitempty" bson:"isActive,omitempty"`
	Timestamps  store.Timestamps `bson:",inline"`
}
