package test

import (
	"github.com/omkarp02/pro/services/utils/store"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Model struct {
	ID          bson.ObjectID      `json:"_id,omitempty" bson:"_id,omitempty"`
	AuditFields *store.AuditFields `json:",inline,omitempty" bson:",inline,omitempty"`
	Timestamps  *store.Timestamps  `json:"timestamp,omitempty" bson:",inline"`
	CreatorId   string
}
