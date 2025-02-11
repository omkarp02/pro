package filter

import (
	"github.com/omkarp02/pro/services/utils/store"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Filter struct {
	ID         bson.ObjectID    `json:"id,omitempty" bson:"_id,omitempty"`
	Type       bson.ObjectID    `json:"type,omitempty" bson:"type,omitempty"`
	Name       string           `json:"name,omitempty" bson:"name,omitempty"`
	Slug       string           `json:"slug,omitempty" bson:"slug,omitempty"`
	Status     string           `json:"status,omitempty" bson:"status,omitempty"`
	Category   string           `json:"category,omitempty" bson:"category,omitempty"`
	Timestamps store.Timestamps `json:"timestamp,omitempty" bson:",inline"`
}

type FilterType struct {
	ID         bson.ObjectID    `json:"id,omitempty" bson:"_id,omitempty"`
	Name       string           `json:"name,omitempty" bson:"name,omitempty"`
	Status     string           `json:"status,omitempty" bson:"status,omitempty"`
	Timestamps store.Timestamps `json:"timestamp,omitempty" bson:",inline"`
}
