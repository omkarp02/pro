package test

import (
	"github.com/omkarp02/pro/services/utils/store"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Model struct {
	ID          bson.ObjectID     `json:"_id,omitempty" bson:"_id,omitempty"`
	AuditFields store.AuditFields `json:"auditFields" bson:",inline"`
	Timestamps  store.Timestamps  `bson:",inline"`
}
