package userprofile

import (
	"github.com/omkarp02/pro/services/utils/store"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID         bson.ObjectID    `json:"_id,omitempty" bson:"_id,omitempty"`
	FullName   string           `json:"fullname,omitempty"`
	Age        int              `json:"age,omitempty"`
	Gender     string           `json:"gender,omitempty"`
	Timestamps store.Timestamps `bson:",inline"`
}
