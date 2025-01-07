package userprofile

import (
	"time"

	"github.com/omkarp02/pro/services/utils/store"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID          bson.ObjectID    `json:"_id,omitempty" bson:"_id,omitempty"`
	FirstName   string           `json:"firstname,omitempty" bson:"firstname,omitempty"`
	LastName    string           `json:"lastname,omitempty" bson:"lastname,omitempty"`
	Email       string           `json:"email,omitempty" bson:"email,omitempty"`
	DateOfBirth time.Time        `json:"dateofbirth,omitempty" bson:"dateofbirth,omitempty"`
	Gender      string           `json:"gender,omitempty" bson:"gender,omitempty"`
	Timestamps  store.Timestamps `json:"timestamp,omitempty" bson:",inline"`
}
