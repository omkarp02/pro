package master

import (
	"github.com/omkarp02/pro/services/utils/store"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Master struct {
	ID          bson.ObjectID     `json:"_id,omitempty" bson:"_id,omitempty"`
	Type        string            `json:"type,omitempty" bson:"type,omitempty"`
	Name        string            `json:"name,omitempty" bson:"name,omitempty"`
	AuditFields store.AuditFields `bson:",inline"`
}
