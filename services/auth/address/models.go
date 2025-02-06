package address

import (
	"github.com/omkarp02/pro/services/utils/store"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Address struct {
	ID         bson.ObjectID    `bson:"_id,omitempty"`
	Address    store.Address    `bson:",inline"`
	IsPrimary  bool             `bson:"isPrimary,omitempty"`
	Type       ADDRESS_TYPE     `bson:"type,omitempty"`
	UserID     bson.ObjectID    `bson:"userId,omitempty"`
	Timestamps store.Timestamps `json:"timestamp,omitempty" bson:",inline"`
}

type ADDRESS_TYPE = string

const (
	TYPE_HOME ADDRESS_TYPE = "home"
	TYPE_WORK ADDRESS_TYPE = "work"
)
