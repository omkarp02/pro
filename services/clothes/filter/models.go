package filter

import (
	"github.com/omkarp02/pro/services/utils/store"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Filter struct {
	ID         bson.ObjectID    `json:"_id,omitempty" bson:"_id,omitempty"`
	Type       bson.ObjectID    `json:"type,omitempty" bson:"type,omitempty"`
	Name       string           `json:"name,omitempty" bson:"name,omitempty"`
	Category   bson.ObjectID    `json:"category,omitempty" bson:"category,omitempty"`
	Timestamps store.Timestamps `bson:",inline"`
}

type FilterType struct {
	ID         bson.ObjectID    `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string           `json:"name,omitempty" bson:"name,omitempty"`
	Timestamps store.Timestamps `bson:",inline"`
}
