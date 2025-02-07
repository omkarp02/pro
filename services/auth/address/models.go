package address

import (
	"github.com/omkarp02/pro/services/utils/store"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Address struct {
	ID         bson.ObjectID    `bson:"_id,omitempty" json:"id,omitempty"`
	Address    store.Address    `bson:",inline" json:"address,inline"`
	IsPrimary  bool             `bson:"isPrimary" json:"isPrimary"`
	Type       ADDRESS_TYPE     `bson:"type,omitempty" json:"type,omitempty"`
	UserID     bson.ObjectID    `bson:"userId,omitempty" json:"userId,omitempty"`
	Timestamps store.Timestamps `json:"timestamp,omitempty" bson:",inline"`
}

type ADDRESS_TYPE = string

const (
	TYPE_HOME ADDRESS_TYPE = "home"
	TYPE_WORK ADDRESS_TYPE = "work"
)
