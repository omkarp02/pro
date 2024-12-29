package owner

import (
	"github.com/omkarp02/pro/services/utils/store"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Owner struct {
	ID         bson.ObjectID    `json:"_id,omitempty" bson:"_id,omitempty"`
	Name       string           `json:"name,omitempty" bson:"name"`
	Email      string           `json:"email,omitempty" bson:"email"`
	Businesses []bson.ObjectID  `json:"businesses,omitempty" bson:"businesses"`
	Timestamps store.Timestamps `bson:",inline"`
}
