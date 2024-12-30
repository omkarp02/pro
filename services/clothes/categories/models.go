package categories

import (
	"github.com/omkarp02/pro/services/utils/store"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Category struct {
	ID         bson.ObjectID    `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string           `json:"name,omitempty" bson:"name"`
	Timestamps store.Timestamps `bson:",inline"`
}
