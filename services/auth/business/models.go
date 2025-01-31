package bussiness

import (
	"github.com/omkarp02/pro/services/utils/store"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Business struct {
	ID          bson.ObjectID     `json:"_id,omitempty" bson:"_id,omitempty"`
	Name        string            `json:"name,omitempty" bson:"name,omitempty"`
	OwnerID     bson.ObjectID     `json:"ownerId,omitempty" bson:"ownerId,omitempty"`
	Category    string            `json:"category,omitempty" bson:"category,omitempty"`
	Description string            `json:"description,omitempty" bson:"description,omitempty"`
	Address     store.Address     `json:"address,omitempty" bson:"address,omitempty"`
	Contacts    []Contact         `json:"contacts,omitempty" bson:"contacts,omitempty"`
	Website     string            `json:"website,omitempty" bson:"website,omitempty"`
	LogoUrl     string            `json:"logoUrl,omitempty" bson:"logoUrl,omitempty"`
	Status      string            `json:"status,omitempty" bson:"status,omitempty"`
	AuditFields store.AuditFields `bson:",inline"`
}

type Contact struct {
	Type  string `json:"type,omitempty" bson:"type,omitempty"`
	Value string `json:"value,omitempty" bson:"value,omitempty"`
}
