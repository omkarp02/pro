package state

import (
	"github.com/omkarp02/pro/services/utils/store"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type State struct {
	ID         bson.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
	Name       string            `json:"name,omitempty" bson:"name,omitempty"`
	Timestamps *store.Timestamps `json:"timestamp,omitempty" bson:",inline"`
}
