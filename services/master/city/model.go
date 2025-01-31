package city

import (
	"github.com/omkarp02/pro/services/utils/store"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type City struct {
	ID         bson.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
	Name       string            `json:"name,omitempty" bson:"name,omitempty"`
	StateId    bson.ObjectID     `json:"stateId,omitempty" bson:"stateId,omitempty"`
	Timestamps *store.Timestamps `json:"timestamp,omitempty" bson:",inline"`
}

type TCityRes struct {
	ID         bson.ObjectID     `json:"id,omitempty" bson:"_id,omitempty"`
	Name       string            `json:"name,omitempty" bson:"name,omitempty"`
	StateId    bson.ObjectID     `json:"stateId,omitempty" bson:"stateId,omitempty"`
	Timestamps *store.Timestamps `json:"timestamp,omitempty" bson:",inline"`
	State      any               `json:"stateDetails" bson:"stateDetails"`
}
